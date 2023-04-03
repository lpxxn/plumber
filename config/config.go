package config

import (
	"io/ioutil"
	"os"

	"github.com/lpxxn/plumber/src/log"
	"gopkg.in/yaml.v3"
)

func readFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func ReadFile(filename string, v interface{}) error {
	body, err := readFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(body, v)
	log.Debugf("config: %v", v)
	return nil
}
