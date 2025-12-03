// Package main demonstrates a stateful multi-agent system in ADK.
// This example combines persistent state management with multi-agent delegation
// for a customer service system that remembers user information and interactions.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/session/database"

	"github.com/muchlist/agent-dev-kit/8-stateful-multi-agent/customer_service_agent/agents"
)

const (
	APP_NAME   = "customer_service"
	MODEL_NAME = "gemini-2.0-flash"
	DB_FILE    = "./customer_service_data.db"
)

// ===== Customer Service Agent Creation =====

// createCustomerServiceAgent creates the root customer service agent that coordinates specialized agents
func createCustomerServiceAgent(_ context.Context, mdl model.LLM, policyAgent, salesAgent, courseSupportAgent, orderAgent agent.Agent) (agent.Agent, error) {
	// Create customer service agent with all sub-agents
	customerServiceAgent, err := llmagent.New(llmagent.Config{
		Name:        "customer_service",
		Model:       mdl,
		Description: "Customer service agent for AI Developer Accelerator community",
		Instruction: `You are the primary customer service agent for the AI Developer Accelerator community.
Your role is to help users with their questions and direct them to the appropriate specialized agent.

**Core Capabilities:**

1. Query Understanding & Routing
   - Understand user queries about policies, course purchases, course support, and orders
   - Direct users to the appropriate specialized agent
   - Maintain conversation context using state

2. State Management
   - Track user interactions in state['interaction_history']
   - Monitor user's purchased courses in state['purchased_courses']
     - Course information is stored as objects with "id" and "purchase_date" properties
   - Use state to provide personalized responses

**User Information:**
<user_info>
Name: {user_name}
</user_info>

**Purchase Information:**
<purchase_info>
Purchased Courses: {purchased_courses}
</purchase_info>

**Interaction History:**
<interaction_history>
{interaction_history}
</interaction_history>

You have access to the following specialized agents:

1. Policy Agent
   - For questions about community guidelines, course policies, refunds
   - Direct policy-related queries here

2. Sales Agent
   - For questions about purchasing the AI Marketing Platform course
   - Handles course purchases and updates state
   - Course price: $149

3. Course Support Agent
   - For questions about course content
   - Only available for courses the user has purchased
   - Check if a course with id "ai_marketing_platform" exists in the purchased courses before directing here

4. Order Agent
   - For checking purchase history and processing refunds
   - Shows courses user has bought
   - Can process course refunds (30-day money-back guarantee)
   - References the purchased courses information

Tailor your responses based on the user's purchase history and previous interactions.
When the user hasn't purchased any courses yet, encourage them to explore the AI Marketing Platform.
When the user has purchased courses, offer support for those specific courses.

When users express dissatisfaction or ask for a refund:
- IMMEDIATELY DELEGATE to the Order Agent - DO NOT process refunds yourself
- The Order Agent has the refund_course tool to actually process the refund
- Mention our 30-day money-back guarantee policy

**IMPORTANT ROUTING RULES:**
- For purchases: DELEGATE to Sales Agent
- For refunds or order history: DELEGATE to Order Agent
- For course content help: DELEGATE to Course Support Agent
- For policy questions: DELEGATE to Policy Agent
- You are a COORDINATOR - always delegate to the appropriate specialist, never handle their tasks directly

Always maintain a helpful and professional tone. If you're unsure which agent to delegate to,
ask clarifying questions to better understand the user's needs.`,
		SubAgents: []agent.Agent{policyAgent, salesAgent, courseSupportAgent, orderAgent},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create customer service agent: %w", err)
	}

	return customerServiceAgent, nil
}

// ===== Main Function =====

func main() {
	godotenv.Load()
	ctx := context.Background()

	// Create the Gemini model
	model, err := gemini.NewModel(ctx, MODEL_NAME, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	fmt.Println("🤖 Creating specialized agents...")

	// Create all specialized agents
	policyAgent, err := agents.NewPolicyAgent(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create policy agent: %v", err)
	}
	fmt.Println("  ✓ Policy Agent created")

	salesAgent, err := agents.NewSalesAgent(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create sales agent: %v", err)
	}
	fmt.Println("  ✓ Sales Agent created")

	courseSupportAgent, err := agents.NewCourseSupportAgent(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create course support agent: %v", err)
	}
	fmt.Println("  ✓ Course Support Agent created")

	orderAgent, err := agents.NewOrderAgent(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create order agent: %v", err)
	}
	fmt.Println("  ✓ Order Agent created")

	// Create customer service manager agent
	fmt.Println("🎯 Creating customer service manager agent...")
	customerServiceAgent, err := createCustomerServiceAgent(ctx, model, policyAgent, salesAgent, courseSupportAgent, orderAgent)
	if err != nil {
		log.Fatalf("Failed to create customer service agent: %v", err)
	}
	fmt.Println("  ✓ Customer Service Agent created")

	// ===== Session Management Setup =====

	fmt.Println("\n📦 Setting up session management...")

	// Create database session service with SQLite
	// This properly persists state changes made by tools
	sessionService, err := database.NewSessionService(
		sqlite.Open(DB_FILE),
		&gorm.Config{
			PrepareStmt: true,
			Logger:      logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		log.Fatalf("Failed to create database session service: %v", err)
	}

	// Initialize database schema
	if err := database.AutoMigrate(sessionService); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	fmt.Println("✅ Connected to database:", DB_FILE)

	// Wrap session service to provide default initial state for new sessions
	initialState := map[string]any{
		"user_name":           "Muchlis",
		"purchased_courses":   []any{},
		"interaction_history": []any{},
	}
	wrappedSessionService := &sessionServiceWithDefaults{
		Service:      sessionService,
		initialState: initialState,
	}

	// ===== Launch with Web/API/WebUI =====

	fmt.Println("\n🚀 Launching Stateful Multi-Agent System...")
	fmt.Println("========================================")

	// Configure and launch the agent with session service
	config := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(customerServiceAgent),
		SessionService: wrappedSessionService,
	}

	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

// sessionServiceWithDefaults wraps a session service to provide default initial state
type sessionServiceWithDefaults struct {
	session.Service
	initialState map[string]any
}

// Create wraps the Create method to ensure initial state is set
func (s *sessionServiceWithDefaults) Create(ctx context.Context, req *session.CreateRequest) (*session.CreateResponse, error) {
	// If no state provided, use initial state
	if req.State == nil || len(req.State) == 0 {
		req.State = s.initialState
	}
	return s.Service.Create(ctx, req)
}
