package repl

// Scope represents the current level of TODO viewing.
type Scope int

const (
	ScopeRoot Scope = iota
	ScopeWorkspace
	ScopeProject
)
