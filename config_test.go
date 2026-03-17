package ptrstruct

import "testing"

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	// Rule toggles (Phase 1: on)
	if !cfg.Receiver {
		t.Error("Receiver should default to true")
	}
	if !cfg.Param {
		t.Error("Param should default to true")
	}
	if !cfg.Result {
		t.Error("Result should default to true")
	}
	if !cfg.Field {
		t.Error("Field should default to true")
	}

	// Phase 2 toggles (off)
	if cfg.InterfaceMethod {
		t.Error("InterfaceMethod should default to false")
	}
	if cfg.FuncType {
		t.Error("FuncType should default to false")
	}
	if cfg.NamedType {
		t.Error("NamedType should default to false")
	}

	// Container checks
	if !cfg.SliceElem {
		t.Error("SliceElem should default to true")
	}
	if !cfg.MapValue {
		t.Error("MapValue should default to true")
	}
	if cfg.MapKey {
		t.Error("MapKey should default to false")
	}
	if cfg.ArrayElem {
		t.Error("ArrayElem should default to false")
	}
	if cfg.ChanElem {
		t.Error("ChanElem should default to false")
	}

	// File filtering
	if !cfg.IgnoreGenerated {
		t.Error("IgnoreGenerated should default to true")
	}
	if cfg.IgnoreTests {
		t.Error("IgnoreTests should default to false")
	}

	// Suppression toggles
	if !cfg.HonorNolint {
		t.Error("HonorNolint should default to true")
	}
	if !cfg.HonorNolintAll {
		t.Error("HonorNolintAll should default to true")
	}

	// Allowlists
	if cfg.AllowTypes != nil {
		t.Error("AllowTypes should default to nil")
	}
	if cfg.AllowPatterns != nil {
		t.Error("AllowPatterns should default to nil")
	}
	if cfg.AllowPackages != nil {
		t.Error("AllowPackages should default to nil")
	}
}
