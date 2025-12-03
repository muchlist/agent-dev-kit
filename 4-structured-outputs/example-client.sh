#!/bin/bash

# Example: Using ADK Email Agent as a Backend API
# Make sure the agent is running: make run/4

BASE_URL="http://localhost:8080"
APP_NAME="email_agent"
USER_ID="user"
SESSION_ID="session-$(date +%s)"

echo "=========================================="
echo "ADK API Backend Example"
echo "=========================================="
echo ""

# Step 1: List available agents
echo "1. Listing available agents..."
curl -s "${BASE_URL}/api/list-apps?relative_path=./" | jq '.'
echo ""

# Step 2: Create a new session
echo "2. Creating new session: ${SESSION_ID}..."
curl -s -X POST "${BASE_URL}/api/apps/${APP_NAME}/users/${USER_ID}/sessions" \
  -H "Content-Type: application/json" \
  -d "{\"session_id\": \"${SESSION_ID}\"}" | jq '.'
echo ""

# Step 3: Send a message to generate an email
echo "3. Generating email..."
RESPONSE=$(curl -s -X POST "${BASE_URL}/api/run" \
  -H "Content-Type: application/json" \
  -d "{
    \"app\": \"${APP_NAME}\",
    \"user_id\": \"${USER_ID}\",
    \"session_id\": \"${SESSION_ID}\",
    \"user_message\": \"Write a professional email to my team about the upcoming project deadline that has been extended by two weeks. The new deadline is December 20th.\"
  }")

echo ""
echo "4. Response received!"
echo ""
echo "Full Response:"
echo "${RESPONSE}" | jq '.'
echo ""

# Step 4: Extract structured output
echo "=========================================="
echo "Structured Output (from state.email):"
echo "=========================================="
echo "${RESPONSE}" | jq '.state.email'
echo ""

# Extract subject and body separately
SUBJECT=$(echo "${RESPONSE}" | jq -r '.state.email.subject')
BODY=$(echo "${RESPONSE}" | jq -r '.state.email.body')

echo "Subject: ${SUBJECT}"
echo ""
echo "Body:"
echo "${BODY}"
echo ""

# Step 5: Get session history
echo "=========================================="
echo "5. Session History:"
echo "=========================================="
curl -s "${BASE_URL}/api/apps/${APP_NAME}/users/${USER_ID}/sessions/${SESSION_ID}" | jq '.messages'
echo ""

# Step 6: List all sessions for this user
echo "=========================================="
echo "6. All user sessions:"
echo "=========================================="
curl -s "${BASE_URL}/api/apps/${APP_NAME}/users/${USER_ID}/sessions" | jq '.'
echo ""

echo "=========================================="
echo "Example completed!"
echo "=========================================="
