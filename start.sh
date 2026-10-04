#!/bin/bash
set -e

APP_DIR="/home/workspace/go_path/health_app"
APP_NAME="health_app"
PORT=8081

cd "$APP_DIR"

# 检查是否已在运行
if curl -s "http://localhost:$PORT/api/health" > /dev/null 2>&1; then
    echo "$APP_NAME is already running on port $PORT"
    exit 0
fi

# 启动应用
nohup ./health_app > /var/log/health_app.log 2>&1 &

echo "$APP_NAME started, PID: $!"
echo "Log: /var/log/health_app.log"
