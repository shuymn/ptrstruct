package ptrstruct_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/shuymn/ptrstruct"
)

func TestAnalyzer_SuggestedFixes(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(
		t,
		testdata,
		ptrstruct.Analyzer,
		"basic",
		"containers",
		"nested",
		"alias",
		"embedded",
	)
}

func TestAnalyzer_SuggestedFixes_WithOptionalFlags(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	a := ptrstruct.NewAnalyzer()
	for _, name := range []string{"map-key", "array-elem", "chan-elem"} {
		if err := a.Flags.Set(name, "true"); err != nil {
			t.Fatalf("set flag %s: %v", name, err)
		}
	}
	analysistest.RunWithSuggestedFixes(t, testdata, a, "fixpaths")
}

func TestCLI_FixAndDiff(t *testing.T) {
	t.Parallel()

	repoRoot, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binDir := t.TempDir()
	binPath := filepath.Join(binDir, "ptrstruct")
	build := exec.Command("go", "build", "-o", binPath, "./cmd/ptrstruct")
	build.Dir = repoRoot
	if buildOut, buildErr := build.CombinedOutput(); buildErr != nil {
		t.Fatalf("build ptrstruct: %v\n%s", buildErr, buildOut)
	}

	workspace := t.TempDir()
	writeFile(t, filepath.Join(workspace, "go.mod"), "module example.com/cli\n\ngo 1.25.0\n")
	writeFile(t, filepath.Join(workspace, "cli.go"), `package cli

type User struct {
	Name string
}

func (u User) Normalize() {}

func Save(users []User) {}

type Response struct {
	Meta User
}
`)

	diff := exec.Command(binPath, "-fix", "-diff", "./...")
	diff.Dir = workspace
	diffOut, err := diff.CombinedOutput()
	if err != nil {
		t.Fatalf("ptrstruct -fix -diff: %v\n%s", err, diffOut)
	}

	diffText := string(diffOut)
	for _, needle := range []string{
		"-func (u User) Normalize() {}",
		"+func (u *User) Normalize() {}",
		"-func Save(users []User) {}",
		"+func Save(users []*User) {}",
		"-\tMeta User",
		"+\tMeta *User",
	} {
		if !strings.Contains(diffText, needle) {
			t.Fatalf("diff output missing %q\n%s", needle, diffText)
		}
	}

	afterDiff, err := os.ReadFile(filepath.Join(workspace, "cli.go"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(afterDiff), "*User") {
		t.Fatalf("-diff should not rewrite files\n%s", afterDiff)
	}

	fix := exec.Command(binPath, "-fix", "./...")
	fix.Dir = workspace
	if fixOut, fixErr := fix.CombinedOutput(); fixErr != nil {
		t.Fatalf("ptrstruct -fix: %v\n%s", fixErr, fixOut)
	}

	got, err := os.ReadFile(filepath.Join(workspace, "cli.go"))
	if err != nil {
		t.Fatal(err)
	}

	want := `package cli

type User struct {
	Name string
}

func (u *User) Normalize() {}

func Save(users []*User) {}

type Response struct {
	Meta *User
}
`
	if string(got) != want {
		t.Fatalf("unexpected rewritten file\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
