package ptrstruct

import (
	"fmt"
	"go/types"
	"regexp"
)

// Classifier checks whether a named type is exempted by the allowlist.
type Classifier struct {
	types    map[string]bool
	patterns []*regexp.Regexp
	packages map[string]bool
}

// NewClassifier creates a Classifier from the allowlist fields of cfg.
// It returns an error if any AllowPatterns entry is not a valid regexp.
func NewClassifier(cfg *Config) (*Classifier, error) {
	c := &Classifier{
		types:    make(map[string]bool, len(cfg.AllowTypes)),
		packages: make(map[string]bool, len(cfg.AllowPackages)),
	}
	for _, t := range cfg.AllowTypes {
		c.types[t] = true
	}
	for _, p := range cfg.AllowPackages {
		c.packages[p] = true
	}
	for _, pat := range cfg.AllowPatterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("ptrstruct: invalid allow-pattern %q: %w", pat, err)
		}
		c.patterns = append(c.patterns, re)
	}
	return c, nil
}

// IsAllowed reports whether the named type is exempted by any allowlist entry.
func (c *Classifier) IsAllowed(named *types.Named) bool {
	fqn := qualifiedName(named)

	if c.types[fqn] {
		return true
	}

	if named.Obj().Pkg() != nil && c.packages[named.Obj().Pkg().Path()] {
		return true
	}

	for _, re := range c.patterns {
		if re.MatchString(fqn) {
			return true
		}
	}

	return false
}

func qualifiedName(named *types.Named) string {
	obj := named.Obj()
	if obj.Pkg() == nil {
		return obj.Name()
	}
	return obj.Pkg().Path() + "." + obj.Name()
}
