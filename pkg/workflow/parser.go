package workflow

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// State represents a workflow phase state.
type State string

const (
	StatePlan      State = "plan"
	StateImplement State = "implement"
	StateVerify    State = "verify"
	StateFinish    State = "finish"
	StateNoTask    State = "no_task"
)

// Phase represents a workflow phase.
type Phase struct {
	Number int
	Title  string
	Steps  []Step
}

// Step represents a workflow step.
type Step struct {
	Number string // "1.1"
	Title  string
	Body   string
}

// SkillRoute maps user intent to a skill.
type SkillRoute struct {
	Intent string
	Skill  string
}

// Workflow represents a parsed workflow.md definition.
type Workflow struct {
	Phases       []Phase
	Breadcrumbs  map[State]string
	SkillRouting []SkillRoute
	DoNotSkip    string
	TaskSystem   string
}

// StateMachine defines valid transitions between workflow states.
type StateMachine struct {
	States      []State
	Transitions map[State][]State
}

// NewStateMachine creates the default 4-phase state machine.
func NewStateMachine() *StateMachine {
	return &StateMachine{
		States: []State{StatePlan, StateImplement, StateVerify, StateFinish},
		Transitions: map[State][]State{
			StatePlan:      {StateImplement},
			StateImplement: {StateVerify},
			StateVerify:    {StateImplement, StateFinish},
			StateFinish:    {},
		},
	}
}

// CanTransition checks if a transition from one state to another is valid.
func (sm *StateMachine) CanTransition(from, to State) bool {
	valid, ok := sm.Transitions[from]
	if !ok {
		return false
	}
	for _, s := range valid {
		if s == to {
			return true
		}
	}
	return false
}

// NextStates returns all valid next states from the given state.
func (sm *StateMachine) NextStates(from State) []State {
	return sm.Transitions[from]
}

// Parser extracts workflow state information from workflow.md content.
type Parser struct{}

var (
	phaseRe      = regexp.MustCompile(`^#{2,3}\s+Phase\s+(\d+):?\s*(.*)$`)
	stepRe       = regexp.MustCompile(`^####\s+(\d+\.\d+)\s+(.*)$`)
	breadcrumbRe = regexp.MustCompile(`\[workflow-state:(\w+)\](.*?)\[/workflow-state:\w+\]`)
	stateRe      = regexp.MustCompile(`\[workflow-state:(\w+)\]`)
)

// Parse reads a workflow definition and extracts the full Workflow structure.
func (p *Parser) Parse(r io.Reader) (*Workflow, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read workflow: %w", err)
	}
	content := string(data)

	wf := &Workflow{
		Breadcrumbs: make(map[State]string),
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentPhase *Phase
	var currentStep *Step
	var inSkillTable bool
	var inDoNotSkip bool
	var inTaskSystem bool
	var doNotSkipLines []string
	var taskSystemLines []string
	var skillHeader bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Parse breadcrumbs
		if matches := breadcrumbRe.FindStringSubmatch(line); len(matches) > 0 {
			state := State(strings.ToLower(matches[1]))
			text := strings.TrimSpace(matches[2])
			wf.Breadcrumbs[state] = text
		}

		// Parse Phase Index breadcrumbs (single-line format)
		if strings.HasPrefix(trimmed, "[workflow-state:") {
			if endIdx := strings.Index(trimmed, "]"); endIdx > 0 {
				stateStr := trimmed[len("[workflow-state:"):endIdx]
				state := State(strings.ToLower(stateStr))
				// Extract text between opening and closing tags
				closeTag := fmt.Sprintf("[/workflow-state:%s]", stateStr)
				if closeIdx := strings.Index(trimmed, closeTag); closeIdx > endIdx {
					text := strings.TrimSpace(trimmed[endIdx+1 : closeIdx])
					wf.Breadcrumbs[state] = text
				}
			}
		}

		// Track sections
		if strings.HasPrefix(trimmed, "### Skill Routing") {
			inSkillTable = true
			skillHeader = true
			continue
		}
		if strings.HasPrefix(trimmed, "### DO NOT skip skills") {
			inDoNotSkip = true
			inSkillTable = false
			continue
		}
		if strings.HasPrefix(trimmed, "### Task System") {
			inTaskSystem = true
			inDoNotSkip = false
			continue
		}

		// Parse phases
		if matches := phaseRe.FindStringSubmatch(trimmed); len(matches) > 0 {
			if currentPhase != nil {
				if currentStep != nil {
					currentPhase.Steps = append(currentPhase.Steps, *currentStep)
					currentStep = nil
				}
				wf.Phases = append(wf.Phases, *currentPhase)
			}
			phaseNum := 0
			fmt.Sscanf(matches[1], "%d", &phaseNum)
			currentPhase = &Phase{
				Number: phaseNum,
				Title:  strings.TrimSpace(matches[2]),
			}
			inSkillTable = false
			inDoNotSkip = false
			inTaskSystem = false
			continue
		}

		// Parse steps
		if matches := stepRe.FindStringSubmatch(trimmed); len(matches) > 0 && currentPhase != nil {
			if currentStep != nil {
				currentPhase.Steps = append(currentPhase.Steps, *currentStep)
			}
			currentStep = &Step{
				Number: matches[1],
				Title:  strings.TrimSpace(matches[2]),
			}
			inSkillTable = false
			inDoNotSkip = false
			inTaskSystem = false
			continue
		}

		// Accumulate step body
		if currentStep != nil && trimmed != "" && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "[workflow-state:") {
			if currentStep.Body != "" {
				currentStep.Body += "\n"
			}
			currentStep.Body += trimmed
		}

		// Parse skill routing table
		if inSkillTable {
			if strings.HasPrefix(trimmed, "|") && !skillHeader {
				// Skip separator lines (e.g., |---|---|)
				if strings.Contains(trimmed, "---") {
					skillHeader = false
					continue
				}
				parts := strings.Split(trimmed, "|")
				if len(parts) >= 3 {
					intent := strings.TrimSpace(parts[1])
					skill := strings.TrimSpace(parts[2])
					if intent != "" && skill != "" && intent != "User intent" {
						wf.SkillRouting = append(wf.SkillRouting, SkillRoute{
							Intent: intent,
							Skill:  skill,
						})
					}
				}
			}
			skillHeader = false
		}

		// Accumulate Do Not Skip
		if inDoNotSkip {
			doNotSkipLines = append(doNotSkipLines, trimmed)
		}

		// Accumulate Task System
		if inTaskSystem {
			taskSystemLines = append(taskSystemLines, trimmed)
		}
	}

	// Finalize last phase/step
	if currentStep != nil && currentPhase != nil {
		currentPhase.Steps = append(currentPhase.Steps, *currentStep)
	}
	if currentPhase != nil {
		wf.Phases = append(wf.Phases, *currentPhase)
	}

	wf.DoNotSkip = strings.TrimSpace(strings.Join(doNotSkipLines, "\n"))
	wf.TaskSystem = strings.TrimSpace(strings.Join(taskSystemLines, "\n"))

	return wf, nil
}

// ExtractState finds the [workflow-state:STATUS] tag in content.
func (p *Parser) ExtractState(content string) (State, error) {
	matches := stateRe.FindStringSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("no workflow state found")
	}
	stateStr := strings.ToLower(matches[1])
	state := State(stateStr)

	// Accept any state (including custom ones)
	return state, nil
}

// InjectPrompt generates an injection prompt for the given state.
// Uses breadcrumbs from the parsed workflow when available.
func (p *Parser) InjectPrompt(state State) (string, error) {
	// Try to read workflow.md for custom breadcrumbs
	wf, err := p.loadWorkflow()
	if err == nil {
		if text, ok := wf.Breadcrumbs[state]; ok && text != "" {
			return fmt.Sprintf("<workflow-state>\n%s\n</workflow-state>", text), nil
		}
	}

	// Fallback to defaults
	switch state {
	case StatePlan:
		return "<workflow-state>\nYou are in the PLAN phase. Work with the user to clarify requirements and produce a PRD.\n</workflow-state>", nil
	case StateImplement:
		return "<workflow-state>\nYou are in the IMPLEMENT phase. Write code according to the PRD and spec.\n</workflow-state>", nil
	case StateVerify:
		return "<workflow-state>\nYou are in the VERIFY phase. Review the implementation against specs and run checks.\n</workflow-state>", nil
	case StateFinish:
		return "<workflow-state>\nYou are in the FINISH phase. Archive the task and update journals.\n</workflow-state>", nil
	case StateNoTask:
		return "<workflow-state>\nYou are in the NO TASK phase. Refer to workflow.md for current step.\n</workflow-state>", nil
	default:
		// Unknown state: generic fallback
		return "<workflow-state>\nRefer to workflow.md for current step.\n</workflow-state>", nil
	}
}

// ExtractStep extracts the body content of a specific workflow step (e.g., "2.1").
func (p *Parser) ExtractStep(content string, stepNum string) (string, error) {
	wf, err := p.parseString(content)
	if err != nil {
		return "", err
	}
	for _, phase := range wf.Phases {
		for _, step := range phase.Steps {
			if step.Number == stepNum {
				return step.Body, nil
			}
		}
	}
	return "", fmt.Errorf("step %s not found", stepNum)
}

// ParseSkillRouting extracts the skill routing table from workflow content.
func (p *Parser) ParseSkillRouting(content string) ([]SkillRoute, error) {
	wf, err := p.parseString(content)
	if err != nil {
		return nil, err
	}
	return wf.SkillRouting, nil
}

// parseString is a helper to parse workflow content from a string.
func (p *Parser) parseString(content string) (*Workflow, error) {
	return p.Parse(strings.NewReader(content))
}

// loadWorkflow attempts to read workflow.md from the default location.
func (p *Parser) loadWorkflow() (*Workflow, error) {
	// This is a best-effort read; the caller handles errors gracefully.
	// In practice, the workflow.md path is resolved at the CLI layer.
	return nil, fmt.Errorf("workflow.md not provided; use Parse() directly")
}
