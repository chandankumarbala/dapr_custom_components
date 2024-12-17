#!/bin/bash

# Define the URL and the JSON payload
  # Replace with your actual JSON payload

for i in {1..10}; do

    URL="http://localhost:3000/send-data"  # Replace with your API endpoint
    JSON_PAYLOAD='{"hello":"hi"}'
    response=$(curl -s -X POST "$URL" -H "Content-Type: application/json" -d "$JSON_PAYLOAD")
    extracted_value=$(echo "$response" | jq -r '.timetaken')

    URL2="http://localhost:3000/trigger-binding"  # Replace with your API endpoint
    JSON_PAYLOAD='{"hello":"hi"}'
    response1=$(curl -s -X POST "$URL2" -H "Content-Type: application/json" -d "$JSON_PAYLOAD")
    extracted_value2=$(echo "$response1" | jq -r '.timetaken')

# Check if jq extraction was successful

        echo "Without binding: ${extracted_value} / with binding: ${extracted_value2}"
        sleep 5
done