package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/kael-k/pongo/v2/pongo"
)

// this test file aims to expose all tests file in tests/assets
// and to validate the schemas inside it

const testRootBrokenPongoSchemas = "../tests/assets/broken-pongo-schema"

func TestBrokenPongoSchemas(t *testing.T) {
	testsDir, err := os.ReadDir(testRootBrokenPongoSchemas)

	if err != nil {
		t.Errorf("cannot test broken schemas, error on read %s: %s", testRootBrokenPongoSchemas, err)
		return
	}

	for _, testFile := range testsDir {
		filename := testFile.Name()
		if filename == "README.md" {
			continue
		}

		path := fmt.Sprintf("%s/%s", testRootBrokenPongoSchemas, testFile.Name())
		rawPonGOJSONSchema, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("cannot test broken schemas, error on read %s: %s", path, err)
			return
		}

		_, _, err = pongo.UnmarshalPongoSchema(rawPonGOJSONSchema)
		if err == nil {
			t.Errorf("expected error on path %s, got no error", path)
		}
	}
}
