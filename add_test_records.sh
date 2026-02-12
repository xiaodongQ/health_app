#!/bin/bash

# 健康指标记录批量添加脚本
# 用于向健康指标管理系统添加指定的记录

# API基础URL
API_BASE_URL="http://localhost:8081/api"

# 用户ID（您可以根据需要修改）
USER_ID=1

# 添加健康记录的函数
add_health_record() {
    local date=$1
    shift
    local json_data=$1
    
    echo "正在添加 $date 的记录..."
    
    response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -d "$json_data" \
        "$API_BASE_URL/health_records")
    
    # 分离响应体和HTTP状态码 (兼容macOS的BSD版head命令)
    body=$(echo "$response" | sed '$d')
    status_code=$(echo "$response" | tail -n 1)
    
    if [ "$status_code" -eq 201 ] || [ "$status_code" -eq 200 ]; then
        echo "✓ $date 记录添加成功"
    else
        echo "✗ $date 记录添加失败，状态码: $status_code"
        echo "错误信息: $body"
    fi
    echo ""
}

# 添加测试记录
# 添加2025年6月8日的记录
echo "开始添加2025年6月8日的记录..."

json_20250608='{
  "user_id": '$USER_ID',
  "record_date": "2025-06-08T00:00:00Z",
  "category": "综合健康检查",
  "results": [
    {"indicator": "肌酐", "value": 95, "unit": "μmol/L"},
    {"indicator": "尿素", "value": 4.5, "unit": "mmol/L"},
    {"indicator": "钾", "value": 4.2, "unit": "mmol/L"},
    {"indicator": "空腹血糖", "value": 5.2, "unit": "mmol/L"},
    {"indicator": "红细胞", "value": 4.5, "unit": "×10^12/L"},
    {"indicator": "血红蛋白", "value": 140, "unit": "g/L"},
    {"indicator": "红细胞压积", "value": 42, "unit": "%"},
    {"indicator": "隐血", "value": 0, "unit": "HPF"},
    {"indicator": "蛋白质", "value": 0.1, "unit": "g/L"},
    {"indicator": "红细胞", "value": 2, "unit": "HPF"},
    {"indicator": "非鳞状上皮细胞", "value": 0.5, "unit": ""},
    {"indicator": "尿蛋白", "value": 5, "unit": "mg/dL"},
    {"indicator": "尿蛋白肌酐比值", "value": 0.5, "unit": "mg/g"}
  ]
}'

add_health_record "2025-06-08" "$json_20250608"

# 添加2025年7月8日的记录
echo "开始添加2025年7月8日的记录..."

json_20250708='{
  "user_id": '$USER_ID',
  "record_date": "2025-07-08T00:00:00Z",
  "category": "综合健康检查",
  "results": [
    {"indicator": "肌酐", "value": 98, "unit": "μmol/L"},
    {"indicator": "尿素", "value": 4.2, "unit": "mmol/L"},
    {"indicator": "钾", "value": 4.1, "unit": "mmol/L"},
    {"indicator": "空腹血糖", "value": 5.0, "unit": "mmol/L"},
    {"indicator": "红细胞", "value": 4.6, "unit": "×10^12/L"},
    {"indicator": "血红蛋白", "value": 142, "unit": "g/L"},
    {"indicator": "红细胞压积", "value": 43, "unit": "%"},
    {"indicator": "隐血", "value": 0, "unit": "HPF"},
    {"indicator": "蛋白质", "value": 0.1, "unit": "g/L"},
    {"indicator": "红细胞", "value": 1, "unit": "HPF"},
    {"indicator": "非鳞状上皮细胞", "value": 0.3, "unit": ""},
    {"indicator": "尿蛋白", "value": 6, "unit": "mg/dL"},
    {"indicator": "尿蛋白肌酐比值", "value": 0.6, "unit": "mg/g"}
  ]
}'

add_health_record "2025-07-08" "$json_20250708"

# 添加2025年8月8日的记录
echo "开始添加2025年8月8日的记录..."

json_20250808='{
  "user_id": '$USER_ID',
  "record_date": "2025-08-08T00:00:00Z",
  "category": "综合健康检查",
  "results": [
    {"indicator": "肌酐", "value": 102, "unit": "μmol/L"},
    {"indicator": "尿素", "value": 4.8, "unit": "mmol/L"},
    {"indicator": "钾", "value": 4.0, "unit": "mmol/L"},
    {"indicator": "空腹血糖", "value": 5.1, "unit": "mmol/L"},
    {"indicator": "红细胞", "value": 4.4, "unit": "×10^12/L"},
    {"indicator": "血红蛋白", "value": 138, "unit": "g/L"},
    {"indicator": "红细胞压积", "value": 41, "unit": "%"},
    {"indicator": "隐血", "value": 0, "unit": "HPF"},
    {"indicator": "蛋白质", "value": 0.1, "unit": "g/L"},
    {"indicator": "红细胞", "value": 2, "unit": "HPF"},
    {"indicator": "非鳞状上皮细胞", "value": 0.4, "unit": ""},
    {"indicator": "尿蛋白", "value": 4, "unit": "mg/dL"},
    {"indicator": "尿蛋白肌酐比值", "value": 0.4, "unit": "mg/g"}
  ]
}'

add_health_record "2025-08-08" "$json_20250808"

# 添加2025年9月8日的记录
echo "开始添加2025年9月8日的记录..."

json_20250908='{
  "user_id": '$USER_ID',
  "record_date": "2025-09-08T00:00:00Z",
  "category": "综合健康检查",
  "results": [
    {"indicator": "肌酐", "value": 99, "unit": "μmol/L"},
    {"indicator": "尿素", "value": 4.6, "unit": "mmol/L"},
    {"indicator": "钾", "value": 4.3, "unit": "mmol/L"},
    {"indicator": "空腹血糖", "value": 5.3, "unit": "mmol/L"},
    {"indicator": "红细胞", "value": 4.5, "unit": "×10^12/L"},
    {"indicator": "血红蛋白", "value": 141, "unit": "g/L"},
    {"indicator": "红细胞压积", "value": 42, "unit": "%"},
    {"indicator": "隐血", "value": 0, "unit": "HPF"},
    {"indicator": "蛋白质", "value": 0.1, "unit": "g/L"},
    {"indicator": "红细胞", "value": 1, "unit": "HPF"},
    {"indicator": "非鳞状上皮细胞", "value": 0.2, "unit": ""},
    {"indicator": "尿蛋白", "value": 5, "unit": "mg/dL"},
    {"indicator": "尿蛋白肌酐比值", "value": 0.5, "unit": "mg/g"}
  ]
}'

add_health_record "2025-09-08" "$json_20250908"

# 添加2025年10月8日的记录
echo "开始添加2025年10月8日的记录..."

json_20251008='{
  "user_id": '$USER_ID',
  "record_date": "2025-10-08T00:00:00Z",
  "category": "综合健康检查",
  "results": [
    {"indicator": "肌酐", "value": 105, "unit": "μmol/L"},
    {"indicator": "尿素", "value": 4.9, "unit": "mmol/L"},
    {"indicator": "钾", "value": 4.2, "unit": "mmol/L"},
    {"indicator": "空腹血糖", "value": 5.0, "unit": "mmol/L"},
    {"indicator": "红细胞", "value": 4.3, "unit": "×10^12/L"},
    {"indicator": "血红蛋白", "value": 139, "unit": "g/L"},
    {"indicator": "红细胞压积", "value": 40, "unit": "%"},
    {"indicator": "隐血", "value": 0, "unit": "HPF"},
    {"indicator": "蛋白质", "value": 0.1, "unit": "g/L"},
    {"indicator": "红细胞", "value": 1, "unit": "HPF"},
    {"indicator": "非鳞状上皮细胞", "value": 0.3, "unit": ""},
    {"indicator": "尿蛋白", "value": 6, "unit": "mg/dL"},
    {"indicator": "尿蛋白肌酐比值", "value": 0.6, "unit": "mg/g"}
  ]
}'

add_health_record "2025-10-08" "$json_20251008"

echo "所有记录添加完成！"
