package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

// HealthIndicator 表示单个健康指标
type HealthIndicator struct {
	Name      string `json:"name"`
	Unit      string `json:"unit"`
	Reference string `json:"reference"` // 参考范围
}

// HealthCategory 表示健康检查类别
type HealthCategory struct {
	Name       string            `json:"name"`
	Indicators []HealthIndicator `json:"indicators"`
}

// GetHealthConfig 返回所有健康检查类别和对应指标的配置
func GetHealthConfig(c *gin.Context) {
	log.Printf("[API] 收到获取健康配置请求: %s %s", c.Request.Method, c.Request.URL.Path)

	// 预定义的健康检查类别和指标
	config := []HealthCategory{
		{
			Name: "生化筛查",
			Indicators: []HealthIndicator{
				{Name: "肌酐", Unit: "μmol/L", Reference: "44-133"},
				{Name: "尿素", Unit: "mmol/L", Reference: "2.5-7.1"},
				{Name: "钾", Unit: "mmol/L", Reference: "3.5-5.3"},
				{Name: "空腹血糖", Unit: "mmol/L", Reference: "3.9-6.1"},
				// {Name: "钠", Unit: "mmol/L", Reference: "137-147"},
				// {Name: "氯", Unit: "mmol/L", Reference: "98-107"},
				// {Name: "钙", Unit: "mmol/L", Reference: "2.1-2.6"},
				// {Name: "磷", Unit: "mmol/L", Reference: "0.85-1.51"},
				// {Name: "镁", Unit: "mmol/L", Reference: "0.72-1.15"},
				// {Name: "尿酸", Unit: "μmol/L", Reference: "208-428"},
				// {Name: "总蛋白", Unit: "g/L", Reference: "65-85"},
				// {Name: "白蛋白", Unit: "g/L", Reference: "40-55"},
				// {Name: "球蛋白", Unit: "g/L", Reference: "20-35"},
				// {Name: "白球比", Unit: "", Reference: "1.2-2.4"},
				// {Name: "总胆红素", Unit: "μmol/L", Reference: "3.4-20.5"},
				// {Name: "直接胆红素", Unit: "μmol/L", Reference: "0-6.8"},
				// {Name: "间接胆红素", Unit: "μmol/L", Reference: "1.7-13.7"},
				// {Name: "谷丙转氨酶", Unit: "U/L", Reference: "7-40"},
				// {Name: "谷草转氨酶", Unit: "U/L", Reference: "13-35"},
				// {Name: "碱性磷酸酶", Unit: "U/L", Reference: "45-125"},
				// {Name: "γ-谷氨酰转肽酶", Unit: "U/L", Reference: "10-60"},
				// {Name: "肌酸激酶", Unit: "U/L", Reference: "24-195"},
				// {Name: "乳酸脱氢酶", Unit: "U/L", Reference: "135-225"},
				// {Name: "α-羟丁酸脱氢酶", Unit: "U/L", Reference: "72-182"},
				// {Name: "总胆固醇", Unit: "mmol/L", Reference: "2.9-5.7"},
				// {Name: "甘油三酯", Unit: "mmol/L", Reference: "0.45-1.7"},
				// {Name: "高密度脂蛋白", Unit: "mmol/L", Reference: "1.0-1.6"},
				// {Name: "低密度脂蛋白", Unit: "mmol/L", Reference: "1.9-3.6"},
				// {Name: "载脂蛋白A1", Unit: "g/L", Reference: "1.2-1.6"},
				// {Name: "载脂蛋白B", Unit: "g/L", Reference: "0.6-1.1"},
				// {Name: "脂蛋白(a)", Unit: "mg/L", Reference: "0-300"},
				// {Name: "同型半胱氨酸", Unit: "μmol/L", Reference: "5-15"},
			},
		},
		{
			Name: "血常规",
			Indicators: []HealthIndicator{
				{Name: "红细胞", Unit: "×10^12/L", Reference: "4.3-5.8"},
				{Name: "血红蛋白", Unit: "g/L", Reference: "130-175"},
				{Name: "红细胞压积", Unit: "%", Reference: "40-50"},
				// {Name: "白细胞", Unit: "×10^9/L", Reference: "3.5-9.5"},
				// {Name: "血小板", Unit: "×10^9/L", Reference: "125-350"},
				// {Name: "中性粒细胞百分比", Unit: "%", Reference: "50-70"},
				// {Name: "淋巴细胞百分比", Unit: "%", Reference: "20-50"},
				// {Name: "单核细胞百分比", Unit: "%", Reference: "3-10"},
				// {Name: "嗜酸性粒细胞百分比", Unit: "%", Reference: "0.4-8"},
				// {Name: "嗜碱性粒细胞百分比", Unit: "%", Reference: "0-1"},
				// {Name: "中性粒细胞", Unit: "×10^9/L", Reference: "1.8-6.3"},
				// {Name: "淋巴细胞", Unit: "×10^9/L", Reference: "1.1-3.2"},
				// {Name: "单核细胞", Unit: "×10^9/L", Reference: "0.1-0.6"},
				// {Name: "嗜酸性粒细胞", Unit: "×10^9/L", Reference: "0.02-0.52"},
				// {Name: "嗜碱性粒细胞", Unit: "×10^9/L", Reference: "0-0.06"},
				// {Name: "平均红细胞体积", Unit: "fL", Reference: "80-100"},
				// {Name: "平均红细胞血红蛋白量", Unit: "pg", Reference: "27-34"},
				// {Name: "平均红细胞血红蛋白浓度", Unit: "g/L", Reference: "320-360"},
				// {Name: "红细胞分布宽度", Unit: "%", Reference: "11.5-14.5"},
				// {Name: "血小板压积", Unit: "%", Reference: "0.1-0.3"},
				// {Name: "平均血小板体积", Unit: "fL", Reference: "7.5-11.5"},
				// {Name: "血小板分布宽度", Unit: "fL", Reference: "9.8-16.1"},
				// {Name: "大血小板比率", Unit: "%", Reference: "13-43"},
			},
		},
		{
			Name: "尿常规",
			Indicators: []HealthIndicator{
				{Name: "蛋白质", Unit: "g/L", Reference: "阴性或弱阳性"},
				{Name: "隐血", Unit: "HPF", Reference: "阴性"},
				{Name: "红细胞", Unit: "HPF", Reference: "0-3"},
				{Name: "非鳞状上皮细胞", Unit: "", Reference: "少量"},
				// {Name: "尿比重", Unit: "", Reference: "1.003-1.030"},
				// {Name: "尿pH", Unit: "", Reference: "4.5-8.0"},
				// {Name: "尿糖", Unit: "", Reference: "阴性"},
				// {Name: "尿酮体", Unit: "", Reference: "阴性"},
				// {Name: "尿白细胞", Unit: "", Reference: "阴性"},
				// {Name: "尿亚硝酸盐", Unit: "", Reference: "阴性"},
				// {Name: "尿胆原", Unit: "", Reference: "弱阳性"},
				// {Name: "胆红素", Unit: "", Reference: "阴性"},
				// {Name: "维生素C", Unit: "", Reference: "阴性"},
				// {Name: "白细胞", Unit: "/HPF", Reference: "0-5"},
				// {Name: "细菌", Unit: "/HPF", Reference: "少量"},
				// {Name: "管型", Unit: "/LPF", Reference: "0-1"},
				// {Name: "结晶", Unit: "/HPF", Reference: "少量或无"},
				// {Name: "粘液丝", Unit: "/LPF", Reference: "0-20"},
			},
		},
		{
			Name: "尿蛋白、尿素、肌酐测定",
			Indicators: []HealthIndicator{
				{Name: "尿蛋白", Unit: "mg/dL", Reference: "0-20"},
				{Name: "尿蛋白肌酐比值", Unit: "mg/g", Reference: "0-30"},
				// {Name: "微量白蛋白", Unit: "mg/L", Reference: "0-20"},
				// {Name: "尿免疫球蛋白G", Unit: "mg/L", Reference: "0-20"},
				// {Name: "尿α1微球蛋白", Unit: "mg/L", Reference: "0-12.5"},
				// {Name: "尿β2微球蛋白", Unit: "mg/L", Reference: "0-300"},
				// {Name: "尿NAG酶", Unit: "U/L", Reference: "0-15"},
				// {Name: "尿转铁蛋白", Unit: "mg/L", Reference: "0-2.8"},
				// {Name: "尿微量清蛋白/肌酐", Unit: "mg/g", Reference: "0-30"},
			},
		},
		{
			Name:       "其他",
			Indicators: []HealthIndicator{
				// {Name: "体重", Unit: "kg", Reference: ""},
				// {Name: "身高", Unit: "cm", Reference: ""},
				// {Name: "血压", Unit: "mmHg", Reference: "<140/90"},
				// {Name: "心率", Unit: "次/分钟", Reference: "60-100"},
				// {Name: "BMI", Unit: "kg/m²", Reference: "18.5-23.9"},
				// {Name: "腰围", Unit: "cm", Reference: "<90(男), <80(女)"},
				// {Name: "臀围", Unit: "cm", Reference: ""},
			},
		},
	}

	c.JSON(200, config)
}
