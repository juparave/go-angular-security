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
# Intervals: month, week, year or day
# The product name can contain spaces
declare -a PRODUCTS=(
  "App Basic 25 month 229.00"
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

  echo "Checking for existing product: '$product_name'..." >&2 # Redirect to stderr
  # Use try-catch block for stripe command
  local list_response list_exit_code
  if ! list_response=$(stripe products list --limit 100 2>&1); then
      list_exit_code=$?
      echo "Error: Failed to list products (Exit Code: $list_exit_code)." >&2
      echo "Stripe CLI output:" >&2
      echo "$list_response" >&2
      exit 1
  fi

  # Use jq to find the product ID
  product_id=$(echo "$list_response" | jq -r ".data[] | select(.name==\"$product_name\") | .id" | head -n 1)

  if [ -n "$product_id" ]; then
    echo "Product '$product_name' already exists with ID: $product_id" >&2 # Redirect to stderr
  else
    echo "Creating product '$product_name'..." >&2 # Redirect to stderr
    local create_response create_exit_code
    # Important: Keep 2>&1 here to capture potential errors from stripe command itself
    if ! create_response=$(stripe products create \
      --name="$product_name" \
      --description="$product_description" 2>&1); then
        create_exit_code=$?
        echo "Error: Failed to create product '$product_name' (Exit Code: $create_exit_code)." >&2
        echo "Stripe CLI output:" >&2
        echo "$create_response" >&2
        exit 1 # Exit script on failure
    fi

    # Check if the response contains a valid '.id' field
    if ! echo "$create_response" | jq -e '.id' > /dev/null; then
        echo "Error: Failed to create product '$product_name'. Invalid response from Stripe." >&2
        echo "Stripe CLI output:" >&2
        echo "$create_response" >&2
        exit 1 # Exit script on failure
    fi

    product_id=$(echo "$create_response" | jq -r '.id')
    echo "Product '$product_name' created with ID: $product_id" >&2 # Redirect to stderr
  fi

  # Return the product ID - THIS is the only echo that should go to stdout
  echo "$product_id"
}
# --- End Helper Functions ---

# --- Main Script ---

# Check dependencies
# ... (dependency checks remain the same) ...
if ! command -v jq &> /dev/null; then echo "Error: 'jq' command not found." >&2; exit 1; fi
if ! command -v bc &> /dev/null; then echo "Error: 'bc' command not found." >&2; exit 1; fi
if ! command -v stripe &> /dev/null; then echo "Error: 'stripe' command not found." >&2; exit 1; fi


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
      echo "Error: Failed to obtain Product ID for '$formatted_name'. Exiting." >&2
      exit 1
  fi
  echo "Using Product ID: $product_id"

  # --- Check for Existing Price ---
  echo "Checking for existing active price for Product ID $product_id..." >&2
  price_id="" # Initialize price_id
  # Use temporary variable for the command output to check exit status with pipefail
  list_prices_output=""
  list_prices_exit_code=0
  # Run the command and capture output; check exit code using pipefail behavior
  list_prices_output=$(stripe prices list --product "$product_id" --active=true --limit 100 | \
                       jq -r --argjson amount "$price_cents" \
                             --arg currency "$CURRENCY" \
                             --arg interval "$interval" \
                             '.data[] | select(.unit_amount == $amount and .currency == $currency and .recurring.interval == $interval) | .id' | \
                       head -n 1) || list_prices_exit_code=$? # Capture exit code on failure

  if [ $list_prices_exit_code -ne 0 ]; then
      # jq returns non-zero if selection fails or input is bad, stripe might too
      echo "Warning: Failed to query/filter existing prices (Exit Code: $list_prices_exit_code). Assuming price needs creation." >&2
      # Potentially inspect list_prices_output if it contains error messages
      price_id="" # Ensure price_id is empty
  else
      price_id="$list_prices_output" # Assign output if command succeeded
  fi

  # --- Create Price IF it doesn't exist ---
  if [ -z "$price_id" ]; then
      echo "No matching active price found. Creating new price..." >&2
      price_response=""
      stripe_exit_code=0
      # Capture stdout and stderr, check exit code
      if ! price_response=$(stripe prices create \
        -d product="$product_id" \
        -d unit_amount="$price_cents" \
        -d currency="$CURRENCY" \
        -d "recurring[interval]"="$interval" 2>&1); then
          stripe_exit_code=$?
          echo "Error: Failed to create price for '$formatted_name' (Exit Code: $stripe_exit_code)." >&2
          echo "Stripe CLI output:" >&2
          echo "$price_response" >&2
          echo "Skipping price creation for '$formatted_name' due to error." >&2
          continue # Skip to the next product in the loop
      fi

      # Validate response contains ID
      if ! echo "$price_response" | jq -e '.id' > /dev/null; then
          echo "Error: Failed to validate price creation response for '$formatted_name'." >&2
          echo "Stripe CLI output:" >&2
          echo "$price_response" >&2
          echo "Skipping price creation for '$formatted_name' due to invalid response." >&2
          continue # Skip to the next product
      fi

      # Extract the ID
      price_id=$(echo "$price_response" | jq -r '.id')

      # Final check (belt and suspenders)
      if [ -z "$price_id" ] || [ "$price_id" = "null" ]; then
          echo "Error: Extracted Price ID is empty or null for '$formatted_name' even after validation." >&2
          echo "Stripe CLI output:" >&2
          echo "$price_response" >&2
          echo "Skipping price creation for '$formatted_name' due to invalid ID." >&2
          continue # Skip to the next product
      fi
      echo "Successfully created new price: $price_id" >&2 # Log creation

  else
      echo "Found existing active price: $price_id" >&2 # Log finding
  fi

  # --- Generate Env Var and Save ---
  env_var_name="STRIPE_PRICE_$(echo "${product_name}_${interval}" | tr ' ' '_' | tr '[:lower:]' '[:upper:]')"
  echo "$env_var_name=$price_id" >> "$ENV_FILE"

  echo "Using Price ID: $price_id for $formatted_name ($interval)"
  echo "Saved as: $env_var_name"

done

echo "--------------------------------------------------"
echo "Script finished."
echo "Price IDs have been written to $ENV_FILE (currency: $CURRENCY):"
cat "$ENV_FILE"
echo "--------------------------------------------------"

