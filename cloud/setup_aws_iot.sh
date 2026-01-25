#!/bin/bash

# AWS IoT setup script
THING_NAME="coffee-maker-001"
POLICY_NAME="CoffeeMakerPolicy"

echo "🔧 Setting up AWS IoT..."

# Create IoT thing
aws iot create-thing --thing-name $THING_NAME

# Create policy
aws iot create-policy \
  --policy-name $POLICY_NAME \
  --policy-document file://iot-policy.json

# Create certificate
CERT_ARN=$(aws iot create-keys-and-certificate \
  --set-as-active \
  --certificate-pem-outfile cert.pem \
  --public-key-outfile public.key \
  --private-key-outfile private.key \
  --query certificateArn --output text)

# Attach policy to certificate
aws iot attach-policy \
  --policy-name $POLICY_NAME \
  --target $CERT_ARN

# Attach certificate to thing
aws iot attach-thing-principal \
  --thing-name $THING_NAME \
  --principal $CERT_ARN

# Get IoT endpoint
ENDPOINT=$(aws iot describe-endpoint --endpoint-type iot:Data-ATS --query endpointAddress --output text)

echo "✅ AWS IoT setup complete!"
echo "   Endpoint: $ENDPOINT"
echo "   Certificates: cert.pem, private.key"
