package consul

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/consul/api"
	"gopkg.in/urfave/cli.v2"
	"hedzr.com/consul-tags/util"
	"net/http"
	"time"
	//"log"
	log "github.com/cihub/seelog"
)

func getConnectionFromFlags(c *cli.Context) (client *api.Client, bkup *kvJSON, err error) {
	// Start with the default Consul API config
	config := api.DefaultConfig()

	// Create a TLS config to be populated with flag-defined certs if applicable
	tlsConf := &tls.Config{}

	// Set scheme and address:port
	config.Scheme = c.String("scheme")
	//config.Address = fmt.Sprintf("%s:%v", c.GlobalString("addr"), c.GlobalInt("port"))
	config.Scheme = c.String("scheme")
	config.Address = c.String("addr")
	//if config.Address == "" {
	//	config.Address = c.GlobalString("consul.addr")
	//}
	log.Debugf("Connecting to %s://%s ...", config.Scheme, config.Address)

	// Populate backup metadata
	bkup = &kvJSON{
		BackupDate: time.Now(),
		Connection: map[string]string{},
	}

	// Check for insecure flag
	if c.Bool("insecure") {
		tlsConf.InsecureSkipVerify = true
		bkup.Connection["insecure"] = "true"
	}

	// Load default system root CAs
	// ignore errors since the TLS config
	// will only be applied if --cert and --key
	// are defined
	tlsConf.ClientCAs, _ = util.LoadSystemRootCAs()

	// If --cert and --key are defined, load them and apply the TLS config
	if len(c.String("cert")) > 0 && len(c.String("key")) > 0 {
		// Make sure scheme is HTTPS when certs are used, regardless of the flag
		config.Scheme = "https"
		bkup.Connection["cert"] = c.String("cert")
		bkup.Connection["key"] = c.String("key")

		// Load cert and key files
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(c.String("cert"), c.String("key"))
		if err != nil {
			log.Criticalf("Could not load cert: %v", err)
		}
		tlsConf.Certificates = append(tlsConf.Certificates, cert)

		// If cacert is defined, add it to the cert pool
		// else just use system roots
		if len(c.String("cacert")) > 0 {
			tlsConf.ClientCAs = util.AddCACert(c.String("cacert"), tlsConf.ClientCAs)
			tlsConf.RootCAs = tlsConf.ClientCAs
			bkup.Connection["cacert"] = c.String("cacert")
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
	if len(c.String("user")) > 0 && len(c.String("pass")) > 0 {
		config.HttpAuth = &api.HttpBasicAuth{
			Username: c.String("user"),
			Password: c.String("pass"),
		}
		bkup.Connection["user"] = c.String("user")
		bkup.Connection["pass"] = c.String("pass")
	}

	// Generate and return the API client
	client, err = api.NewClient(config)
	if err != nil {
		log.Criticalf("Error: %v", err)
		fmt.Println("Failed!")
	} else {
		fmt.Println("successfully")
	}
	return client, bkup, nil
}
