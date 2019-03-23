package implint

import (
	"path/filepath"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParser(Config{
		Levels: [][]string{
			{"worker", "handler"},
			{"fetcher", "saver"},
			{"db/io", "http"},
		},
	})
	absDir, err := filepath.Abs("example")
	if err != nil {
		t.FailNow()
	}
	errs := parser.Parse(absDir)

	if len(errs) != 4 {
		t.Error("Wrong error count")
		t.FailNow()
	}
	assertErr(t,
		ErrNoLevel(filepath.Join(absDir, "db")),
		errs[0],
	)
	assertErr(t,
		ErrBadImport{
			from: filepath.Join(absDir, "db/io"),
			imp: filepath.Join(absDir, "http"),
			file: filepath.Join(absDir, "db/io/io.go"),
			fromLevel: 2,
			impLevel: 2,
		},
		errs[1],
	)
	assertErr(t,
		ErrBadImport{
			from: filepath.Join(absDir, "http"),
			imp: absDir,
			file: filepath.Join(absDir, "http/http.go"),
			fromLevel: 2,
			impLevel: 0,
		},
		errs[2],
	)
	assertErr(t,
		ErrBadImport{
			from: filepath.Join(absDir, "saver"),
			imp: filepath.Join(absDir, "fetcher"),
			file: filepath.Join(absDir, "saver/saver.go"),
			fromLevel: 1,
			impLevel: 1,
		},
		errs[3],
	)
}

func assertErr(t *testing.T, expected, actual error) {
	if expected != actual {
		t.Errorf("Wrong error.\nExpected: \"%s\"\nActual: \"%s\"", expected, actual)
		t.FailNow()
	}
}
