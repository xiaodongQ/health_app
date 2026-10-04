#!/bin/bash

# 健康指标记录批量添加脚本
# 用法: ./add_test_records.sh

DB_FILE="${DB_FILE:-health_app.db}"

# 检查 sqlite3
if ! command -v sqlite3 &> /dev/null; then
    echo "错误: 需要安装 sqlite3"
    exit 1
fi

# 检查数据库
if [ ! -f "$DB_FILE" ]; then
    echo "错误: 数据库文件 $DB_FILE 不存在"
    exit 1
fi

# 添加记录的函数
# 参数: user_id category record_date notes result1 result2 ...
# result 格式: "indicator|value|unit|reference|abnormal"
add_record() {
    local user_id=$1
    local category=$2
    local record_date=$3
    local notes=$4
    shift 4
    local results=("$@")

    # 插入记录
    local record_id=$(sqlite3 "$DB_FILE" "
        INSERT INTO health_records (user_id, category, record_date, notes)
        VALUES ($user_id, '$category', '$record_date', '$notes');
        SELECT last_insert_rowid();
    ")

    # 插入指标结果
    for result in "${results[@]}"; do
        IFS='|' read -r indicator value unit ref abnormal <<< "$result"
        sqlite3 "$DB_FILE" "
            INSERT INTO test_results (record_id, indicator, value, unit, reference, abnormal)
            VALUES ($record_id, '$indicator', $value, '$unit', '$ref', $abnormal);
        "
    done

    echo "添加记录 ID=$record_id 日期=$record_date"
}

echo "=== 添加健康记录 ==="
echo ""

# 2026-01-18 记录
add_record 1 "综合健康检查" "2026-01-18" "" \
    "肌酐|435|μmol/L|57-97|0" \
    "尿素|30|mmol/L|3.10-8.00|0" \
    "钾|3.29|mmol/L|3.5-5.3|0" \
    "空腹血糖|6.34|mmol/L|3.9-6.1|0" \
    "红细胞|3.88|×10^12/L|4.0-5.5|0" \
    "血红蛋白|119|g/L|130-175|0" \
    "红细胞压积|33.9|%|0.40-0.50|0" \
    "隐血|80|HPF|0-5|0" \
    "蛋白质|1.0|g/L|0-0.2|0" \
    "尿蛋白|1.6|mg/dL|0-0.15|0" \
    "尿蛋白肌酐比值|2.17|mg/g|0-30|0"

# 2026-02-15 记录
add_record 1 "综合健康检查" "2026-02-15" "" \
    "肌酐|412|μmol/L|57-97|0" \
    "尿素|28.5|mmol/L|3.10-8.00|0" \
    "钾|3.45|mmol/L|3.5-5.3|0" \
    "空腹血糖|5.98|mmol/L|3.9-6.1|0" \
    "甘油三酯|1.68|mmol/L|0-1.7|0" \
    "总胆固醇|4.82|mmol/L|2.8-5.2|0" \
    "红细胞|3.95|×10^12/L|4.0-5.5|0" \
    "血红蛋白|121|g/L|130-175|0" \
    "红细胞压积|34.8|%|0.40-0.50|0" \
    "隐血|60|HPF|0-5|0" \
    "蛋白质|0.80|g/L|0-0.2|0" \
    "尿蛋白|1.20|mg/dL|0-0.15|0" \
    "尿蛋白肌酐比值|1.85|mg/g|0-30|0"

# 2026-03-22 记录
add_record 1 "综合健康检查" "2026-03-22" "" \
    "肌酐|398|μmol/L|57-97|0" \
    "尿素|26.2|mmol/L|3.10-8.00|0" \
    "钾|3.52|mmol/L|3.5-5.3|0" \
    "空腹血糖|5.76|mmol/L|3.9-6.1|0" \
    "甘油三酯|1.55|mmol/L|0-1.7|0" \
    "总胆固醇|4.65|mmol/L|2.8-5.2|0" \
    "高密度脂蛋白|1.18|mmol/L|1.0-1.8|0" \
    "低密度脂蛋白|2.98|mmol/L|0-3.4|0" \
    "红细胞|4.01|×10^12/L|4.0-5.5|0" \
    "血红蛋白|124|g/L|130-175|0" \
    "红细胞压积|35.6|%|0.40-0.50|0" \
    "隐血|40|HPF|0-5|0" \
    "蛋白质|0.50|g/L|0-0.2|0" \
    "尿蛋白|0.92|mg/dL|0-0.15|0" \
    "尿蛋白肌酐比值|1.38|mg/g|0-30|0"

echo ""
echo "=== 当前记录列表 ==="
sqlite3 -header -column "$DB_FILE" "
    SELECT hr.id, hr.record_date, hr.category, COUNT(tr.id) as indicators
    FROM health_records hr
    LEFT JOIN test_results tr ON hr.id = tr.record_id
    WHERE hr.user_id = 1
    GROUP BY hr.id
    ORDER BY hr.record_date DESC;
"

echo ""
echo "添加完成!"
