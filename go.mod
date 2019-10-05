module github.com/hedzr/consul-tags

go 1.12

// replace github.com/hedzr/common v0.0.0 => ../common

// replace github.com/hedzr/cmdr v0.0.0 => ../cmdr
// replace github.com/hedzr/cmdr v0.2.25 => ../cmdr

require (
	github.com/hashicorp/consul/api v1.1.0
	github.com/hashicorp/go-rootcerts v1.0.1 // indirect
	github.com/hedzr/cmdr v1.5.3
	github.com/hedzr/logex v1.0.3
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/yaml.v2 v2.2.2
)
