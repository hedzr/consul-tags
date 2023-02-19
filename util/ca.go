/*
 * Copyright © 2019 Hedzr Yeh.
 */

package util

import (
	"crypto/x509"
	"io/ioutil"
	"os"

	"github.com/hedzr/cmdr"
)

var certDirectories = []string{
	"/system/etc/security/cacerts",     // Android
	"/usr/local/share/ca-certificates", // Debian derivatives
	"/etc/pki/ca-trust/source/anchors", // RedHat derivatives
	"/etc/ca-certificates",             // Misc alternatives
	"/usr/share/ca-certificates",       // Misc alternatives
}

func AddCACert(path string, roots *x509.CertPool) *x509.CertPool {
	f, err := os.Open(path)
	if err != nil {
		cmdr.Logger.Fatalf("Could not open CA cert: %v", err)
		return roots
	}

	fBytes, err := ioutil.ReadAll(f)
	if err != nil {
		cmdr.Logger.Fatalf("Failed to read CA cert: %v", err)
		return roots
	}

	if !roots.AppendCertsFromPEM(fBytes) {
		cmdr.Logger.Fatalf("Could not add CA to CA pool: %v", err)
	}
	return roots
}

func LoadSystemRootCAs() (systemRoots *x509.CertPool, err error) {
	systemRoots = x509.NewCertPool()

	for _, directory := range certDirectories {
		fis, err := ioutil.ReadDir(directory)
		if err != nil {
			continue
		}
		for _, fi := range fis {
			data, err := ioutil.ReadFile(directory + "/" + fi.Name())
			if err == nil && systemRoots.AppendCertsFromPEM(data) {
				cmdr.Logger.Debugf("Loaded Root CA %s", fi.Name())
			}
		}
	}

	return
}
