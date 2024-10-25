/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	cmdrv2 "github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/pkg/logz"

	"github.com/hedzr/consul-tags/util"
)

func getConnectionFromFlags(prefix string) (client *api.Client, bkup *kvJSON, err error) {
	// Start with the default Consul API config
	config := api.DefaultConfig()

	// Create a TLS config to be populated with flag-defined certs if applicable
	tlsConf := &tls.Config{}

	cs := cmdrv2.CmdStore().WithPrefix(prefix)

	// Set scheme and address:port
	config.Scheme = cs.MustString("scheme")
	// config.Address = fmt.Sprintf("%s:%v", c.GlobalString("app.ms.addr"), c.GlobalInt("app.ms.port"))
	// config.Scheme = cmdr.GetString("app.ms.scheme")
	config.Address = cs.MustString("addr")
	if !strings.Contains(config.Address, ":") {
		config.Address = fmt.Sprintf("%s:%v", config.Address, cs.MustInt("port"))
	}
	logz.Debug("Connecting to consul agent", "sheme", config.Scheme, "addr", config.Address)

	// Populate backup metadata
	bkup = &kvJSON{
		BackupDate: time.Now(),
		Connection: map[string]string{},
	}

	// Check for insecure flag
	if cs.MustBool("insecure") {
		tlsConf.InsecureSkipVerify = true
		bkup.Connection["insecure"] = "true"
	}

	// Load default system root CAs
	// ignore errors since the TLS config
	// will only be applied if --cert and --key
	// are defined
	tlsConf.ClientCAs, _ = util.LoadSystemRootCAs()

	// If --cert and --key are defined, load them and apply the TLS config
	if len(cs.MustString("cert")) > 0 && len(cs.MustString("key")) > 0 {
		// Make sure scheme is HTTPS when certs are used, regardless of the flag
		config.Scheme = "https"
		bkup.Connection["cert"] = cs.MustString("cert")
		bkup.Connection["key"] = cs.MustString("key")

		// Load cert and key files
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(cs.MustString("cert"), cs.MustString("key"))
		if err != nil {
			logz.Fatal("Could not load cert:", "err", err)
		}
		tlsConf.Certificates = append(tlsConf.Certificates, cert)

		// If cacert is defined, add it to the cert pool
		// else just use system roots
		if len(cs.MustString("cacert")) > 0 {
			tlsConf.ClientCAs = util.AddCACert(cs.MustString("cacert"), tlsConf.ClientCAs)
			tlsConf.RootCAs = tlsConf.ClientCAs
			bkup.Connection["cacert"] = cs.MustString("cacert")
		}
	}

	bkup.Connection["host"] = config.Scheme + "://" + config.Address

	if config.Scheme == "https" {
		// Set Consul's transport to the TLS config
		config.HttpClient.Transport = &http.Transport{
			TLSClientConfig: tlsConf,
		}
	}

	// Check for HTTP auth flags
	if len(cs.MustString("username")) > 0 && len(cs.MustString("password")) > 0 {
		config.HttpAuth = &api.HttpBasicAuth{
			Username: cs.MustString("username"),
			Password: cs.MustString("password"),
		}
		bkup.Connection["user"] = cs.MustString("username")
		bkup.Connection["pass"] = cs.MustString("password")
	}

	// Generate and return the API client
	client, err = api.NewClient(config)
	if err != nil {
		logz.Fatal("Error: %v", err)
		fmt.Println("Failed!")
	} else {
		fmt.Println("successfully")
	}
	return
}
