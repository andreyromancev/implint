/*
 *    Copyright 2019 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var cache = map[string]int64{}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No directory provided")
		os.Exit(1)
	}
	parseDir(os.Args[1])
}

func parseDir(dir string) int64 {
	level := levelFromFile(filepath.Join(dir, "lint.go"))
	cache[dir] = level

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return -1
	}

	packages, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.ImportsOnly)
	fatal(err)
	for _, p := range packages {
		for fk, f := range p.Files {
			for _, i := range f.Imports {
				path := strings.Trim(i.Path.Value, "\"")
				path = strings.Replace(path, "github.com/insolar/insolar/", "", -1)
				if path == dir {
					continue
				}

				importLevel := levelFromPath(path)
				if importLevel == -1 || level == -1 {
					continue
				}

				if importLevel <= level {
					fmt.Println("Potential import loop detected!")
					fmt.Printf("Package         | \"%s\" (level %d)\n", dir, level)
					fmt.Printf("Imports package | \"%s\" (level %d)\n", path, importLevel)
					fmt.Printf("In file         | \"%s\"\n", filepath.Base(fk))
					fmt.Println()
				}
			}
		}
	}

	files, err := ioutil.ReadDir(dir)
	fatal(err)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		parseDir(filepath.Join(dir, f.Name()))
	}

	return level
}

func levelFromFile(path string) int64 {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return -1
	}

	lintFile, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	fatal(err)

	var level int64
	for _, c := range lintFile.Comments {
		level = levelFromText(c.Text())
		if level == -1 {
			continue
		}
	}

	return level
}

func levelFromText(text string) int64 {
	text = strings.Replace(text, " ", "", -1)
	text = strings.TrimSuffix(text, "\n")
	parts := strings.Split(text, ":")
	if len(parts) != 2 || parts[0] != "level" {
		return -1
	}
	level, err := strconv.ParseInt(parts[1], 0, 64)
	if err != nil {
		return -1
	}
	return level
}

func levelFromPath(path string) int64 {
	level, ok := cache[path]
	if !ok {
		return parseDir(path)
	}

	return level
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
