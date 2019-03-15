package implint

import (
	"errors"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type Config struct {
	Levels [][]string
}

func NewConfigYML(filename string) (conf Config, err error) {
	file, err := ioutil.ReadFile(filename)
	if err == io.EOF {
		err = errors.New("config not found")
	}
	if err != nil {
		return
	}

	err = yaml.Unmarshal(file, &conf)
	return
}
