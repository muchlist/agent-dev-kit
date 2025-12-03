# ADK Generated API Reference

When you run an ADK agent with `api` mode, the launcher automatically generates a complete REST API. Here's the full API reference:

## Base URL
```
http://localhost:8080
```

## Core Agent Endpoints

### 1. List Available Agents
```http
GET /api/list-apps?relative_path=./
```
**Response:**
```json
{
  "apps": [
    {
      "name": "email_agent",
      "description": "Generates professional emails with structured subject and body"
    }
  ]
}
```

### 2. Run Agent (Server-Sent Events - Streaming)
```http
POST /api/run_sse
Content-Type: application/json
```
**Request Body:**
```json
{
  "app": "email_agent",
  "user_id": "user",
  "session_id": "session-123",
  "user_message": "Write an email about project deadline extension"
}
```
**Response:** Server-Sent Events stream with agent responses

### 3. Run Agent (Non-Streaming)
```http
POST /api/run
Content-Type: application/json
```
**Request Body:**
```json
{
  "app": "email_agent",
  "user_id": "user",
  "session_id": "session-123",
  "user_message": "Write an email about project deadline extension"
}
```
**Response:**
```json
{
  "response": "...",
  "session_id": "session-123",
  "state": {
    "email": {
      "subject": "Project Deadline Extension",
      "body": "Dear Team,\n\nI hope this message finds you well..."
    }
  }
}
```

## Session Management Endpoints

### 4. Create New Session
```http
POST /api/apps/{app_name}/users/{user_id}/sessions
Content-Type: application/json
```
**Request Body:**
```json
{
  "session_id": "optional-custom-id"
}
```
**Response:**
```json
{
  "session_id": "generated-or-custom-id",
  "created_at": "2025-12-03T16:32:01Z"
}
```

### 5. List User Sessions
```http
GET /api/apps/{app_name}/users/{user_id}/sessions
```
**Response:**
```json
{
  "sessions": [
    {
      "session_id": "session-123",
      "created_at": "2025-12-03T16:32:01Z",
      "last_activity": "2025-12-03T16:35:22Z"
    }
  ]
}
```

### 6. Get Session Details
```http
GET /api/apps/{app_name}/users/{user_id}/sessions/{session_id}
```
**Response:**
```json
{
  "session_id": "session-123",
  "messages": [
    {
      "role": "user",
      "content": "Write an email..."
    },
    {
      "role": "agent",
      "content": "{\"subject\": \"...\", \"body\": \"...\"}"
    }
  ],
  "state": {
    "email": {
      "subject": "Project Deadline Extension",
      "body": "Dear Team,..."
    }
  }
}
```

### 7. Delete Session
```http
DELETE /api/apps/{app_name}/users/{user_id}/sessions/{session_id}
```
**Response:**
```json
{
  "success": true
}
```

## Debugging Endpoints

### 8. Get Session Trace
```http
GET /api/debug/trace/session/{session_id}
```
**Response:** Detailed execution trace of the session

### 9. Health Check
```http
GET /api/health
```
**Response:**
```json
{
  "status": "healthy"
}
```

## Evaluation Endpoints (For Testing)

### 10. List Evaluation Sets
```http
GET /api/apps/{app_name}/eval_sets
```

### 11. Get Evaluation Results
```http
GET /api/apps/{app_name}/eval_results
```

## Complete cURL Examples

### Example 1: Simple Interaction
```bash
# Create a session
curl -X POST http://localhost:8080/api/apps/email_agent/users/user/sessions \
  -H "Content-Type: application/json" \
  -d '{"session_id": "my-session-1"}'

# Send a message
curl -X POST http://localhost:8080/api/run \
  -H "Content-Type: application/json" \
  -d '{
    "app": "email_agent",
    "user_id": "user",
    "session_id": "my-session-1",
    "user_message": "Write a professional email to my team about upcoming project deadline extension by two weeks"
  }'

# Get session history
curl http://localhost:8080/api/apps/email_agent/users/user/sessions/my-session-1
```

### Example 2: Streaming Response (SSE)
```bash
curl -X POST http://localhost:8080/api/run_sse \
  -H "Content-Type: application/json" \
  -d '{
    "app": "email_agent",
    "user_id": "user",
    "session_id": "my-session-2",
    "user_message": "Draft an email to client requesting additional information"
  }'
```

### Example 3: Using with Structured Outputs
```bash
# Send request
RESPONSE=$(curl -s -X POST http://localhost:8080/api/run \
  -H "Content-Type: application/json" \
  -d '{
    "app": "email_agent",
    "user_id": "user",
    "session_id": "my-session-3",
    "user_message": "Create a meeting invitation email for marketing department"
  }')

# Extract structured output from state
echo $RESPONSE | jq '.state.email'
# Output:
# {
#   "subject": "Meeting Invitation: Marketing Strategy Discussion",
#   "body": "Dear Marketing Team,\n\n..."
# }
```

## Using as Backend API

### Node.js/JavaScript Example
```javascript
async function generateEmail(prompt) {
  const response = await fetch('http://localhost:8080/api/run', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      app: 'email_agent',
      user_id: 'user',
      session_id: `session-${Date.now()}`,
      user_message: prompt
    })
  });

  const data = await response.json();

  // Access structured output
  const emailContent = data.state.email;
  console.log('Subject:', emailContent.subject);
  console.log('Body:', emailContent.body);

  return emailContent;
}

// Usage
const email = await generateEmail(
  'Write an email about project status update'
);
```

### Python Example
```python
import requests
import json

def generate_email(prompt):
    response = requests.post(
        'http://localhost:8080/api/run',
        json={
            'app': 'email_agent',
            'user_id': 'user',
            'session_id': f'session-{time.time()}',
            'user_message': prompt
        }
    )

    data = response.json()

    # Access structured output
    email_content = data['state']['email']
    print(f"Subject: {email_content['subject']}")
    print(f"Body: {email_content['body']}")

    return email_content

# Usage
email = generate_email('Write an email about project status update')
```

### Go Example
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type RunRequest struct {
    App         string `json:"app"`
    UserID      string `json:"user_id"`
    SessionID   string `json:"session_id"`
    UserMessage string `json:"user_message"`
}

type EmailContent struct {
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

type RunResponse struct {
    Response  string                 `json:"response"`
    SessionID string                 `json:"session_id"`
    State     map[string]interface{} `json:"state"`
}

func generateEmail(prompt string) (*EmailContent, error) {
    req := RunRequest{
        App:         "email_agent",
        UserID:      "user",
        SessionID:   fmt.Sprintf("session-%d", time.Now().Unix()),
        UserMessage: prompt,
    }

    body, _ := json.Marshal(req)

    resp, err := http.Post(
        "http://localhost:8080/api/run",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result RunResponse
    json.NewDecoder(resp.Body).Decode(&result)

    // Extract structured output
    emailData := result.State["email"].(map[string]interface{})
    email := &EmailContent{
        Subject: emailData["subject"].(string),
        Body:    emailData["body"].(string),
    }

    fmt.Printf("Subject: %s\n", email.Subject)
    fmt.Printf("Body: %s\n", email.Body)

    return email, nil
}

func main() {
    email, _ := generateEmail("Write an email about project status update")
    fmt.Println(email)
}
```

### Streaming with EventSource (Browser)
```javascript
function generateEmailStreaming(prompt) {
  const eventSource = new EventSource('/api/run_sse', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      app: 'email_agent',
      user_id: 'user',
      session_id: `session-${Date.now()}`,
      user_message: prompt
    })
  });

  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Chunk:', data);

    if (data.done) {
      // Access final structured output
      const email = data.state.email;
      console.log('Subject:', email.subject);
      console.log('Body:', email.body);
      eventSource.close();
    }
  };

  eventSource.onerror = (error) => {
    console.error('Error:', error);
    eventSource.close();
  };
}
```

## Running Only API Server (No UI)

If you only want the backend API without the web UI:

```bash
# Run API only
go run 4-structured-outputs/email_agent/main.go api

# Or with custom port
go run 4-structured-outputs/email_agent/main.go api --port 3000
```

This starts only the API server without the web UI, making it perfect for:
- Microservices architecture
- Backend integration
- Mobile app backends
- Third-party integrations

## CORS Configuration

For cross-origin requests, you may need to configure CORS. ADK handles this automatically, but you can customize it if needed.

## Authentication

The default ADK launcher doesn't include authentication. For production use, you should:
1. Add authentication middleware
2. Validate API keys/tokens
3. Implement rate limiting
4. Use HTTPS

## Session State Management

Sessions are stored in memory by default. For production:
- Implement persistent storage (Redis, PostgreSQL)
- Use the `OutputKey` mechanism to extract important data
- Consider session expiration policies

## Best Practices

1. **Always use unique session IDs** for different conversations
2. **Extract structured data from state** using the `OutputKey` you defined
3. **Handle SSE properly** for streaming responses
4. **Implement error handling** for API failures
5. **Use appropriate timeouts** for long-running agent operations
6. **Monitor session count** to manage memory usage

## Next Steps

- Integrate the API into your application
- Build a custom frontend
- Add authentication and authorization
- Implement session persistence
- Deploy to production with proper monitoring
