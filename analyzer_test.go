package ptrstruct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/shuymn/ptrstruct"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "basic", "ok", "containers", "generics")
}

func TestAnalyzer_Suppress(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "suppress")
}

func TestAnalyzer_FileNolint(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "filenolint")
}

func TestAnalyzer_TypeBlock(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "typeblock")
}

func TestAnalyzer_Allow(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	a := ptrstruct.NewAnalyzer()
	if err := a.Flags.Set("allow-types", "time.Time"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, testdata, a, "allow")
}

func TestAnalyzer_Alias(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "alias")
}

func TestAnalyzer_Nested(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "nested")
}

func TestAnalyzer_OnePer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "oneper")
}

func TestAnalyzer_Generated(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "generated")
}

func TestAnalyzer_Embedded(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "embedded")
}

func TestAnalyzer_IgnoreTestsDefault(t *testing.T) {
	t.Parallel()

	// Default: IgnoreTests=false, so test files ARE checked.
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ptrstruct.Analyzer, "ignoretests")
}

func TestAnalyzer_IgnoreTestsEnabled(t *testing.T) {
	t.Parallel()

	// With IgnoreTests=true, _test.go files should be skipped.
	// skiptests/_test.go has a violation but no // want comment;
	// if the analyzer checked it, analysistest would fail.
	testdata := analysistest.TestData()
	a := ptrstruct.NewAnalyzer()
	if err := a.Flags.Set("ignore-tests", "true"); err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, testdata, a, "skiptests")
}
