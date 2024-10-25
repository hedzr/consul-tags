/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"unicode/utf8"

	"github.com/hashicorp/consul/api"
	cmdrv2 "github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"gopkg.in/hedzr/errors.v3"
	"gopkg.in/yaml.v3"
)

func Backup() (err error) {
	var (
		client       *api.Client
		backupResult *kvJSON
		output       string
		prefixV2     = "kv"
	)

	cs := cmdrv2.CmdStore().WithPrefix(prefixV2)
	output = cs.MustString("backup.output")
	if len(output) == 0 {
		logz.Error("ERROR: need -o output-file")
		return errors.New("Need -o output-file")
	}

	// Get KV client
	client, backupResult, err = getConnectionFromFlags(prefixV2)
	if err != nil {
		return
	}

	logz.Debug("Connected: ", "client", client)
	kv := client.KV()

	// Dump all
	pairs, _, err := kv.List(cs.MustString("prefix"), &api.QueryOptions{})
	if err != nil {
		logz.Fatal("ERROR:", "err", err)
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

	logz.Debug("Dumping to file", "output", output)
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
				logz.Fatal("Error:", "err", err)
			}
			if err = os.WriteFile(pathname, outBytes, 0664); err != nil {
				logz.Fatal("Error:", "err", err)
			}
		case ".yml", ".yaml":
			outBytes, err := yaml.Marshal(bkup)
			if err != nil {
				logz.Fatal("Error:", "err", err)
			}
			if err = os.WriteFile(pathname, outBytes, 0664); err != nil {
				logz.Fatal("Error:", "err", err)
			}
		}
	} else {
		outBytes, err := json.MarshalIndent(bkup, "", "  ")
		if err != nil {
			logz.Fatal("Error:", "err", err)
		}
		fmt.Printf("%s\n", string(outBytes))
	}
}
