package ptrstruct

import (
	"go/types"
	"slices"
	"strings"
)

// Violation describes a struct-by-value occurrence found by the type walker.
type Violation struct {
	Path     string       // Human-readable path, e.g. "slice element", "map value".
	TypeName string       // Short name of the struct type found, e.g. "User".
	Named    *types.Named // The named type if the violation is on a named struct; nil for anonymous structs.
	steps    []pathStep
}

// walker holds traversal state for a single FindViolation call.
// Cycle detection uses seen keyed on *types.Named because Go's type system
// only permits recursion through named types.
type walker struct {
	cfg  *Config
	cls  *Classifier
	seen map[*types.Named]bool
}

// FindViolation walks t recursively and returns the first value-struct
// violation, or nil if the type is clean.
// cls may be nil if no allowlist is configured.
func FindViolation(t types.Type, cfg *Config, cls *Classifier) *Violation {
	w := &walker{cfg: cfg, cls: cls, seen: make(map[*types.Named]bool)}
	return w.walk(t, nil)
}

func (w *walker) walk(t types.Type, path []pathStep) *Violation {
	t = types.Unalias(t)

	switch tt := t.(type) {
	case *types.Pointer:
		return w.walkPointer(tt, path)
	case *types.Named:
		return w.walkNamed(tt, path)
	case *types.Struct:
		return walkStruct(tt, path)
	case *types.Slice:
		return w.walkSlice(tt, path)
	case *types.Map:
		return w.walkMap(tt, path)
	case *types.Array:
		return w.walkArray(tt, path)
	case *types.Chan:
		return w.walkChan(tt, path)
	default:
		return nil
	}
}

func (w *walker) walkPointer(tt *types.Pointer, path []pathStep) *Violation {
	elem := types.Unalias(tt.Elem())

	if isStructType(elem) {
		return nil
	}

	return w.walk(elem, appendPath(path, pathPointer))
}

func (w *walker) walkNamed(tt *types.Named, path []pathStep) *Violation {
	if w.seen[tt] {
		return nil
	}
	w.seen[tt] = true

	if w.cls != nil && w.cls.IsAllowed(tt) {
		return nil
	}

	under := tt.Underlying()
	if st, ok := under.(*types.Struct); ok {
		if st.NumFields() == 0 {
			return nil
		}
		return newViolation(path, tt.Obj().Name(), tt)
	}

	return w.walk(under, path)
}

func walkStruct(tt *types.Struct, path []pathStep) *Violation {
	if tt.NumFields() == 0 {
		return nil
	}
	return newViolation(path, "struct{...}", nil)
}

func (w *walker) walkSlice(tt *types.Slice, path []pathStep) *Violation {
	if !w.cfg.SliceElem {
		return nil
	}
	return w.walk(tt.Elem(), appendPath(path, pathSliceElement))
}

func (w *walker) walkMap(tt *types.Map, path []pathStep) *Violation {
	if w.cfg.MapKey {
		if v := w.walk(tt.Key(), appendPath(path, pathMapKey)); v != nil {
			return v
		}
	}
	if !w.cfg.MapValue {
		return nil
	}
	return w.walk(tt.Elem(), appendPath(path, pathMapValue))
}

func (w *walker) walkArray(tt *types.Array, path []pathStep) *Violation {
	if !w.cfg.ArrayElem {
		return nil
	}
	return w.walk(tt.Elem(), appendPath(path, pathArrayElement))
}

func (w *walker) walkChan(tt *types.Chan, path []pathStep) *Violation {
	if !w.cfg.ChanElem {
		return nil
	}
	return w.walk(tt.Elem(), appendPath(path, pathChanElement))
}

type pathStep string

// Path segment constants for violation paths.
const (
	pathPointer      pathStep = "pointer"
	pathSliceElement pathStep = "slice element"
	pathMapKey       pathStep = "map key"
	pathMapValue     pathStep = "map value"
	pathArrayElement pathStep = "array element"
	pathChanElement  pathStep = "chan element"
)

// isStructType reports whether t is a struct (named or anonymous).
func isStructType(t types.Type) bool {
	switch tt := t.(type) {
	case *types.Named:
		_, ok := tt.Underlying().(*types.Struct)
		return ok
	case *types.Struct:
		return true
	default:
		return false
	}
}

func newViolation(path []pathStep, typeName string, named *types.Named) *Violation {
	steps := slices.Clone(path)
	return &Violation{
		Path:     formatPath(steps),
		TypeName: typeName,
		Named:    named,
		steps:    steps,
	}
}

func appendPath(base []pathStep, segment pathStep) []pathStep {
	return append(slices.Clone(base), segment)
}

func formatPath(path []pathStep) string {
	if len(path) == 0 {
		return ""
	}

	parts := make([]string, len(path))
	for i, step := range path {
		parts[i] = string(step)
	}
	return strings.Join(parts, " -> ")
}
