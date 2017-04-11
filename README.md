# Consul Service Tags Modifier

[![Build Status](https://travis-ci.org/hedzr/consul-tags.svg?branch=master)](https://travis-ci.org/hedzr/consul-tags)

`consul-tags` can add, remove tag(s) of a consul service (one or all its instances).

Here is a first release for key functionality.


## Install

### Binary

Download binary from [Release](../../releases/latest) page.

### Docker Hub

WAIT A MINUTE.

### Go Build

clone the repo and build:

```bash
go build -o consul-tags github.com/hedzr/consul-tags
```

mixin:

```bash
go get -u github.com/hedzr/consul-tags/objects/consul
```


## Usage


```bash
# list tags
consul-tags ms --name test-redis tags ls
# add tags
consul-tags ms --name test-redis tags --add a,c,e
consul-tags ms --name test-redis tags --add a --add c --add e
# remove tags
consul-tags ms --name test-redis tags --rm a,c,e
consul-tags ms --name test-redis tags --rm a --rm c --rm e
# by id
consul-tags ms --id test-redis-1 tags ls

# toggle master/slave
consul-tags ms --name test-redis tags toggle --addr 10.7.13.1:6379 --set role=master --reset role=slave
consul-tags ms --name test-mq tags toggle --addr 10.7.16.3:5672 --set role=leader,type=ram --reset role=peer,type=disk

# get commands and sub-commands help
consul-tags ms -h
# get help: -h or --help
consul-tags -h
```

Default consul address is `consul.ops.local:8500`, can be overriden with environment variable `CONSUL_ADDR` (host:port), too. Such as:

```
CONSUL_ADDR=127.0.0.1:8500 consul-tags ms --name=consul tags ls
# or
export CONSUL_HOST=ns2.company.domain:8500
consul-tags --help
```

### COMMANDS:

```
     kv           K/V pair Operations, ...
     ms, service  Microservice Operations, ...
     help, h      Shows a list of commands or help for one command
```

### GLOBAL OPTIONS:

```
   --addr HOST[:PORT], -a HOST[:PORT]  Consul address and port: HOST[:PORT] (No leading 'http(s)://') (default: "consul.ops.local") [$CONSUL_ADDR]
   --port value, -p value              Consul port (default: 8500) [$CONSUL_PORT]
   --prefix value                      Root key prefix (default: "/") [$CONSUL_PREFIX]
   --cacert value, -r value            Client CA cert [$CONSUL_CA_CERT]
   --cert value, -t value              Client cert [$CONSUL_CERT]
   --scheme value, -s value            Consul connection scheme (HTTP or HTTPS) (default: "http") [$CONSUL_SCHEME]
   --insecure, -K                      Skip TLS host verification (default: false) [$CONSUL_INSECURE]
   --username value, -U value          HTTP Basic auth user [$CONSUL_USER]
   --password value, -P value          HTTP Basic auth password [$CONSUL_PASS]
   --key value, -Y value               Client key [$CONSUL_KEY]
   --help, -h                          show help (default: false)
   --init-completion value             generate completion code. Value must be 'bash' or 'zsh'
   --version, -v                       print the version (default: false)
```

## Shell completion

with bonus from `urfave/cli`, you can install auto-completion for this binary:

```bash
eval "`consul-tags --init-completion bash`"
consul-tags --init-completion bash >> ~/.bashrc
```

or put it to `/etc/bash_completion.d/`:

```bash
cp bin/consul-tags /usr/local/bin/
/usr/local/bin/consul-tags --init-completion bash | sudo tee /etc/bash_completion.d/consul-tags
source /etc/bash_completion.d/consul-tags
```

see also [Shell Completion](https://github.com/urfave/cli/tree/v2#shell-completion).

zsh: TODO
Mac: TODO


## Thanks

my plan is building a suite of devops. consul operations is part of it. exception consul tags modifying operation, these codes ported in:

- [colebrumley/consul-kv-backup](https://github.com/colebrumley/consul-kv-backup)

and some repositories are good:

- https://github.com/shreyu86/consul-backup
- https://github.com/kailunshi/consul-backup

### Third-party repos

Much more and mutable, not list here. But one of them:

- https://github.com/urfave/cli/tree/v2



## License

MIT.




