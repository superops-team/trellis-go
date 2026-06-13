package template

import (
	"text/template"
)

// RenderContext provides the data for template placeholder substitution.
type RenderContext struct {
	PlatformID    string
	PlatformName  string
	Developer     string
	ProjectName   string
	ExecutorAI    string
	AgentCapable  bool
	HasHooks      bool
	CLIFlag       string
	Extra         map[string]any
}

// FuncMap returns the template function map derived from the context.
func (ctx RenderContext) FuncMap() template.FuncMap {
	return template.FuncMap{
		"PlatformID":    func() string { return ctx.PlatformID },
		"PlatformName":  func() string { return ctx.PlatformName },
		"Developer":     func() string { return ctx.Developer },
		"ProjectName":   func() string { return ctx.ProjectName },
		"ExecutorAI":    func() string { return ctx.ExecutorAI },
		"AgentCapable":  func() bool { return ctx.AgentCapable },
		"HasHooks":      func() bool { return ctx.HasHooks },
		"CLIFlag":       func() string { return ctx.CLIFlag },
	}
}

// ToMap converts the context to a flat map for direct key access.
func (ctx RenderContext) ToMap() map[string]any {
	m := map[string]any{
		"PlatformID":    ctx.PlatformID,
		"PlatformName":  ctx.PlatformName,
		"Developer":     ctx.Developer,
		"ProjectName":   ctx.ProjectName,
		"ExecutorAI":    ctx.ExecutorAI,
		"AgentCapable":  ctx.AgentCapable,
		"HasHooks":      ctx.HasHooks,
		"CLIFlag":       ctx.CLIFlag,
	}
	for k, v := range ctx.Extra {
		m[k] = v
	}
	return m
}
