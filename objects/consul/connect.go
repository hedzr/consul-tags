/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hedzr/consul-tags/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func getConnectionFromFlags() (client *api.Client, bkup *kvJSON, err error) {
	// Start with the default Consul API config
	config := api.DefaultConfig()

	// Create a TLS config to be populated with flag-defined certs if applicable
	tlsConf := &tls.Config{}

	// Set scheme and address:port
	config.Scheme = viper.GetString("app.ms.scheme")
	// config.Address = fmt.Sprintf("%s:%v", c.GlobalString("app.ms.addr"), c.GlobalInt("app.ms.port"))
	// config.Scheme = viper.GetString("app.ms.scheme")
	config.Address = viper.GetString("app.ms.addr")
	// if config.Address == "" {
	// 	config.Address = c.GlobalString("consul.addr")
	// }
	logrus.Debugf("Connecting to %s://%s ...", config.Scheme, config.Address)

	// Populate backup metadata
	bkup = &kvJSON{
		BackupDate: time.Now(),
		Connection: map[string]string{},
	}

	// Check for insecure flag
	if viper.GetBool("app.ms.insecure") {
		tlsConf.InsecureSkipVerify = true
		bkup.Connection["insecure"] = "true"
	}

	// Load default system root CAs
	// ignore errors since the TLS config
	// will only be applied if --cert and --key
	// are defined
	tlsConf.ClientCAs, _ = util.LoadSystemRootCAs()

	// If --cert and --key are defined, load them and apply the TLS config
	if len(viper.GetString("app.ms.cert")) > 0 && len(viper.GetString("app.ms.key")) > 0 {
		// Make sure scheme is HTTPS when certs are used, regardless of the flag
		config.Scheme = "https"
		bkup.Connection["cert"] = viper.GetString("app.ms.cert")
		bkup.Connection["key"] = viper.GetString("app.ms.key")

		// Load cert and key files
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(viper.GetString("app.ms.cert"), viper.GetString("app.ms.key"))
		if err != nil {
			logrus.Fatalf("Could not load cert: %v", err)
		}
		tlsConf.Certificates = append(tlsConf.Certificates, cert)

		// If cacert is defined, add it to the cert pool
		// else just use system roots
		if len(viper.GetString("app.ms.cacert")) > 0 {
			tlsConf.ClientCAs = util.AddCACert(viper.GetString("app.ms.cacert"), tlsConf.ClientCAs)
			tlsConf.RootCAs = tlsConf.ClientCAs
			bkup.Connection["cacert"] = viper.GetString("app.ms.cacert")
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
	if len(viper.GetString("app.ms.user")) > 0 && len(viper.GetString("app.ms.pass")) > 0 {
		config.HttpAuth = &api.HttpBasicAuth{
			Username: viper.GetString("app.ms.user"),
			Password: viper.GetString("app.ms.pass"),
		}
		bkup.Connection["user"] = viper.GetString("app.ms.user")
		bkup.Connection["pass"] = viper.GetString("app.ms.pass")
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
