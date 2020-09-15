module github.com/hedzr/consul-tags

go 1.12

// replace github.com/hedzr/common v0.0.0 => ../common

// replace github.com/hedzr/cmdr v0.0.0 => ../cmdr
// replace github.com/hedzr/cmdr v0.2.25 => ../cmdr

require (
	github.com/hashicorp/consul/api v1.1.0
	github.com/hashicorp/go-rootcerts v1.0.1 // indirect
	github.com/hedzr/cmdr v1.7.11
	github.com/hedzr/logex v1.2.9
	gopkg.in/hedzr/errors.v2 v2.1.0
	gopkg.in/yaml.v2 v2.3.0
)
