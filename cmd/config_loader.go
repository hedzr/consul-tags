package cmd

import (
	"fmt"
	"os"
	"path"
	//"time"

	"io/ioutil"
	"net/http"
	"net/url"

	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"

	"gopkg.in/yaml.v2"
)

type yamlSourceContext struct {
	FilePath string
}

//
//// InputSourceContext is an interface used to allow
//// other input sources to be implemented as needed.
//type InputSourceContext interface {
//	Int(name string) (int, error)
//	Duration(name string) (time.Duration, error)
//	Float64(name string) (float64, error)
//	String(name string) (string, error)
//	StringSlice(name string) ([]string, error)
//	IntSlice(name string) ([]int, error)
//	Generic(name string) (cli.Generic, error)
//	Bool(name string) (bool, error)
//	BoolT(name string) (bool, error)
//}

// NewYamlSourceFromFile creates a new Yaml InputSourceContext from a filepath.
func NewYamlSourceFromFile(file string) (altsrc.InputSourceContext, error) {
	var results map[interface{}]interface{}
	ysc := &yamlSourceContext{FilePath: file}
	if _, err := os.Stat(ysc.FilePath); os.IsNotExist(err) {
		home := path.Dir(os.Getenv("HOME"))
		path := path.Join(home, file)
		ysc = &yamlSourceContext{FilePath: path}
	}
	err := readCommandYaml(ysc.FilePath, &results)
	if err != nil {
		return &MapInputSource{}, nil // fmt.Errorf("Skip loading configuration fiile '%s': '%v'", ysc.FilePath, err.Error())
	}
	return &MapInputSource{valueMap: results}, nil
}

// NewYamlSourceFromFlagFunc creates a new Yaml InputSourceContext from a provided flag name and source context.
func NewYamlSourceFromFlagFunc(flagFileName string) func(context *cli.Context) (altsrc.InputSourceContext, error) {
	return func(context *cli.Context) (altsrc.InputSourceContext, error) {
		filePath := context.String(flagFileName)
		return NewYamlSourceFromFile(filePath)
	}
}

func readCommandYaml(filePath string, container interface{}) (err error) {
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

func loadDataFrom(filePath string) ([]byte, error) {
	u, err := url.Parse(filePath)
	if err != nil {
		return nil, err
	}

	if u.Host != "" {
		// i have a host, now do i support the scheme?
		switch u.Scheme {
		case "http", "https":
			res, err := http.Get(filePath)
			if err != nil {
				return nil, err
			}
			return ioutil.ReadAll(res.Body)
		default:
			return nil, fmt.Errorf("scheme of %s is unsupported", filePath)
		}
	} else if u.Path != "" {
		// i dont have a host, but I have a path. I am a local file.
		if _, notFoundFileErr := os.Stat(filePath); notFoundFileErr != nil {
			return nil, fmt.Errorf("Cannot read from file: '%s' because it does not exist.", filePath)
		}
		return ioutil.ReadFile(filePath)
	} else {
		return nil, fmt.Errorf("unable to determine how to load from path %s", filePath)
	}
}
