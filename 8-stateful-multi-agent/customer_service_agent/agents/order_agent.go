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

// ===== Order Agent Tool Structures =====

type getCurrentTimeArgs struct{}

type getCurrentTimeResults struct {
	CurrentTime string `json:"current_time"`
}

type refundCourseArgs struct{}

type refundCourseResults struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	CourseID  string `json:"course_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// ===== Tool Implementations =====

// getCurrentTime returns the current time in YYYY-MM-DD HH:MM:SS format
func getCurrentTime(ctx tool.Context, input getCurrentTimeArgs) (getCurrentTimeResults, error) {
	fmt.Println("--- Tool: get_current_time called ---")
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return getCurrentTimeResults{
		CurrentTime: currentTime,
	}, nil
}

// refundCourse simulates refunding the AI Marketing Platform course
// Updates state by removing the course from purchased_courses
func refundCourse(ctx tool.Context, input refundCourseArgs) (refundCourseResults, error) {
	fmt.Println("--- Tool: refund_course called ---")

	courseID := "ai_marketing_platform"
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	state := ctx.State()

	// Get current purchased courses
	var purchasedCourses []Course
	if val, err := state.Get("purchased_courses"); err == nil {
		if courses, ok := val.([]any); ok {
			for _, c := range courses {
				if courseMap, ok := c.(map[string]any); ok {
					course := Course{
						ID:           fmt.Sprintf("%v", courseMap["id"]),
						PurchaseDate: fmt.Sprintf("%v", courseMap["purchase_date"]),
					}
					purchasedCourses = append(purchasedCourses, course)
				}
			}
		}
	}

	// Check if user owns the course
	found := false
	for _, course := range purchasedCourses {
		if course.ID == courseID {
			found = true
			break
		}
	}

	if !found {
		return refundCourseResults{
			Status:  "error",
			Message: "You don't own this course, so it can't be refunded.",
		}, nil
	}

	// Create new list without the course to be refunded
	var newPurchasedCourses []map[string]any
	for _, course := range purchasedCourses {
		if course.ID != courseID {
			newPurchasedCourses = append(newPurchasedCourses, map[string]any{
				"id":            course.ID,
				"purchase_date": course.PurchaseDate,
			})
		}
	}

	// Update purchased courses in state
	state.Set("purchased_courses", newPurchasedCourses)

	// Get current interaction history
	var interactionHistory []map[string]any
	if val, err := state.Get("interaction_history"); err == nil {
		if history, ok := val.([]any); ok {
			for _, h := range history {
				if hMap, ok := h.(map[string]any); ok {
					interactionHistory = append(interactionHistory, hMap)
				}
			}
		}
	}

	// Add refund to interaction history
	interactionHistory = append(interactionHistory, map[string]any{
		"action":    "refund_course",
		"course_id": courseID,
		"timestamp": currentTime,
	})

	// Update interaction history in state
	state.Set("interaction_history", interactionHistory)

	return refundCourseResults{
		Status:    "success",
		Message:   "Successfully refunded the AI Marketing Platform course! Your $149 will be returned to your original payment method within 3-5 business days.",
		CourseID:  courseID,
		Timestamp: currentTime,
	}, nil
}

// ===== Agent Creation =====

// NewOrderAgent creates a specialized agent for order management and refunds
func NewOrderAgent(ctx context.Context, mdl model.LLM) (agent.Agent, error) {
	// Create get_current_time tool
	getCurrentTimeTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_current_time",
			Description: "Get the current time in the format YYYY-MM-DD HH:MM:SS",
		},
		getCurrentTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create get_current_time tool: %w", err)
	}

	// Create refund_course tool
	refundCourseTool, err := functiontool.New(
		functiontool.Config{
			Name:        "refund_course",
			Description: "Simulates refunding the AI Marketing Platform course and updates state",
		},
		refundCourse)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund_course tool: %w", err)
	}

	// Create order agent
	orderAgent, err := llmagent.New(llmagent.Config{
		Name:        "order_agent",
		Model:       mdl,
		Description: "Order agent for viewing purchase history and processing refunds",
		Instruction: `You are the order agent for the AI Developer Accelerator community.
Your role is to help users view their purchase history, course access, and process refunds.

<user_info>
Name: {user_name}
</user_info>

<purchase_info>
Purchased Courses: {purchased_courses}
</purchase_info>

<interaction_history>
{interaction_history}
</interaction_history>

When users ask about their purchases:
1. Check their course list from the purchase info above
   - Course information is stored as objects with "id" and "purchase_date" properties
2. Format the response clearly showing:
   - Which courses they own
   - When they were purchased (from the course.purchase_date property)

When users request a refund:
1. Verify they own the course they want to refund ("ai_marketing_platform")
2. If they own it:
   - **CRITICAL**: You MUST call the refund_course tool to actually process the refund
   - DO NOT just say the refund is processed - actually call the tool
   - After calling the tool, confirm the refund was successful
   - Remind them the money will be returned to their original payment method
   - If it's been more than 30 days, inform them that they are not eligible for a refund
3. If they don't own it:
   - Inform them they don't own the course, so no refund is needed

**IMPORTANT**: The refund_course tool is the ONLY way to remove courses from the user's account.
You must call it for every refund request, not just acknowledge the request.

Course Information:
- ai_marketing_platform: "Fullstack AI Marketing Platform" ($149)

Example Response for Purchase History:
"Here are your purchased courses:
1. Fullstack AI Marketing Platform
   - Purchased on: 2024-04-21 10:30:00
   - Full lifetime access"

Example Response for Refund:
"I've processed your refund for the Fullstack AI Marketing Platform course.
Your $149 will be returned to your original payment method within 3-5 business days.
The course has been removed from your account."

If they haven't purchased any courses:
- Let them know they don't have any courses yet
- Suggest talking to the sales agent about the AI Marketing Platform course

Remember:
- Be clear and professional
- Mention our 30-day money-back guarantee if relevant
- Direct course questions to course support
- Direct purchase inquiries to sales`,
		Tools: []tool.Tool{refundCourseTool, getCurrentTimeTool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order agent: %w", err)
	}

	return orderAgent, nil
}
