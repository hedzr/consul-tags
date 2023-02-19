/*
 * Copyright © 2019 Hedzr Yeh.
 */

package objects

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config map[string]interface{}

var (
	Configurations Config
)

func Has(key string) bool {
	if _, ok := Configurations[key]; ok {
		// do something here
		return true
	}
	return false
}

func Get(parent *Config, key string) *Config {
	if parent == nil {
		parent = &Configurations
	}
	if val, ok := (*parent)[key]; ok {
		// return val.(map[string]interface{})
		// fmt.Printf("val = %v\n", val)

		// z := val.(Config)
		z := make(Config)
		j, err := yaml.Marshal(&val)
		if err != nil {
			fmt.Println(err)
		}
		err = yaml.Unmarshal(j, &z)
		if err != nil {
			fmt.Println(err)
		}

		// fmt.Printf("z = %v\n", z)
		return &z
	}
	return nil
}

func Set(parent *Config, key string, value interface{}) {
	if parent == nil {
		parent = &Configurations
	}
	(*parent)[key] = value
}

func GetAs(container interface{}, parent *Config, key string) *Config {
	if parent == nil {
		parent = &Configurations
	}
	if val, ok := (*parent)[key]; ok {
		// fmt.Printf("   val = %v\n", val)
		j, err := yaml.Marshal(&val)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Printf("   j = %v\n", j)
		err = yaml.Unmarshal(j, container)
		if err != nil {
			fmt.Println(err)
		}

		// return val.(map[string]interface{})
		// z := val.(Config)
		z := make(Config)
		err = yaml.Unmarshal(j, &z)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Printf("   z: %v\n", z)
		return &z
	}

	return nil
}
