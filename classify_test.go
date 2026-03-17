package ptrstruct

import (
	"go/types"
	"testing"
)

func TestNewClassifier_InvalidPattern(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()
	cfg.AllowPatterns = []string{"[invalid"}
	_, err := NewClassifier(cfg)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestClassifier_IsAllowed_ByType(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()
	cfg.AllowTypes = []string{"example.com/app.User"}
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	if !cls.IsAllowed(user) {
		t.Error("User should be allowed by type")
	}

	profile := newNamedStruct("example.com/app", "app", "Profile",
		types.NewVar(0, nil, "Bio", types.Typ[types.String]),
	)
	if cls.IsAllowed(profile) {
		t.Error("Profile should not be allowed")
	}
}

func TestClassifier_IsAllowed_ByPackage(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()
	cfg.AllowPackages = []string{"example.com/external"}
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	ext := newNamedStruct("example.com/external", "external", "Foo",
		types.NewVar(0, nil, "V", types.Typ[types.Int]),
	)
	if !cls.IsAllowed(ext) {
		t.Error("Foo from allowed package should be allowed")
	}

	internal := newNamedStruct("example.com/internal", "internal", "Bar",
		types.NewVar(0, nil, "V", types.Typ[types.Int]),
	)
	if cls.IsAllowed(internal) {
		t.Error("Bar from non-allowed package should not be allowed")
	}
}

func TestClassifier_IsAllowed_ByPattern(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()
	cfg.AllowPatterns = []string{`\.Null[A-Z]\w*$`}
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	nullStr := newNamedStruct("database/sql", "sql", "NullString",
		types.NewVar(0, nil, "String", types.Typ[types.String]),
		types.NewVar(0, nil, "Valid", types.Typ[types.Bool]),
	)
	if !cls.IsAllowed(nullStr) {
		t.Error("NullString should match pattern")
	}

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	if cls.IsAllowed(user) {
		t.Error("User should not match pattern")
	}
}

func TestClassifier_EmptyConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()
	cls, err := NewClassifier(cfg)
	if err != nil {
		t.Fatal(err)
	}

	user := newNamedStruct("example.com/app", "app", "User",
		types.NewVar(0, nil, "Name", types.Typ[types.String]),
	)
	if cls.IsAllowed(user) {
		t.Error("nothing should be allowed with empty config")
	}
}
