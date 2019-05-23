/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type testFn func() (bool, error)
type errorFn func(error)

func WaitForResult(try testFn, fail errorFn) {
	var err error
	wait := baseWait
	for retries := 100; retries > 0; retries-- {
		var success bool
		success, err = try()
		if success {
			time.Sleep(25 * time.Millisecond)
			return
		}

		time.Sleep(wait)
		wait *= 2
		if wait > maxWait {
			wait = maxWait
		}
	}
	fail(err)
}

type configCallback func(c *api.Config)

func MakeClient() *api.Client {
	return MakeClientWithConfig(nil)
}

func MakeACLClient() *api.Client {
	return MakeClientWithConfig(
		// t,
		func(clientConfig *api.Config) {
			clientConfig.Token = "root"
		},
		// , func(serverConfig *testutil.TestServerConfig) {
		// 	serverConfig.ACLMasterToken = "root"
		// 	serverConfig.ACLDatacenter = "dc1"
		// 	serverConfig.ACLDefaultPolicy = "deny"
		// }
	)
}

func MakeClientWithConfig(cb1 configCallback) *api.Client {

	// Make client config
	conf := api.DefaultConfig()
	if cb1 != nil {
		cb1(conf)
	}

	// // Create server
	// server := testutil.NewTestServerConfig(t, cb2)
	// conf.Address = server.HTTPAddr

	// Create client
	client, err := api.NewClient(conf)
	if err != nil {
		logrus.Fatalf("err: %v", err)
	}

	return client // , server
}

func QueryService(name string, catalog *api.Catalog) ([]*api.CatalogService, error) {
	// metaQ := map[string]string{"Name": name}
	services, meta, err := catalog.Service(name, "", nil) // &api.QueryOptions{NodeMeta: metaQ})
	if err != nil {
		return nil, err
	}

	if meta.LastIndex == 0 {
		return nil, fmt.Errorf("Bad: %v", meta)
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("Bad: %v", services)
	}
	return services, nil
}

func QueryServiceByID(serviceID string, client *api.Client) (as *api.AgentService, err error) {
	var res *api.AgentService = nil
	WaitForResult(func() (bool, error) {
		cn, err := client.Agent().Services()
		if err != nil {
			return false, err
		}

		for id, s := range cn {
			if strings.EqualFold(id, serviceID) {
				res = s
				return true, nil
			}
		}

		return false, fmt.Errorf("Bad: cannot found service '#%s'", serviceID)
	}, func(err error) {
		logrus.Fatalf("err: %v", err)
	})
	return res, nil
}

func AgentServiceToCatalogService(as *api.AgentService, client *api.Client) (res *api.CatalogService, err error) {
	var cn []*api.CatalogService = nil
	WaitForResult(func() (bool, error) {
		catalog := client.Catalog()
		cn, _, err = catalog.Service(as.Service, "", nil)
		if err != nil {
			return false, err
		}
		for _, cs := range cn {
			if cs.ServiceID == as.ID {
				res = cs
				return true, nil
			}
		}
		return false, fmt.Errorf("Bad: cannot found service '#%s' inside catalog", as.ID)
	}, func(err error) {
		logrus.Fatalf("err: %v", err)
	})
	return
}

func CatalogNodeGetService(cn *api.CatalogNode, serviceName string) *api.AgentService {
	for _, val := range cn.Services {
		if strings.EqualFold(val.Service, serviceName) {
			return val
		}
	}
	return nil
}

func NodeToAgent(registrar *Registrar, node string) *api.CatalogNode {
	cn, qm, err := registrar.FirstClient.Catalog().Node(node, nil)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	} else {
		logrus.Debugf("    QueryMeta: %v", qm)
		// cn.Node.Address
		return cn
	}

	fmt.Println("Querying nodes...")
	WaitForResult(func() (bool, error) {
		// meta := map[string]string{"somekey": "somevalue"}
		// catalogrus.Nodes(&QueryOptions{NodeMeta: meta})
		nodes, meta, err := registrar.FirstClient.Catalog().Nodes(nil)
		if err != nil {
			return false, err
		}

		if meta.LastIndex == 0 {
			return false, fmt.Errorf("Bad: %v", meta)
		}

		if len(nodes) == 0 {
			return false, fmt.Errorf("Bad: %v", nodes)
		}

		if _, ok := nodes[0].TaggedAddresses["wan"]; !ok {
			return false, fmt.Errorf("Bad: %v\n", nodes[0])
		}

		for _, node := range nodes {
			logrus.Debugf("    Nodes[i]: %v", node)
		}

		return true, nil
	}, func(err error) {
		logrus.Fatalf("err: %v", err)
	})
	return nil
}

func GetConsulApiEntryPoint(registrar *Registrar) *api.CatalogService {
	var err error = nil
	registrar.Clients, err = QueryService(SERVICE_CONSUL_API, registrar.FirstClient.Catalog())
	if err != nil {
		logrus.Fatalf("err: %v", err)
		return nil
	} else {
		// registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
		logrus.Tracef("    Using '%s', %s:%d", registrar.Clients[0].ServiceID, registrar.Clients[0].Address, registrar.Clients[0].ServicePort)
		registrar.CurrentClient = registrar.Clients[0]
		return registrar.CurrentClient
	}

	// consulapi := findConsulApi(base)
	// if len(consulapi) > 0 {
	// 	registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
	// 	fmt.Printf("    Using '%s', %s:%d\n", registrarId, registrarAddr, registrarPort)
	// }
}

func findConsulApi(base *Base) []*api.CatalogService {
	services, err := QueryService(SERVICE_CONSUL_API, base.FirstClient.Catalog())
	if err != nil {
		logrus.Fatalf("err: %v", err)
		return nil
	} else {
		for i, service := range services {
			logrus.Tracef("    Service[%d, %s]: %v\n", i, service.ServiceID, service)
		}
		return services
	}
}
