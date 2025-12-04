# Go ADK Execution Modes: Launcher vs Runner

This guide explains the different ways to run agents in Google's Agent Development Kit (ADK) for Go, helping you choose the right approach for your needs.

## Quick Overview

| Method | When to Use | Example Command |
|--------|------------|-----------------|
| **Launcher Pattern** | Development, Testing, Production | `go run main.go web api webui` |
| **Runner Pattern** | Custom Applications, Programmatic Control | `runner.New().Run()` |

---

## 1. Launcher Pattern (Most Common)

The launcher pattern is used in all examples in this repository. It provides a comprehensive framework for running agents with multiple execution modes.

### Basic Structure

```go
import (
    "google.golang.org/adk/cmd/launcher"
    "google.golang.org/adk/cmd/launcher/full"
)

func main() {
    // Create your agent
    a, err := llmagent.New(llmagent.Config{
        Name:        "my_agent",
        Model:       model,
        Instruction: "Your agent instructions...",
    })

    // Configure launcher with your agent
    config := &launcher.Config{
        AgentLoader: agent.NewSingleLoader(a),
    }

    // Execute with command-line arguments
    l := full.NewLauncher()
    if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
        log.Fatalf("Run failed: %v", err)
    }
}
```

### Available Commands

All Go examples in this repository support these commands:

```bash
# Web UI (Recommended for development)
go run main.go web api webui
# Opens at: http://localhost:8080

# Terminal interaction
go run main.go run

# API server only
go run main.go api

# See all available options
go run main.go help
```

### Makefile Commands

This repository provides convenient shortcuts:

```bash
make run/1   # Basic agent with web UI
make run/2   # Tool agent with web UI
make run/7   # Multi-agent system with web UI
make run/10  # Sequential agent with web UI
```

**Key Features:**
- ✅ **Web UI**: Interactive chat interface at http://localhost:8080
- ✅ **REST API**: HTTP endpoints for programmatic access
- ✅ **CLI Mode**: Terminal-based interaction
- ✅ **Production Ready**: Built for deployment
- ✅ **Easy Testing**: Quick development cycle

---

## 2. Runner Pattern (Programmatic Control)

The runner pattern gives you direct programmatic control over agent execution, perfect for custom applications.

### Basic Structure

```go
import "google.golang.org/adk/runner"

func main() {
    // Create session service for state management
    sessionService, _ := session.NewInMemoryService()

    // Create runner with direct control
    r, err := runner.New(runner.Config{
        AppName:        "MyCustomApp",
        Agent:          myAgent,
        SessionService: sessionService,
    })

    // Run with specific context
    for event, err := range r.Run(ctx, userID, sessionID, message, runConfig) {
        if event.Content != nil {
            fmt.Printf("Response: %s\n", event.Content.Parts[0].Text)
        }
    }
}
```

### Custom Session Management

```go
// Create session with initial state
newSession := sessionService.CreateSession(
    app_name: "Customer Support",
    user_id: "user123",
    state: map[string]any{
        "user_name": "John Doe",
        "preferences": map[string]string{
            "language": "en",
            "timezone": "UTC",
        },
    },
)

// Run with that session
for event := range r.Run(ctx, "user123", newSession.ID, "Hello!") {
    // Process events with full control
}
```

**Key Features:**
- ✅ **Full Control**: Direct access to execution flow
- ✅ **Custom Sessions**: Programmatic session management
- ✅ **Event Processing**: Handle each event individually
- ✅ **State Access**: Direct state manipulation
- ✅ **Integration**: Perfect for custom applications

---

## 3. When to Use Each

### 🚀 **Use Launcher Pattern When:**

#### **Development & Testing**
```bash
# Quick interactive testing
make run/10  # Launch sequential agent with web UI
```

#### **Simple Deployments**
- One-command deployment
- Built-in web interface
- Standardized configuration

#### **When You Need:**
- Interactive testing and debugging
- Visual agent interaction
- Quick prototyping
- Built-in monitoring tools

**Perfect for:**
- All examples in this repository
- Development environments
- Simple customer service bots
- Quick proof of concepts

### 💻 **Use Runner Pattern When:**

#### **Custom Applications**
```go
// Custom customer service integration
func handleCustomerRequest(customerID, message string) {
    for event := range r.Run(ctx, customerID, sessionID, message, config) {
        // Add custom business logic
        if event.Content != nil {
            response := event.Content.Parts[0].Text
            saveToDatabase(customerID, response)
            triggerWebhook(customerID, response)
        }
    }
}
```

#### **When You Need:**
- Full control over execution flow
- Custom session management
- Integration with existing systems
- Automated testing
- Batch processing

**Perfect for:**
- Custom web applications with agents
- Automated report generation
- Integration with databases and APIs
- Data processing pipelines
- Enterprise applications

---

## 4. Code Examples

### Launcher Pattern (Standard Approach)

```go
// Used in: main.go files of all examples
package main

import (
    "google.golang.org/adk/cmd/launcher"
    "google.golang.org/adk/cmd/launcher/full"
)

func main() {
    // Create agent (example from sequential agent)
    sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
        AgentConfig: agent.Config{
            Name:        "LeadQualificationPipeline",
            SubAgents:   []agent.Agent{validator, scorer, recommender},
        },
    })

    // Standard launcher setup
    config := &launcher.Config{
        AgentLoader: agent.NewSingleLoader(sequentialAgent),
    }

    l := full.NewLauncher()
    l.Execute(ctx, config, os.Args[1:])
}
```

### Runner Pattern (Custom Application)

```go
// Used for: Custom integrations
package main

import "google.golang.org/adk/runner"

func main() {
    // Create runner with custom session service
    sessionService, _ := session.NewInMemoryService()

    r, _ := runner.New(runner.Config{
        AppName:        "LeadProcessor",
        Agent:          sequentialAgent,
        SessionService: sessionService,
    })

    // Process leads programmatically
    leads := getLeadsFromDatabase()
    for _, lead := range leads {
        processLead(r, lead)
    }
}

func processLead(r *runner.Runner, Lead) {
    // Create session for this lead
    session := r.SessionService.CreateSession("LeadProcessor", lead.ID, lead.InitialState)

    // Run qualification pipeline
    events := r.Run(ctx, lead.ID, session.ID, lead.Info, runner.RunConfig{})

    // Process results
    qualification := extractQualificationResult(events)
    updateLeadInDatabase(lead.ID, qualification)
    notifySalesTeam(lead.ID, qualification)
}
```

---

## 5. Migration from Launcher to Runner

### Step 1: Replace Launcher Setup

**Launcher (Current):**
```go
config := &launcher.Config{
    AgentLoader: agent.NewSingleLoader(a),
}
l := full.NewLauncher()
l.Execute(ctx, config, os.Args[1:])
```

**Runner (Target):**
```go
sessionService, _ := session.NewInMemoryService()
r, _ := runner.New(runner.Config{
    AppName:        "MyApp",
    Agent:          a,
    SessionService: sessionService,
})
```

### Step 2: Replace Command Execution

**Launcher (CLI-driven):**
```bash
go run main.go web api webui
```

**Runner (Programmatic):**
```go
for event := range r.Run(ctx, userID, sessionID, message, config) {
    // Handle events
}
```

### Step 3: Add Session Management

```go
// Create session with initial state
session := sessionService.CreateSession(
    app_name: "MyApp",
    user_id: userID,
    state: map[string]any{
        "user_preferences": preferences,
        "conversation_history": []string{},
    },
)

// Use session in runner
r.Run(ctx, userID, session.ID, message, config)
```

---

## 6. Best Practices

### For Launcher Pattern
- **Use for development**: Quick testing with web UI
- **Standard deployment**: Build binary with launcher for production
- **Monitoring**: Built-in logging and metrics

### For Runner Pattern
- **Error handling**: Always check error returns
- **Session cleanup**: Proper session lifecycle management
- **Resource management**: Handle concurrent requests properly

### Performance Considerations
- **Launcher**: Optimized for single-user development
- **Runner**: Better for high-volume, multi-user applications

---

## 📚 **Quick Decision Guide**

| Scenario | Recommended Approach | Why |
|----------|---------------------|-----|
| **Learning ADK** | Launcher Pattern | Web UI for interactive testing |
| **Simple Chatbot** | Launcher Pattern | Quick deployment with built-in UI |
| **Customer Service System** | Runner Pattern | Custom integration with ticketing |
| **Data Processing Pipeline** | Runner Pattern | Batch processing and automation |
| **Internal Tool** | Launcher Pattern | Fast development and deployment |
| **Enterprise Application** | Runner Pattern | Full control and integration |

Choose the approach that fits your development and deployment needs! 🚀