/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/consul-tags/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func getConnectionFromFlags(prefix string) (client *api.Client, bkup *kvJSON, err error) {
	// Start with the default Consul API config
	config := api.DefaultConfig()

	// Create a TLS config to be populated with flag-defined certs if applicable
	tlsConf := &tls.Config{}

	// Set scheme and address:port
	config.Scheme = cmdr.GetStringP(prefix, "scheme")
	// config.Address = fmt.Sprintf("%s:%v", c.GlobalString("app.ms.addr"), c.GlobalInt("app.ms.port"))
	// config.Scheme = cmdr.GetString("app.ms.scheme")
	config.Address = cmdr.GetStringP(prefix, "addr")
	if !strings.Contains(config.Address, ":") {
		config.Address = fmt.Sprintf("%s:%v", config.Address, cmdr.GetIntP(prefix, "port"))
	}
	logrus.Debugf("Connecting to %s://%s ...", config.Scheme, config.Address)

	// Populate backup metadata
	bkup = &kvJSON{
		BackupDate: time.Now(),
		Connection: map[string]string{},
	}

	// Check for insecure flag
	if cmdr.GetBoolP(prefix, "insecure") {
		tlsConf.InsecureSkipVerify = true
		bkup.Connection["insecure"] = "true"
	}

	// Load default system root CAs
	// ignore errors since the TLS config
	// will only be applied if --cert and --key
	// are defined
	tlsConf.ClientCAs, _ = util.LoadSystemRootCAs()

	// If --cert and --key are defined, load them and apply the TLS config
	if len(cmdr.GetStringP(prefix, "cert")) > 0 && len(cmdr.GetStringP(prefix, "key")) > 0 {
		// Make sure scheme is HTTPS when certs are used, regardless of the flag
		config.Scheme = "https"
		bkup.Connection["cert"] = cmdr.GetStringP(prefix, "cert")
		bkup.Connection["key"] = cmdr.GetStringP(prefix, "key")

		// Load cert and key files
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(cmdr.GetStringP(prefix, "cert"), cmdr.GetStringP(prefix, "key"))
		if err != nil {
			logrus.Fatalf("Could not load cert: %v", err)
		}
		tlsConf.Certificates = append(tlsConf.Certificates, cert)

		// If cacert is defined, add it to the cert pool
		// else just use system roots
		if len(cmdr.GetStringP(prefix, "cacert")) > 0 {
			tlsConf.ClientCAs = util.AddCACert(cmdr.GetStringP(prefix, "cacert"), tlsConf.ClientCAs)
			tlsConf.RootCAs = tlsConf.ClientCAs
			bkup.Connection["cacert"] = cmdr.GetStringP(prefix, "cacert")
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
	if len(cmdr.GetStringP(prefix, "username")) > 0 && len(cmdr.GetStringP(prefix, "password")) > 0 {
		config.HttpAuth = &api.HttpBasicAuth{
			Username: cmdr.GetStringP(prefix, "username"),
			Password: cmdr.GetStringP(prefix, "password"),
		}
		bkup.Connection["user"] = cmdr.GetStringP(prefix, "username")
		bkup.Connection["pass"] = cmdr.GetStringP(prefix, "password")
	}

	// Generate and return the API client
	client, err = api.NewClient(config)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
		fmt.Println("Failed!")
	} else {
		fmt.Println("successfully")
	}
	return
}
