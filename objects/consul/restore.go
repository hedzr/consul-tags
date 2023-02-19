/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/consul/api"
	"github.com/hedzr/cmdr"
	"gopkg.in/yaml.v3"
)

func Restore() (err error) {
	var (
		client *api.Client
		bkup   *kvJSON
		v      string
		prefix = "app.kv"
	)
	// Get KV client
	client, bkup, err = getConnectionFromFlags(prefix)
	if err != nil {
		return
	}

	kv := client.KV()

	// Get backup JSON from file
	bkup, err = readBackupFile(cmdr.GetStringP(prefix, "input"))
	if err != nil {
		return fmt.Errorf("error getting data: %v", err)
	}

	// restore file contents
	for k, ve := range bkup.Values {
		switch ve.Encoding {
		case "base64":
			vd, err := base64.StdEncoding.DecodeString(ve.Str)
			if err != nil {
				return fmt.Errorf("error decoding the value of key '%s': %v", k, err)
			}
			v = string(vd)
		case "utf8", "":
			v = ve.Str
		default:
			return fmt.Errorf("unknown encoding '%v' for key '%s'", ve.Encoding, k)
		}

		cmdr.Logger.Debugf("Restoring key '%s'", k)
		if _, err := kv.Put(&api.KVPair{
			Key:   k,
			Value: []byte(v),
		}, &api.WriteOptions{}); err != nil {
			return fmt.Errorf("error writing key %s: %v", k, err)
		}
	}
	return nil
}

func readBackupFile(pathname string) (bkup *kvJSON, err error) {
	var f *os.File
	f, err = os.Open(pathname)
	if err != nil {
		return
	}

	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	ext := path.Ext(pathname)
	switch ext {
	case ".json":
		err = json.Unmarshal(b, &bkup)
	case ".yml", ".yaml":
		// https://github.com/go-yaml/yaml
		// https://gopkg.in/yaml.v3
		err = yaml.Unmarshal(b, &bkup)
	}
	return
}
