package registry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/secscan/secscan/backend/internal/scanner"
)

type Registry struct {
	modules map[string]scanner.Scanner
	order   []string
}

func New() *Registry {
	return &Registry{modules: make(map[string]scanner.Scanner)}
}

func (r *Registry) Register(module scanner.Scanner) {
	name := strings.ToLower(strings.TrimSpace(module.Name()))
	if _, exists := r.modules[name]; !exists {
		r.order = append(r.order, name)
	}
	r.modules[name] = module
}

func (r *Registry) Names() []string {
	names := append([]string(nil), r.order...)
	sort.Strings(names)
	return names
}

func (r *Registry) Select(requested []string) ([]scanner.Scanner, []string, error) {
	names := r.order
	if len(requested) > 0 {
		names = make([]string, 0, len(requested))
		seen := make(map[string]struct{}, len(requested))
		for _, item := range requested {
			name := strings.ToLower(strings.TrimSpace(item))
			if name == "" {
				continue
			}
			if _, exists := r.modules[name]; !exists {
				return nil, nil, fmt.Errorf("unknown scanner module: %s", item)
			}
			if _, exists := seen[name]; !exists {
				names = append(names, name)
				seen[name] = struct{}{}
			}
		}
	}

	if len(names) == 0 {
		return nil, nil, fmt.Errorf("at least one scanner module is required")
	}

	selected := make([]scanner.Scanner, 0, len(names))
	selectedNames := make([]string, 0, len(names))
	for _, name := range names {
		selected = append(selected, r.modules[name])
		selectedNames = append(selectedNames, name)
	}
	return selected, selectedNames, nil
}
