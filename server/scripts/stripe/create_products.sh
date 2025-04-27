#!/usr/bin/env bash
# Prerequisite: You need jq installed.

# Exit immediately if a command exits with a non-zero status.
set -e
# Treat unset variables as an error when substituting.
set -u
# Pipeliines fail if any command fails, not just the last one.
set -o pipefail

echo "Creating Premium product..."
# Create Premium product and capture its ID
# The 'jq -r .id' command parses the JSON output and extracts the raw string value of the 'id' field.
PREMIUM_PRODUCT_ID=$(stripe products create \
  --name="Billing Guide: Premium Service" \
  --description="Premium service with extra features" | jq -r '.id')

# Check if PREMIUM_PRODUCT_ID was captured
if [ -z "$PREMIUM_PRODUCT_ID" ]; then
  echo "Error: Failed to capture Premium Product ID."
  exit 1
fi
echo "Premium Product created with ID: $PREMIUM_PRODUCT_ID"

echo "Creating Basic product..."
# Create Basic product and capture its ID
BASIC_PRODUCT_ID=$(stripe products create \
  --name="Billing Guide: Basic Service" \
  --description="Basic service with minimum features" | jq -r '.id')

# Check if BASIC_PRODUCT_ID was captured
if [ -z "$BASIC_PRODUCT_ID" ]; then
  echo "Error: Failed to capture Basic Product ID."
  exit 1
fi
echo "Basic Product created with ID: $BASIC_PRODUCT_ID"

echo "Creating Premium price for Product ID: $PREMIUM_PRODUCT_ID..."
# Create Premium price using the captured PREMIUM_PRODUCT_ID
stripe prices create \
  -d product="$PREMIUM_PRODUCT_ID" \
  -d unit_amount=1500 \
  -d currency=usd \
  -d "recurring[interval]"=month

echo "Premium Price created."

echo "Creating Basic price for Product ID: $BASIC_PRODUCT_ID..."
# Create Basic price using the captured BASIC_PRODUCT_ID
stripe prices create \
  -d product="$BASIC_PRODUCT_ID" \
  -d unit_amount=500 \
  -d currency=usd \
  -d "recurring[interval]"=month

echo "Basic Price created."

echo "Script finished successfully."
