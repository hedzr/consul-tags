/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hedzr/cmdr"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"unicode/utf8"
)

func Backup() (err error) {
	var (
		client       *api.Client
		backupResult *kvJSON
		output       string
		prefix       = "app.kv"
	)

	output = cmdr.GetStringP(prefix, "backup.output")
	if len(output) == 0 {
		logrus.Fatal("ERROR: need -o output-file")
		return errors.New("Need -o output-file")
	}

	// Get KV client
	client, backupResult, err = getConnectionFromFlags(prefix)
	if err != nil {
		return
	}

	logrus.Debugf("Connected: %v", client)
	kv := client.KV()

	// Dump all
	pairs, _, err := kv.List(cmdr.GetStringP(prefix, "prefix"), &api.QueryOptions{})
	if err != nil {
		logrus.Fatalf("ERROR: %v", err)
		return
	}
	bkup := map[string]valueEnc{}
	for _, p := range pairs {
		validUtf8 := utf8.Valid(p.Value)
		if validUtf8 {
			bkup[p.Key] = valueEnc{"", string(p.Value)}
		} else {
			sEnc := base64.StdEncoding.EncodeToString(p.Value)
			bkup[p.Key] = valueEnc{"base64", sEnc}
		}
	}
	backupResult.Values = bkup

	logrus.Debugf("Dumping to %s", output)
	// Send results to outfile (if defined) or stdout
	dumpOutput(output, backupResult)

	return
}

func dumpOutput(pathname string, bkup *kvJSON) {
	if len(pathname) > 0 {
		ext := path.Ext(pathname)
		// fmt.Printf("EXT: %s", ext)
		switch ext {
		case ".json":
			outBytes, err := json.Marshal(bkup)
			if err != nil {
				logrus.Fatalf("Error: %v", err)
			}
			if err = ioutil.WriteFile(pathname, outBytes, 0664); err != nil {
				logrus.Fatalf("Error: %v", err)
			}
		case ".yml", ".yaml":
			outBytes, err := yaml.Marshal(bkup)
			if err != nil {
				logrus.Fatalf("Error: %v", err)
			}
			if err = ioutil.WriteFile(pathname, outBytes, 0664); err != nil {
				logrus.Fatalf("Error: %v", err)
			}
		}
	} else {
		outBytes, err := json.MarshalIndent(bkup, "", "  ")
		if err != nil {
			logrus.Fatalf("Error: %v", err)
		}
		fmt.Printf("%s\n", string(outBytes))
	}
}
