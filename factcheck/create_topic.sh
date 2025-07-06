#!/bin/bash

# Create a Topic via the factcheck API
# Based on the Topic struct and database schema

API_URL="http://localhost:8080/topics"

# Generate a UUID for the topic ID
TOPIC_ID=$(uuidgen)

# Current timestamp for created_at
CREATED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Creating a new topic..."

# Create the JSON payload
JSON_PAYLOAD=$(cat << EOF
{
  "id": "$TOPIC_ID",
  "name": "Test Fact-Check Topic",
  "result": "",
  "result_status": "TOPIC_RESULT_NONE",
  "created_at": "$CREATED_AT",
  "updated_at": null
}
EOF
)

echo "Request payload:"
echo "$JSON_PAYLOAD"
echo ""

# Make the API call
echo "Making POST request to $API_URL..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "$JSON_PAYLOAD" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s

echo ""
echo "Done!"