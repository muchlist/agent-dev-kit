package agents

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

// ===== Agent Creation =====

// NewNewsAnalyst creates a specialized agent for news analysis
// Note: This agent uses GoogleSearch (a built-in tool), so it should be wrapped
// as an AgentTool when used by the manager due to Go ADK's limitation:
// a single agent can only use ONE built-in tool, and cannot mix built-in tools with custom tools
func NewNewsAnalyst(ctx context.Context, mdl model.LLM) (agent.Agent, error) {
	// Create news analyst agent with Google Search tool
	newsAnalyst, err := llmagent.New(llmagent.Config{
		Name:        "news_analyst",
		Model:       mdl,
		Description: "News analyst agent that searches and summarizes current news",
		Instruction: `You are a helpful assistant that can analyze news articles and provide a summary of the news.

When asked about news:
1. Use the google_search tool to search for relevant news articles
2. Summarize the key findings from the search results
3. Focus on recent and relevant information
4. If the user asks for news using a relative time (e.g., "today", "this week"),
   mention that you're searching for the most recent information

Tips for effective news search:
- Include relevant keywords and context in your search query
- For technology news, specify "tech" or "technology" in the query
- For recent news, include time-related terms like "latest", "recent", or "2024"

Example searches:
- "latest artificial intelligence news 2024"
- "recent Google product announcements"
- "technology industry trends this week"`,
		Tools: []tool.Tool{geminitool.GoogleSearch{}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create news analyst agent: %w", err)
	}

	return newsAnalyst, nil
}
