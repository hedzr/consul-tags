# Consul Service Tags Modifier

![Go](https://github.com/hedzr/consul-tags/workflows/Go/badge.svg)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/consul-tags.svg?label=release)](https://github.com/hedzr/consul-tags/releases)
[![go.dev](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/consul-tags)
[![license](https://img.shields.io/github/license/hedzr/consul-tags.svg)](https://pkg.go.dev/github.com/hedzr/consul-tags)
[![Docker Automated build](https://img.shields.io/docker/automated/hedzr/consul-tags.svg)](https://hub.docker.com/r/hedzr/consul-tags)
[![Docker Pulls](https://img.shields.io/docker/pulls/hedzr/consul-tags.svg)](https://hub.docker.com/r/hedzr/consul-tags)
<!-- [![GitHub tag](https://img.shields.io/github/tag/hedzr/consul-tags.svg)]() -->
<!-- [![ImageLayers Size](https://img.shields.io/imagelayers/image-size/hedzr/consul-tags/latest.svg)]() -->

<!-- [![GitHub version](https://badge.fury.io/gh/hedzr%2Fconsul-tags.svg)](https://badge.fury.io/gh/hedzr%2Fconsul-tags)
-->

`consul-tags` can add, remove tag(s) of a consul service (one or all its instances).

Here is a first release for key functionality.

## News

- v0.8.0
  - Move to cmdr.v2 and rewrite app
  - trying to fix #8

- See [CHANGELOG](https://github.com/hedzr/consul-tags/blob/master/CHANGELOG)

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
# export CT_APP_MS_ADDR=localhost:8500
export ADDR=localhost:8500

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
consul-tags ms tags ls --id test-redis-6379
consul-tags ms tags ls --id test-redis-6380

# toggle master/slave
# if test-redis nodes are 10.7.13.1,...
consul-tags ms tags toggle --name test-redis --service-addr 10.7.13.1:6379 --set role=master --reset role=slave
consul-tags ms tags tog --name test-mq --addr 10.7.16.3:5672 --set role=leader,type=ram --reset role=peer,type=disk
# if test-redis nodes are localhost:6379,...
consul-tags ms tags toggle --name test-redis --service-addr localhost:6380 --set role=master --reset role=slave

# get commands and sub-commands help
consul-tags ms -h
# get help: -h or --help
consul-tags -h
```

Default consul address is `consul.ops.local:8500. But it also can be overridden with environment variable `ADDR` (host:port). Such as:

```bash
export ADDR=localhost:8500
consul-tags ms tags ls --name=consul
# or
ADDR=localhost:8500 consul-tags --help
# or
consul-tags ms tags ls --name=consul --addr 127.0.0.1:8500
```

To have a see about which envvars are valid, check out the command's help screen. For example, `ms tags ls --help`:

![image-20241025095348132](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/upgit/2024/10/20241025_1729825365.png)

In this case, envvar `ADDR` takes effect to the flag `--addr`, and both flags supported envvar also print the tip at its description section.

## Shell completion

> Since v0.8.0, these addons are still in developing in cmdr.v2

```bash
consul-tags gen sh --zsh
consul-tags gen sh --bash
```

## All Commands

Use cmdr builtin `~~tree` to have a bird's eye:

![Screenshot 2024-10-25 at 09.21.43](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/upgit/2024/10/20241025_1729819334.png)



## License

Apache 2.0

> Since v0.8.0, we moved license from MIT to Apache 2.0. It's hardly a little bit effects for you, basically.
