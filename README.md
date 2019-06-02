# Consul Service Tags Modifier

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/consul-tags.svg?label=release)](https://github.com/hedzr/consul-tags/releases)
[![Build Status](https://travis-ci.org/hedzr/consul-tags.svg?branch=master)](https://travis-ci.org/hedzr/consul-tags)
[![license](https://img.shields.io/github/license/hedzr/consul-tags.svg)]()
[![Docker Automated buil](https://img.shields.io/docker/automated/hedzr/consul-tags.svg)]()
[![Docker Pulls](https://img.shields.io/docker/pulls/hedzr/consul-tags.svg)]()
<!-- [![GitHub tag](https://img.shields.io/github/tag/hedzr/consul-tags.svg)]() -->
<!-- [![ImageLayers Size](https://img.shields.io/imagelayers/image-size/hedzr/consul-tags/latest.svg)]() -->

<!-- [![GitHub version](https://badge.fury.io/gh/hedzr%2Fconsul-tags.svg)](https://badge.fury.io/gh/hedzr%2Fconsul-tags)
-->

`consul-tags` can add, remove tag(s) of a consul service (one or all its instances).

Here is a first release for key functionality.



## News

### v0.5.1

- new release has been testing and released soon.
- it has been rewrote and optimized.

### v0.5.0 is a pre-release

- rewrote by [`cmdr`](https://github.com/hedzr/cmdr)
- pre-released for v0.5.1 (final)



## Install

### Binary

Download binary from [Release](../../releases/latest) page.

### Docker Hub (paused since v0.5.x)

```bash
docker pull hedzr/consul-tags
docker run -it --rm hedzr/consul-tags --addr 192.168.0.71:8500 ms --name test-redis tags ls
```

Replace `192.168.0.71` with your consul center ip or name.

DON'T use `127.0.0.1` with dockerize release.

> latest: master branch and based on golang:alpine


### Go Build

clone the repo and build:

```bash
cd $GOPATH/github.com/hedzr/consul-tags
go mod download
go build -o consul-tags ./cli/main.go 
```

mixin:

```bash
go get -u github.com/hedzr/consul-tags
```


## Usage


```bash

# run consul demo instance for testing
./build.sh consul run &

# use the local consul demo instance as default addrress, see also `--addr` in `consul-tags ms --help`
export CT_APP_MS_ADDR=localhost:8500

# list services
consul-tags ms ls

# list tags
consul-tags ms --name test-redis tags ls
consul-tags ms tags ls --name test-redis

# modify tags
consul-tags ms tags modify --name test-redis tags --add a,c,e
consul-tags ms tags mod --name test-redis --add a --add c --add e
consul-tags ms tags mod --name test-redis --rm a,c,e
consul-tags ms tags mod --name test-redis --rm a --rm c --rm e
# by id
consul-tags ms tags ls --id test-redis-1

# toggle master/slave
consul-tags ms tags toggle --name test-redis --addr 10.7.13.1:6379 --set role=master --reset role=slave
consul-tags ms tags tog --name test-mq --addr 10.7.16.3:5672 --set role=leader,type=ram --reset role=peer,type=disk

# get commands and sub-commands help
consul-tags ms -h
# get help: -h or --help
consul-tags -h
```

Default consul address is `consul.ops.local:8500`, can be overridden with environment variable `CT_APP_MS_ADDR` (host:port), too. Such as:

```bash
export CT_APP_MS_ADDR=localhost:8500
consul-tags ms tags ls --name=consul
# or
CT_APP_MS_ADDR=localhost:8500 consul-tags --help
# or
consul-tags ms tags ls --name=consul --addr 127.0.0.1:8500
```

### COMMANDS:

```
     kv              K/V pair Operations, ...
     ms, service, m  Microservice Operations, ...
     help, h         Shows a list of commands or help for one command
```

### GLOBAL OPTIONS:

```
   --help, -h                          show help (default: false)
   --version, -v                       print the version (default: false)
```

## Shell completion

TODO

```bash

```

or put it to `/etc/bash_completion.d/`:

```bash

```



## Thanks

my plan is building a suite of devops. consul operations is part of it.


and some repositories are good:

- [colebrumley/consul-kv-backup](https://github.com/colebrumley/consul-kv-backup)
- <https://github.com/shreyu86/consul-backup>
- <https://github.com/kailunshi/consul-backup>

### Third-party repos




## License

MIT.




