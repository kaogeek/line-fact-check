#!/bin/bash

# Create a Message via the factcheck API
# Based on the Message struct and database schema

# Check if topic ID is provided
if [ $# -eq 0 ]; then
    echo "Error: Topic ID is required"
    echo "Usage: $0 <topic_id>"
    echo "Example: $0 550e8400-e29b-41d4-a716-446655440000"
    exit 1
fi

TOPIC_ID="$1"
API_URL="http://localhost:8080/messages"

# Current timestamp for created_at
CREATED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Creating a new message for topic: $TOPIC_ID"

# Create the JSON payload
JSON_PAYLOAD=$(cat << EOF
{
  "topic_id": "$TOPIC_ID",
  "text": "This is a test message for fact-checking",
  "type": "TYPE_TEXT"
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