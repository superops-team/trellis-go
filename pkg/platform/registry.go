package platform

import (
	"fmt"
	"sort"
	"sync"
)

// Registry manages the set of supported AI coding platforms.
type Registry struct {
	mu        sync.RWMutex
	platforms map[string]Platform
}

// NewRegistry creates a new registry pre-populated with all built-in platforms.
func NewRegistry() *Registry {
	r := &Registry{
		platforms: make(map[string]Platform),
	}
	for _, p := range builtins {
		_ = r.register(p) // builtins are guaranteed valid
	}
	return r
}

// Register adds a platform definition to the registry.
// Returns an error if the platform ID is already registered or invalid.
func (r *Registry) Register(p Platform) error {
	if err := p.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.register(p)
}

func (r *Registry) register(p Platform) error {
	if _, exists := r.platforms[p.ID]; exists {
		return fmt.Errorf("platform %q is already registered", p.ID)
	}
	r.platforms[p.ID] = p
	return nil
}

// Get retrieves a platform by its ID.
func (r *Registry) Get(id string) (Platform, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.platforms[id]
	return p, ok
}

// All returns all registered platforms, sorted by ID.
func (r *Registry) All() []Platform {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Platform, 0, len(r.platforms))
	for _, p := range r.platforms {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// ByClass returns platforms filtered by class, sorted by ID.
func (r *Registry) ByClass(c Class) []Platform {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Platform, 0)
	for _, p := range r.platforms {
		if p.Class == c {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// IDs returns all platform IDs, sorted alphabetically.
func (r *Registry) IDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.platforms))
	for id := range r.platforms {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
