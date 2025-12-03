package agents

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// ===== Course Structure =====

// Course represents a purchased course
type Course struct {
	ID           string `json:"id"`
	PurchaseDate string `json:"purchase_date"`
}

// ===== Sales Agent Tool Structures =====

type purchaseCourseArgs struct{}

type purchaseCourseResults struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	CourseID  string `json:"course_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// ===== Tool Implementation =====

// purchaseCourse simulates purchasing the AI Marketing Platform course
// Updates state with purchase information
func purchaseCourse(ctx tool.Context, input purchaseCourseArgs) (purchaseCourseResults, error) {
	fmt.Println("--- Tool: purchase_course called ---")

	courseID := "ai_marketing_platform"
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	state := ctx.State()

	// Get current purchased courses
	var purchasedCourses []Course
	if val, err := state.Get("purchased_courses"); err == nil {
		if courses, ok := val.([]interface{}); ok {
			for _, c := range courses {
				if courseMap, ok := c.(map[string]interface{}); ok {
					course := Course{
						ID:           fmt.Sprintf("%v", courseMap["id"]),
						PurchaseDate: fmt.Sprintf("%v", courseMap["purchase_date"]),
					}
					purchasedCourses = append(purchasedCourses, course)
				}
			}
		}
	}

	// Check if user already owns the course
	for _, course := range purchasedCourses {
		if course.ID == courseID {
			return purchaseCourseResults{
				Status:  "error",
				Message: "You already own this course!",
			}, nil
		}
	}

	// Add the new course
	purchasedCourses = append(purchasedCourses, Course{
		ID:           courseID,
		PurchaseDate: currentTime,
	})

	// Convert to []map[string]any for state storage
	var coursesForState []map[string]any
	for _, course := range purchasedCourses {
		coursesForState = append(coursesForState, map[string]any{
			"id":            course.ID,
			"purchase_date": course.PurchaseDate,
		})
	}

	// Update purchased courses in state
	state.Set("purchased_courses", coursesForState)

	// Get current interaction history
	var interactionHistory []map[string]interface{}
	if val, err := state.Get("interaction_history"); err == nil {
		if history, ok := val.([]interface{}); ok {
			for _, h := range history {
				if hMap, ok := h.(map[string]interface{}); ok {
					interactionHistory = append(interactionHistory, hMap)
				}
			}
		}
	}

	// Add purchase to interaction history
	interactionHistory = append(interactionHistory, map[string]interface{}{
		"action":    "purchase_course",
		"course_id": courseID,
		"timestamp": currentTime,
	})

	// Update interaction history in state
	state.Set("interaction_history", interactionHistory)

	return purchaseCourseResults{
		Status:    "success",
		Message:   "Successfully purchased the AI Marketing Platform course!",
		CourseID:  courseID,
		Timestamp: currentTime,
	}, nil
}

// ===== Agent Creation =====

// NewSalesAgent creates a specialized agent for course sales
func NewSalesAgent(ctx context.Context, mdl model.LLM) (agent.Agent, error) {
	// Create purchase_course tool
	purchaseCourseTool, err := functiontool.New(
		functiontool.Config{
			Name:        "purchase_course",
			Description: "Simulates purchasing the AI Marketing Platform course and updates state",
		},
		purchaseCourse)
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase_course tool: %w", err)
	}

	// Create sales agent
	salesAgent, err := llmagent.New(llmagent.Config{
		Name:        "sales_agent",
		Model:       mdl,
		Description: "Sales agent for the AI Marketing Platform course",
		Instruction: `You are a sales agent for the AI Developer Accelerator community, specifically handling sales
for the Fullstack AI Marketing Platform course.

<user_info>
Name: {user_name}
</user_info>

<purchase_info>
Purchased Courses: {purchased_courses}
</purchase_info>

<interaction_history>
{interaction_history}
</interaction_history>

Course Details:
- Name: Fullstack AI Marketing Platform
- Price: $149
- Value Proposition: Learn to build AI-powered marketing automation apps
- Includes: 6 weeks of group support with weekly coaching calls

When interacting with users:
1. Check if they already own the course (check purchased_courses above)
   - Course information is stored as objects with "id" and "purchase_date" properties
   - The course id is "ai_marketing_platform"
2. If they own it:
   - Remind them they have access
   - Ask if they need help with any specific part
   - Direct them to course support for content questions

3. If they don't own it:
   - Explain the course value proposition
   - Mention the price ($149)
   - If they want to purchase:
       - Use the purchase_course tool
       - Confirm the purchase
       - Ask if they'd like to start learning right away

4. After any interaction:
   - The state will automatically track the interaction
   - Be ready to hand off to course support after purchase

Remember:
- Be helpful but not pushy
- Focus on the value and practical skills they'll gain
- Emphasize the hands-on nature of building a real AI application`,
		Tools: []tool.Tool{purchaseCourseTool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sales agent: %w", err)
	}

	return salesAgent, nil
}
