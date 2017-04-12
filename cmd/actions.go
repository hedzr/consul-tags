package cmd

import (
	"fmt"
	"github.com/hedzr/consul-tags/objects"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

var DefaultAction = func(c *cli.Context) error {
	name := "Guy"
	if c.NArg() > 0 {
		name = c.Args().Get(0)
	}

	if c.Bool("ginger-crouton") {
		return cli.Exit("it is not in the soup", 86)
	}

	//if c.String("lang") == "spanish" {
	if language == "spanish" {
		fmt.Println("Hola", name)
	} else {
		fmt.Printf("Hello %q, type '%s --help' for command line usage.\n", name, c.App.Name)
	}
	return nil
}

func LoadConfigFile(c *cli.Context) {
	filePath := c.String("config")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		home := path.Dir(os.Getenv("HOME"))
		filePath = path.Join(home, filePath)
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		readConfigYaml(filePath, &objects.Configurations)
	}

	//fmt.Printf("HASH: %v\n\n", objects.Configurations["gitlab-cli"])
}

func readConfigYaml(filePath string, container *objects.Config) (err error) {
	b, err := loadDataFrom(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, container)
	if err != nil {
		return err
	}

	err = nil
	return
}
