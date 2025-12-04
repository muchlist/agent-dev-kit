## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run/1: run the basic-agent
run/1:
	go run 1-basic-agent/greeting_agent/main.go web api webui

## run/2: run the tool-agent
run/2:
	go run 2-tool-agent/tool_agent/main.go web api webui

## run/3: run the dad-joke-agent (uses Gemini, not LiteLLM)
run/3:
	go run 3-litellm-agent/dad_joke_agent/main.go web api webui

## run/4: run the email-agent with structured outputs
run/4:
	go run 4-structured-outputs/email_agent/main.go web api webui

## run/5: run the question-answering-agent with sessions and state
run/5:
	go run 5-sessions-and-state/question_answering_agent/main.go

## run/6: run the memory-agent with persistent database storage
run/6:
	go run 6-persistent-storage/memory_agent/main.go

## run/7: run the multi-agent manager system with specialized agents
run/7:
	go run 7-multi-agent/manager_agent/main.go web api webui

## run/8: run the stateful multi-agent customer service system
run/8:
	go run 8-stateful-multi-agent/customer_service_agent/main.go web api webui

## run/9a: run the before/after agent callbacks example
run/9a:
	go run 9-callbacks/before_after_agent/main.go web api webui

## run/9b: run the before/after model callbacks example
run/9b:
	go run 9-callbacks/before_after_model/main.go web api webui

## run/9c: run the before/after tool callbacks example
run/9c:
	go run 9-callbacks/before_after_tool/main.go web api webui