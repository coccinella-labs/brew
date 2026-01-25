#!/bin/bash

echo "🔧 Building firmware..."
make clean && make

echo "☁️ Deploying AWS infrastructure..."
aws cloudformation deploy \
  --template-file infrastructure.yaml \
  --stack-name coffee-maker-stack \
  --capabilities CAPABILITY_IAM

echo "🔌 Setting up IoT certificates..."
./setup_aws_iot.sh

echo "📱 Starting services..."
go run web_dashboard.go &
WEB_PID=$!

go run iot_controller.go &
IOT_PID=$!

go run cloud_analytics.go &
ANALYTICS_PID=$!

go run serial_controller.go &
SERIAL_PID=$!

echo "✅ Complete IoT coffee maker system running!"
echo "   Dashboard: http://localhost:8080"
echo "   IoT: Connected to AWS"
echo "   Analytics: Logging to DynamoDB"
echo "   Press Ctrl+C to stop"

# Cleanup on exit
trap "kill $WEB_PID $IOT_PID $ANALYTICS_PID $SERIAL_PID 2>/dev/null" EXIT
wait
