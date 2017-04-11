package consul

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"unicode/utf8"
	//"log"
	log "github.com/cihub/seelog"
	"github.com/hashicorp/consul/api"
)

func Backup(c *cli.Context) (err error) {
	if c.String("outfile") == "" {
		log.Critical("ERROR: need -o outfile")
		return errors.New("Need -o outfile")
	}

	// Get KV client
	client, backupResult, err := getConnectionFromFlags(c)
	if err != nil {
		return err
	}

	log.Infof("Connected: %v", client)
	kv := client.KV()

	// Dump all
	pairs, _, err := kv.List(c.String("prefix"), &api.QueryOptions{})
	if err != nil {
		log.Criticalf("ERROR: %v", err)
		return err
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

	log.Debugf("Dumping to %s", c.String("outfile"))
	// Send results to outfile (if defined) or stdout
	dumpOutput(c.String("outfile"), backupResult)

	return nil
}

func dumpOutput(pathname string, bkup *kvJSON) {
	if len(pathname) > 0 {
		ext := path.Ext(pathname)
		//fmt.Printf("EXT: %s", ext)
		switch ext {
		case ".json":
			outBytes, err := json.Marshal(bkup)
			if err != nil {
				log.Criticalf("Error: %v", err)
			}
			if err = ioutil.WriteFile(pathname, outBytes, 0664); err != nil {
				log.Criticalf("Error: %v", err)
			}
		case ".yml", ".yaml":
			outBytes, err := yaml.Marshal(bkup)
			if err != nil {
				log.Criticalf("Error: %v", err)
			}
			if err = ioutil.WriteFile(pathname, outBytes, 0664); err != nil {
				log.Criticalf("Error: %v", err)
			}
		}
	} else {
		outBytes, err := json.MarshalIndent(bkup, "", "  ")
		if err != nil {
			log.Criticalf("Error: %v", err)
		}
		fmt.Printf("%s\n", string(outBytes))
	}
}
