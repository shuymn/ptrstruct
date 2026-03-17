package ptrstruct

import "go/types"

// Violation describes a struct-by-value occurrence found by the type walker.
type Violation struct {
	Path     string       // Human-readable path, e.g. "slice element", "map value".
	TypeName string       // Short name of the struct type found, e.g. "User".
	Named    *types.Named // The named type if the violation is on a named struct; nil for anonymous structs.
}

// FindViolation walks t recursively and returns the first value-struct
// violation, or nil if the type is clean.
// cls may be nil if no allowlist is configured.
func FindViolation(t types.Type, cfg *Config, cls *Classifier) *Violation {
	return findViolation(t, cfg, cls, "")
}

func findViolation(t types.Type, cfg *Config, cls *Classifier, path string) *Violation {
	t = types.Unalias(t)

	switch tt := t.(type) {
	case *types.Pointer:
		return findViolationPointer(tt, cfg, cls, path)
	case *types.Named:
		return findViolationNamed(tt, cfg, cls, path)
	case *types.Struct:
		return findViolationStruct(tt, path)
	case *types.Slice:
		return findViolationSlice(tt, cfg, cls, path)
	case *types.Map:
		return findViolationMap(tt, cfg, cls, path)
	case *types.Array:
		return findViolationArray(tt, cfg, cls, path)
	case *types.Chan:
		return findViolationChan(tt, cfg, cls, path)
	default:
		return nil
	}
}

func findViolationPointer(tt *types.Pointer, cfg *Config, cls *Classifier, path string) *Violation {
	elem := types.Unalias(tt.Elem())

	if isStructType(elem) {
		return nil
	}

	return findViolation(elem, cfg, cls, appendPath(path, pathPointer))
}

func findViolationNamed(tt *types.Named, cfg *Config, cls *Classifier, path string) *Violation {
	if cls != nil && cls.IsAllowed(tt) {
		return nil
	}

	under := tt.Underlying()
	if st, ok := under.(*types.Struct); ok {
		if st.NumFields() == 0 {
			return nil
		}
		return &Violation{
			Path:     path,
			TypeName: tt.Obj().Name(),
			Named:    tt,
		}
	}

	return findViolation(under, cfg, cls, path)
}

func findViolationStruct(tt *types.Struct, path string) *Violation {
	if tt.NumFields() == 0 {
		return nil
	}
	return &Violation{
		Path:     path,
		TypeName: "struct{...}",
		Named:    nil,
	}
}

func findViolationSlice(tt *types.Slice, cfg *Config, cls *Classifier, path string) *Violation {
	if !cfg.SliceElem {
		return nil
	}
	return findViolation(tt.Elem(), cfg, cls, appendPath(path, pathSliceElement))
}

func findViolationMap(tt *types.Map, cfg *Config, cls *Classifier, path string) *Violation {
	if cfg.MapKey {
		if v := findViolation(tt.Key(), cfg, cls, appendPath(path, pathMapKey)); v != nil {
			return v
		}
	}
	if !cfg.MapValue {
		return nil
	}
	return findViolation(tt.Elem(), cfg, cls, appendPath(path, pathMapValue))
}

func findViolationArray(tt *types.Array, cfg *Config, cls *Classifier, path string) *Violation {
	if !cfg.ArrayElem {
		return nil
	}
	return findViolation(tt.Elem(), cfg, cls, appendPath(path, pathArrayElement))
}

func findViolationChan(tt *types.Chan, cfg *Config, cls *Classifier, path string) *Violation {
	if !cfg.ChanElem {
		return nil
	}
	return findViolation(tt.Elem(), cfg, cls, appendPath(path, pathChanElement))
}

// Path segment constants for violation paths.
const (
	pathPointer      = "pointer"
	pathSliceElement = "slice element"
	pathMapKey       = "map key"
	pathMapValue     = "map value"
	pathArrayElement = "array element"
	pathChanElement  = "chan element"
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

func appendPath(base, segment string) string {
	if base == "" {
		return segment
	}
	return base + " -> " + segment
}
