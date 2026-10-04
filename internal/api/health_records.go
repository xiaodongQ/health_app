package api

import (
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"health_app/configs"
	"health_app/internal/models"
	"health_app/internal/repository"
	"health_app/internal/service"

	"github.com/gin-gonic/gin"
)

var healthRecordService *service.HealthRecordService

// getHealthRecordService 获取健康记录服务实例（懒加载）
func getHealthRecordService() *service.HealthRecordService {
	if healthRecordService == nil {
		healthRecordRepo := repository.NewHealthRecordRepository()
		healthRecordService = service.NewHealthRecordService(healthRecordRepo)
	}
	return healthRecordService
}

// getReferenceByIndicator 根据指标名称从配置中查找参考值
func getReferenceByIndicator(indicatorName string) (reference string, found bool) {
	for _, category := range configs.DefaultHealthConfig {
		for _, ind := range category.Indicators {
			if ind.Name == indicatorName {
				return ind.Reference, true
			}
		}
	}
	return "", false
}

// parseReferenceRange 解析参考值字符串，返回下限、上限和是否有效
// 支持格式: "44-133", "3.5-5.3", "阴性", "阳性", "0-3" 等
func parseReferenceRange(ref string) (min, max float64, isRange bool, isNegative bool) {
	ref = strings.TrimSpace(ref)

	// 处理"阴性/阳性"类枚举值
	negativeKeywords := []string{"阴性", "弱阳性", "阳性"}
	for _, kw := range negativeKeywords {
		if ref == kw || ref == kw+"或"+kw {
			return 0, 0, false, true
		}
	}

	// 尝试解析数值范围，格式如 "44-133" 或 "3.5-7.1"
	rangeRegex := regexp.MustCompile(`^([<>]?)\s*([\d.]+)\s*-\s*([\d.]+)$`)
	matches := rangeRegex.FindStringSubmatch(ref)
	if len(matches) == 4 {
		min, _ = strconv.ParseFloat(matches[2], 64)
		max, _ = strconv.ParseFloat(matches[3], 64)
		return min, max, true, false
	}

	// 处理 ">X" 或 "<X" 格式
	gtRegex := regexp.MustCompile(`^>\s*([\d.]+)$`)
	if matches := gtRegex.FindStringSubmatch(ref); len(matches) == 2 {
		min, _ = strconv.ParseFloat(matches[1], 64)
		return min, 0, false, false
	}
	ltRegex := regexp.MustCompile(`^<\s*([\d.]+)$`)
	if matches := ltRegex.FindStringSubmatch(ref); len(matches) == 2 {
		max, _ = strconv.ParseFloat(matches[1], 64)
		return 0, max, false, false
	}

	return 0, 0, false, false
}

// calculateAbnormal 根据参考范围判断数值是否异常
func calculateAbnormal(indicatorName string, value float64) bool {
	ref, found := getReferenceByIndicator(indicatorName)
	if !found || ref == "" {
		return false
	}

	min, max, isRange, isNegative := parseReferenceRange(ref)

	// 枚举类型：阴性/阳性类，只有0或极小值才算正常
	if isNegative {
		return value != 0
	}

	// 范围类型
	if isRange {
		return value < min || value > max
	}

	// 单边范围
	if min > 0 {
		return value <= min
	}
	if max > 0 {
		return value >= max
	}

	return false
}

// enrichResultsWithAbnormal 为结果列表填充异常标记和参考值
func enrichResultsWithAbnormal(results []models.TestResult) {
	for i := range results {
		ref, found := getReferenceByIndicator(results[i].Indicator)
		if found {
			results[i].Reference = ref
		}
		results[i].Abnormal = calculateAbnormal(results[i].Indicator, results[i].Value)
	}
}

// AddHealthRecord 添加健康记录（包括多个检测指标）
func AddHealthRecord(c *gin.Context) {
	log.Printf("[API] 收到添加健康记录请求: %s %s", c.Request.Method, c.Request.URL.Path)

	// 只读取一次请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[API] 添加健康记录失败: 读取请求体错误 - %v", err)
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}

	// 解析完整请求数据（用字符串接收record_date避免time.Time的RFC3339严格解析）
	var requestData struct {
		models.HealthRecord
		RecordDateStr string            `json:"record_date"`
		Results       []models.TestResult `json:"results"`
	}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		log.Printf("[API] 添加健康记录失败: 解析请求体错误 - %v", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	record := requestData.HealthRecord

	// 手动解析简单日期格式 "2006-01-02"
	if parsed, err := time.Parse("2006-01-02", requestData.RecordDateStr); err == nil {
		record.RecordDate = parsed
	}

	// 验证必需字段
	if record.UserID == 0 || record.RecordDate.IsZero() || record.Category == "" {
		log.Printf("[API] 添加健康记录失败: 缺少必需字段 - UserID: %d, RecordDate: %v, Category: %s", record.UserID, record.RecordDate, record.Category)
		c.JSON(400, gin.H{"error": "User ID, Record Date, and Category are required"})
		return
	}

	// 自动计算异常标记和参考值
	enrichResultsWithAbnormal(requestData.Results)

	// 使用服务层添加记录
	service := getHealthRecordService()
	recordID, err := service.AddHealthRecord(&record, requestData.Results)
	if err != nil {
		log.Printf("[API] 添加健康记录失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to add health record: " + err.Error()})
		return
	}

	log.Printf("[API] 健康记录添加成功: record_id=%d, user_id=%d", recordID, record.UserID)
	c.JSON(201, gin.H{
		"message": "Health record added successfully",
		"id":      recordID,
	})
}

// GetHealthRecords 获取用户的所有健康记录
func GetHealthRecords(c *gin.Context) {
	log.Printf("[API] 收到获取健康记录请求: %s %s", c.Request.Method, c.Request.URL.Path)

	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[API] 获取健康记录失败: 无效的用户ID - %s", userIDStr)
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	// 使用服务层获取健康记录
	service := getHealthRecordService()
	records, err := service.GetHealthRecords(userID)
	if err != nil {
		log.Printf("[API] 获取健康记录失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve health records: " + err.Error()})
		return
	}

	log.Printf("[API] 健康记录查询成功: user_id=%d, 共 %d 条记录", userID, len(records))
	c.JSON(200, records)
}

// GetTrendData 获取特定指标的趋势数据
func GetTrendData(c *gin.Context) {
	log.Printf("[API] 收到获取趋势数据请求: %s %s", c.Request.Method, c.Request.URL.Path)

	userIDStr := c.Query("user_id")
	indicator := c.Query("indicator")
	unit := c.Query("unit") // 新增单位参数

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[API] 获取趋势数据失败: 无效的用户ID - %s", userIDStr)
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	if indicator == "" {
		log.Printf("[API] 获取趋势数据失败: 缺少指标参数")
		c.JSON(400, gin.H{"error": "Indicator parameter is required"})
		return
	}

	// 使用服务层获取趋势数据
	service := getHealthRecordService()
	trends, err := service.GetTrendData(userID, indicator, unit)
	if err != nil {
		log.Printf("[API] 获取趋势数据失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve trend data: " + err.Error()})
		return
	}

	log.Printf("[API] 趋势数据查询成功: user_id=%d, indicator=%s, unit=%s, 共 %d 条记录", userID, indicator, unit, len(trends))
	c.JSON(200, trends)
}

// GetAllIndicators 获取所有指标名称
func GetAllIndicators(c *gin.Context) {
	log.Printf("[API] 收到获取所有指标名称请求: %s %s", c.Request.Method, c.Request.URL.Path)

	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[API] 获取所有指标名称失败: 无效的用户ID - %s", userIDStr)
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	// 使用服务层获取所有指标名称
	service := getHealthRecordService()
	indicators, err := service.GetAllIndicators(userID)
	if err != nil {
		log.Printf("[API] 获取所有指标名称失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve indicators: " + err.Error()})
		return
	}

	log.Printf("[API] 指标名称查询成功: user_id=%d, 共 %d 个指标", userID, len(indicators))
	c.JSON(200, indicators)
}

// GetHealthRecordByID 根据ID获取单个健康记录
func GetHealthRecordByID(c *gin.Context) {
	log.Printf("[API] 收到获取单个健康记录请求: %s %s", c.Request.Method, c.Request.URL.Path)

	recordIDStr := c.Param("id")
	recordID, err := strconv.Atoi(recordIDStr)
	if err != nil {
		log.Printf("[API] 获取健康记录失败: 无效的记录ID - %s", recordIDStr)
		c.JSON(400, gin.H{"error": "Invalid Record ID"})
		return
	}

	// 使用服务层获取健康记录
	service := getHealthRecordService()
	record, err := service.GetHealthRecordByID(recordID)
	if err != nil {
		log.Printf("[API] 获取健康记录失败: 服务层错误 - %v", err)
		c.JSON(404, gin.H{"error": "Health record not found"})
		return
	}

	log.Printf("[API] 单个健康记录查询成功: record_id=%d", recordID)
	c.JSON(200, record)
}

// UpdateHealthRecord 更新健康记录
func UpdateHealthRecord(c *gin.Context) {
	log.Printf("[API] 收到更新健康记录请求: %s %s", c.Request.Method, c.Request.URL.Path)

	recordIDStr := c.Param("id")
	recordID, err := strconv.Atoi(recordIDStr)
	if err != nil {
		log.Printf("[API] 更新健康记录失败: 无效的记录ID - %s", recordIDStr)
		c.JSON(400, gin.H{"error": "Invalid Record ID"})
		return
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[API] 更新健康记录失败: 读取请求体错误 - %v", err)
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}

	// 解析请求数据（用字符串接收record_date避免time.Time的RFC3339严格解析）
	var requestData struct {
		models.HealthRecord
		RecordDateStr string            `json:"record_date"`
		Results       []models.TestResult `json:"results"`
	}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		log.Printf("[API] 更新健康记录失败: 解析请求体错误 - %v", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	record := requestData.HealthRecord

	// 手动解析简单日期格式 "2006-01-02"
	if parsed, err := time.Parse("2006-01-02", requestData.RecordDateStr); err == nil {
		record.RecordDate = parsed
	}

	// 验证必需字段
	if record.RecordDate.IsZero() || record.Category == "" {
		log.Printf("[API] 更新健康记录失败: 缺少必需字段 - RecordDate: %v, Category: %s", record.RecordDate, record.Category)
		c.JSON(400, gin.H{"error": "Record Date and Category are required"})
		return
	}

	// 自动计算异常标记和参考值
	enrichResultsWithAbnormal(requestData.Results)

	// 使用服务层更新记录
	service := getHealthRecordService()
	err = service.UpdateHealthRecord(recordID, &record, requestData.Results)
	if err != nil {
		log.Printf("[API] 更新健康记录失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to update health record: " + err.Error()})
		return
	}

	log.Printf("[API] 健康记录更新成功: record_id=%d", recordID)
	c.JSON(200, gin.H{
		"message": "Health record updated successfully",
		"id":      recordID,
	})
}

// GetHealthRecordByDate 获取指定用户和日期的健康记录
func GetHealthRecordByDate(c *gin.Context) {
	log.Printf("[API] 收到获取指定日期健康记录请求: %s %s", c.Request.Method, c.Request.URL.Path)

	userIDStr := c.Query("user_id")
	dateStr := c.Query("date")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("[API] 获取健康记录失败: 无效的用户ID - %s", userIDStr)
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	if dateStr == "" {
		log.Printf("[API] 获取健康记录失败: 缺少日期参数")
		c.JSON(400, gin.H{"error": "Date parameter is required"})
		return
	}

	// 使用服务层获取指定用户和日期的健康记录
	service := getHealthRecordService()
	record, err := service.GetHealthRecordByDate(userID, dateStr)
	if err != nil {
		if record == nil {
			// 如果没有找到记录，返回空响应而不是错误
			log.Printf("[API] 未找到指定日期的健康记录: user_id=%d, date=%s", userID, dateStr)
			c.JSON(200, nil)
			return
		}
		log.Printf("[API] 获取健康记录失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve health record: " + err.Error()})
		return
	}

	log.Printf("[API] 指定日期健康记录查询成功: user_id=%d, date=%s", userID, dateStr)
	c.JSON(200, record)
}

// DeleteHealthRecord 删除健康记录
func DeleteHealthRecord(c *gin.Context) {
	log.Printf("[API] 收到删除健康记录请求: %s %s", c.Request.Method, c.Request.URL.Path)

	recordIDStr := c.Param("id")
	recordID, err := strconv.Atoi(recordIDStr)
	if err != nil {
		log.Printf("[API] 删除健康记录失败: 无效的记录ID - %s", recordIDStr)
		c.JSON(400, gin.H{"error": "Invalid Record ID"})
		return
	}

	// 使用服务层删除记录
	service := getHealthRecordService()
	err = service.DeleteHealthRecord(recordID)
	if err != nil {
		log.Printf("[API] 删除健康记录失败: 服务层错误 - %v", err)
		c.JSON(500, gin.H{"error": "Failed to delete health record: " + err.Error()})
		return
	}

	log.Printf("[API] 健康记录删除成功: record_id=%d", recordID)
	c.JSON(200, gin.H{
		"message": "Health record deleted successfully",
		"id":      recordID,
	})
}
