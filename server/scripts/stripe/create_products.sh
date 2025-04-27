#!/usr/bin/env bash
# Prerequisite: You need jq and bc installed.

# Exit immediately if a command exits with a non-zero status.
set -e
# Treat unset variables as an error when substituting.
set -u
# Pipeliines fail if any command fails, not just the last one.
set -o pipefail

# --- Configuration ---
# Currency to use for all products (change as needed)
CURRENCY="mxn"

# Define product data
# Format: Each line contains: "Product full name" interval price
# The product name can contain spaces
declare -a PRODUCTS=(
  "App Basic 25 month 230.00"
  "App Pro 100 month 499.00"
  "App Basic 25 year 2300.00"
  "App Pro 100 year 4990.00"
)

# Environment variable output file
ENV_FILE="stripe_prices.env"
# --- End Configuration ---

# --- Helper Functions ---

# Function to create product if it doesn't exist and return its ID
create_or_get_product() {
  local product_name="$1"
  local product_description="$2"
  local product_id="" # Initialize variable

  # Check if product exists using Stripe CLI's filtering if possible, or jq
  # Note: Stripe CLI list filtering might be limited. jq is more reliable here.
  echo "Checking for existing product: '$product_name'..." >&2 # Redirect to stderr
  product_id=$(stripe products list --limit 100 | jq -r ".data[] | select(.name==\"$product_name\") | .id" | head -n 1) # head -n 1 ensures only one ID if duplicates somehow exist

  if [ -n "$product_id" ]; then
    echo "Product '$product_name' already exists with ID: $product_id" >&2 # Redirect to stderr
  else
    echo "Creating product '$product_name'..." >&2 # Redirect to stderr
    # Capture response and check for errors
    local create_response
    # Important: Keep 2>&1 here to capture potential errors from stripe command itself
    create_response=$(stripe products create \
      --name="$product_name" \
      --description="$product_description" 2>&1)
    local exit_code=$?

    # Check exit code AND if the response contains a valid '.id' field
    if [ $exit_code -ne 0 ] || ! echo "$create_response" | jq -e '.id' > /dev/null; then
        # These error messages should also go to stderr
        echo "Error: Failed to create product '$product_name'." >&2
        echo "Stripe CLI output (Exit Code: $exit_code):" >&2
        echo "$create_response" >&2
        exit 1 # Exit script on failure
    fi

    product_id=$(echo "$create_response" | jq -r '.id')
    echo "Product '$product_name' created with ID: $product_id" >&2 # Redirect to stderr
  fi

  # Return the product ID - THIS is the only echo that should go to stdout
  echo "$product_id"
}

# --- Main Script ---

# Check dependencies
if ! command -v jq &> /dev/null; then
    echo "Error: 'jq' command not found. Please install it (e.g., brew install jq, sudo apt install jq)." >&2
    exit 1
fi
if ! command -v bc &> /dev/null; then
    echo "Error: 'bc' command not found. Please install it (e.g., sudo apt install bc, sudo yum install bc)." >&2
    exit 1
fi
if ! command -v stripe &> /dev/null; then
    echo "Error: 'stripe' command not found. Please install the Stripe CLI and ensure it's configured." >&2
    exit 1
fi

# Create or clear the .env file
echo "# Stripe Price IDs (Currency: $CURRENCY)" > "$ENV_FILE"
echo "# Generated on $(date)" >> "$ENV_FILE"

# Process each product
for product_info in "${PRODUCTS[@]}"; do
  # --- Parse product info ---
  price="${product_info##* }"
  temp="${product_info% *}"
  interval="${temp##* }"
  product_name="${temp% *}"
  # --- End Parsing ---

  formatted_name="$product_name"
  formatted_description="$product_name with ${interval} billing"
  price_cents=$(echo "$price * 100" | bc | sed 's/\..*$//') # Convert to integer cents

  echo "--------------------------------------------------"
  echo "Processing: $formatted_name ($interval) at $price $CURRENCY (Price in cents: $price_cents)"

  # --- Create or Get Product ---
  product_id=$(create_or_get_product "$formatted_name" "$formatted_description")
  if [ -z "$product_id" ]; then
      # The function create_or_get_product should exit on failure, but double-check
      echo "Error: Failed to obtain Product ID for '$formatted_name'. Exiting." >&2
      exit 1
  fi
  echo "Using Product ID: $product_id"

  # --- Create Price ---
  echo "Creating price for Product ID: $product_id..."
  price_response=$(stripe prices create \
    -d product="$product_id" \
    -d unit_amount="$price_cents" \
    -d currency="$CURRENCY" \
    -d "recurring[interval]"="$interval" 2>&1) # Capture stdout and stderr
  stripe_exit_code=$? # Capture exit code immediately

  # --- Validate Price Creation and Extract ID ---
  # Check exit code AND if the response contains a valid '.id' field using jq -e
  # jq -e exits with 0 if the filter produces a value other than null or false
  if [ $stripe_exit_code -ne 0 ] || ! echo "$price_response" | jq -e '.id' > /dev/null; then
      echo "Error: Failed to create or validate price for '$formatted_name'." >&2
      echo "Stripe CLI output (Exit Code: $stripe_exit_code):" >&2
      echo "$price_response" >&2
      # Decide whether to exit or continue with the next product
      # exit 1 # Option 1: Stop the whole script
      echo "Skipping price creation for '$formatted_name' due to error." >&2
      continue # Option 2: Skip to the next product in the loop
  fi

  # If validation passed, extract the ID
  price_id=$(echo "$price_response" | jq -r '.id')

  # Final check (belt and suspenders) - is the extracted ID non-empty and not "null"?
   if [ -z "$price_id" ] || [ "$price_id" = "null" ]; then
       echo "Error: Extracted Price ID is empty or null for '$formatted_name' even after validation." >&2
       echo "Stripe CLI output:" >&2
       echo "$price_response" >&2
       # exit 1 # Option 1: Stop the whole script
       echo "Skipping price creation for '$formatted_name' due to invalid ID." >&2
       continue # Option 2: Skip to the next product
   fi

  # --- Generate Env Var and Save ---
  env_var_name="STRIPE_PRICE_$(echo "${product_name}_${interval}" | tr ' ' '_' | tr '[:lower:]' '[:upper:]')"
  echo "$env_var_name=$price_id" >> "$ENV_FILE"

  echo "Successfully created price for $formatted_name ($interval): $price_id"
  echo "Saved as: $env_var_name"

done

echo "--------------------------------------------------"
echo "Script finished."
echo "Price IDs have been written to $ENV_FILE (currency: $CURRENCY):"
cat "$ENV_FILE"
echo "--------------------------------------------------"
