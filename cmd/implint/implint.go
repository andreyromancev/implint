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
	"flag"
	"fmt"
	"github.com/andreyromancev/implint"
	"os"
	"path/filepath"
)

const config = ".implint.yml"
var dir string

func init() {
	const (
		dirDefault = ""
		dirUsage = "Directory for linting. Should contain config file. Defaults to working directory."
	)
	flag.StringVar(&dir, "dir", dirDefault, dirUsage)
}

func main() {
	flag.Parse()

	absDir, err := filepath.Abs(dir)
	fatal(err)

	conf, err := implint.NewConfigYML(filepath.Join(absDir, config))
	if err != nil {
		fatal(err)
	}

	parser := implint.NewParser(conf)
	errs := parser.Parse(absDir)
	for _, err := range errs {
		fmt.Println(err)
		fmt.Println()
	}
	if len(errs) != 0 {
		os.Exit(1)
	}
}

func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
