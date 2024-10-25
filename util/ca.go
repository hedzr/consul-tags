/*
 * Copyright © 2019 Hedzr Yeh.
 */

package util

import (
	"crypto/x509"
	"io/ioutil"
	"os"

	"github.com/hedzr/cmdr/v2/pkg/logz"
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
		logz.Fatal("Could not open CA cert:", "err", err)
		return roots
	}

	fBytes, err := ioutil.ReadAll(f)
	if err != nil {
		logz.Fatal("Failed to read CA cert:", "err", err)
		return roots
	}

	if !roots.AppendCertsFromPEM(fBytes) {
		logz.Fatal("Could not add CA to CA pool:", "err", err)
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
				logz.Debug("Loaded Root CA", "name", fi.Name())
			}
		}
	}

	return
}
