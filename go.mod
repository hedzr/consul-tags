module github.com/hedzr/consul-tags

go 1.12

replace github.com/hedzr/common v0.0.0 => ../common

require (
	github.com/hashicorp/consul v1.4.4
	github.com/hedzr/cmdr v0.2.1
	github.com/hedzr/common v0.0.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	github.com/takama/daemon v0.0.0-20180403113744-aa76b0035d12
	gopkg.in/urfave/cli.v2 v2.0.0-20180128182452-d3ae77c26ac8
	gopkg.in/yaml.v2 v2.2.2
)
