package parsers

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type driver struct {
	Config []string `yaml:"drivers"`
}

func ParseDriver(dir string) []string {
	var path strings.Builder
	path.WriteString(dir)
	path.WriteString("/")
	path.WriteString("configs/drivers.yaml")

	config, err := os.ReadFile(path.String())
	mapping := driver{}
	if err != nil {
		log.Println(err)

		return []string{}
	}

	err = yaml.Unmarshal(config, &mapping)
	if err != nil {
		log.Println(err)

		return []string{}
	}

	return mapping.Config
}
