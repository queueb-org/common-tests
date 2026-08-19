package tests

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

var (
	locationPrefix = ""
)

func goldenName(t testing.TB) string {
	name := strings.Replace(t.Name(), "/", "__", -1)
	return fmt.Sprintf("%stestdata/%s.golden", locationPrefix, name)
}

// Golden is a simple helper to run golden like tests.
// -update flag provides test contents update (init run).
func Golden(t testing.TB, contents []byte) string {
	name := goldenName(t)

	if *update {
		if err := os.WriteFile(name, contents, defaultFilePerm); err != nil {
			t.Fatalf("could not write golden file: %s\n", name)
		}

		t.SkipNow()
		return ""
	}

	contents, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("could not read golden file: %s\n", name)
	}
	return string(contents)
}
