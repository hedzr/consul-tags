module github.com/hedzr/consul-tags

go 1.12

// replace github.com/hedzr/common v0.0.0 => ../common

// replace github.com/hedzr/cmdr v0.0.0 => ../cmdr
// replace github.com/hedzr/cmdr v0.2.25 => ../cmdr

require (
	github.com/hashicorp/consul/api v1.8.1
	github.com/hedzr/cmdr v1.8.0
	github.com/hedzr/logex v1.3.13
	gopkg.in/hedzr/errors.v2 v2.1.3
	gopkg.in/yaml.v2 v2.4.0
)

// exclude golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529
// exclude golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
// exclude golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae
