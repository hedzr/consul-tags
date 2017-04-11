package consul

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	//"log"
	log "github.com/cihub/seelog"
	"github.com/hashicorp/consul/api"
)

func Restore(c *cli.Context) (err error) {
	// Get KV client
	client, _, err := getConnectionFromFlags(c)
	if err != nil {
		return err
	}
	kv := client.KV()

	// Get backup JSON from file
	bkup, err := readBackupFile(c.Args().First())
	if err != nil {
		return fmt.Errorf("Error getting data: %v", err)
	}

	// restore file contents
	var v string
	for k, ve := range bkup.Values {
		switch ve.Encoding {
		case "base64":
			vd, err := base64.StdEncoding.DecodeString(ve.Str)
			if err != nil {
				return fmt.Errorf("Error decoding the value of key '%s': %v", k, err)
			}
			v = string(vd)
		case "utf8", "":
			v = ve.Str
		default:
			return fmt.Errorf("Unknown encoding '%v' for key '%s'", ve.Encoding, k)
		}

		log.Debugf("Restoring key '%s'", k)
		if _, err := kv.Put(&api.KVPair{
			Key:   k,
			Value: []byte(v),
		}, &api.WriteOptions{}); err != nil {
			return fmt.Errorf("Error writing key %s: %v", k, err)
		}
	}
	return nil
}

func readBackupFile(pathname string) (bkup *kvJSON, err error) {
	var f *os.File
	f, err = os.Open(pathname)
	defer f.Close()
	if err != nil {
		return
	}
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
		// https://gopkg.in/yaml.v2
		err = yaml.Unmarshal(b, &bkup)
	}
	return
}
