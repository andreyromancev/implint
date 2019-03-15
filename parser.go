package implint

import (
	"errors"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Parser struct {
	config Config
	levels map[string]int
}

func NewParser(config Config) *Parser {
	return &Parser{
		config: config,
		levels: map[string]int{},
	}
}

func (pr *Parser) Parse(dir string) (errs []error) {
	for i, paths := range pr.config.Levels {
		for _, p := range paths {
			abs, err := filepath.Abs(filepath.Join(dir, p))
			if err != nil {
				errs = append(errs, fmt.Errorf("no such package: %s", p))
			}
			pr.levels[abs] = i
		}
	}
	pr.levels[dir] = 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if path == dir {
			return nil
		}
		errs = append(errs, pr.parseDir(path)...)
		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}

	return
}

func (pr *Parser) parseDir(dir string) (errs []error) {
	level, ok := pr.levels[dir]
	if !ok {
		errs = append(errs, ErrNoLevel(dir))
		return
	}

	packages, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.ImportsOnly | parser.ParseComments)
	if err != nil {
		errs = append(errs, errors.New(fmt.Sprintf("failed to parse \"%s\" [%v]", dir, err)))
		return
	}

	for _, p := range packages {
		for filePath, f := range p.Files {
			for _, i := range f.Imports {
				importDir := import2Dir(i.Path.Value)
				importLevel, ok := pr.levels[importDir]
				if !ok {
					continue
				}

				if importLevel <= level {
					errs = append(errs, ErrBadImport{
						from: dir,
						imp: importDir,
						fromLevel: level,
						impLevel: importLevel,
						file: filePath,
					})
				}
			}
		}
	}

	return
}

func goPath() (p string) {
	p = os.Getenv("GOPATH")
	if p == "" {
		p = build.Default.GOPATH
	}
	return
}

func import2Dir(imp string) (p string) {
	p = strings.Trim(imp, `"`)
	p = filepath.Join(goPath(), "src", p)
	return
}
