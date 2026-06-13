package workflow

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Parser extracts workflow state information from workflow.md content.
type Parser struct{}

// Parse reads a workflow definition and extracts the state machine.
func (p *Parser) Parse(r io.Reader) (*StateMachine, error) {
	// For now, return the default state machine.
	// In the future, this could parse custom workflow definitions.
	return NewStateMachine(), nil
}

// ExtractState finds the [workflow-state:STATUS] tag in content.
func (p *Parser) ExtractState(content string) (State, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "[workflow-state:"); idx != -1 {
			start := idx + len("[workflow-state:")
			end := strings.Index(line[start:], "]")
			if end == -1 {
				continue
			}
			stateStr := strings.TrimSpace(line[start : start+end])
			state := State(strings.ToLower(stateStr))
			// Validate against known states
			switch state {
			case StatePlan, StateImplement, StateVerify, StateFinish:
				return state, nil
			default:
				return "", fmt.Errorf("unknown workflow state: %s", stateStr)
			}
		}
	}
	return "", fmt.Errorf("no workflow state found")
}

// InjectPrompt generates an injection prompt for the given state.
func (p *Parser) InjectPrompt(state State) (string, error) {
	switch state {
	case StatePlan:
		return "You are in the PLAN phase. Work with the user to clarify requirements and produce a PRD.", nil
	case StateImplement:
		return "You are in the IMPLEMENT phase. Write code according to the PRD and spec.", nil
	case StateVerify:
		return "You are in the VERIFY phase. Review the implementation against specs and run checks.", nil
	case StateFinish:
		return "You are in the FINISH phase. Archive the task and update journals.", nil
	default:
		return "", fmt.Errorf("unknown state: %s", state)
	}
}
