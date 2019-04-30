/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"github.com/hashicorp/consul/api"
	"time"
)

const (
	baseWait = 1 * time.Millisecond
	maxWait  = 100 * time.Millisecond

	SERVICE_CONSUL_API = "consulapi"
	SERVICE_DB         = "test-rds"
	SERVICE_MQ         = "test-mq"
	SERVICE_CACHE      = "test-redis"

	KEY_WAS_SETUP   = "ops/config/common"
	VALUE_WAS_SETUP = "---"
)

type Base struct {
	FirstClient *api.Client
}

type Registrar struct {
	Base
	Clients       []*api.CatalogService
	CurrentClient *api.CatalogService
}

type Discoverable interface {
	Connect(hostOrIp string, port int) error
	Register(serviceName string, ip string, port int) error
	Deregister(serviceName string, ip string, port int) error
	Disco(serviceName string) ([]Service, error)
}

type Service interface {
	GetName() string
	GetPort() int
	GetId() string
	GetTags() []string
}
