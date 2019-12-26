/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hedzr/consul-tags/util"
	"github.com/hedzr/logex"
	"github.com/sirupsen/logrus"
	_ "log"
	"strconv"
	"strings"
	"testing"
)

func TestQueryConsulService(t *testing.T) {
	defer logex.CaptureLog(t).Release()

	registrar := getTestRegistrar()
	// client, err := api.NewClient(conf)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Querying DB services...")
	// db, err := QueryService(SERVICE_DB, registrar.FirstClient.Catalog())
	// if err != nil {
	// 	fmt.Errorf("Error: %v\n", err)
	// } else {
	// 	//registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
	// 	fmt.Printf("    Using '%s', %s:%d\n", db[0].ServiceID, db[0].Address, db[0].ServicePort)
	// }

	fmt.Printf("Querying 'consul' services...")
	consulService, err := QueryService("consul", registrar.FirstClient.Catalog())
	if err != nil {
		logrus.Fatalf("Error: %v", err)
		panic(err)
	} else {
		// registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
		// fmt.Printf("    Using '%s', %s:%d\n", userServices[0].ServiceID, userServices[0].Address, userServices[0].ServicePort)
		for i, cs := range consulService {
			fmt.Printf("    #%d. '%s', %s:%d, %v, %v\n", i, cs.ServiceID, cs.ServiceAddress, cs.ServicePort, cs.ServiceTags, cs.NodeMeta)
		}
	}
}

func TestQueryConsulapiService(t *testing.T) {
	registrar := getTestRegistrar()
	// client, err := api.NewClient(conf)
	// if err != nil {
	// 	panic(err)
	// }

	cc := GetConsulApiEntryPoint(registrar)
	if cc == nil {
		err := fmt.Errorf("Error: GetConsulApiEntryPoint() retun nil; 'consulapi' service NOT FOUND.")
		logrus.Fatal(err)
		panic(err)
	}
	fmt.Printf("%v", cc)
}

func TestConsulConnection(t *testing.T) {
	registrar := getTestRegistrar()

	if err := someTests(registrar); err != nil {
		t.Logf("TestConsulConnection() return ERROR: %v", err)
		t.Fail()
	}
}

func getTestRegistrar() *Registrar {
	return getRegistrarImpl(DEFAULT_CONSUL_LOCALHOST+":"+strconv.Itoa(DEFAULT_CONSUL_PORT), DEFAULT_CONSUL_SCHEME)
}

func someTests(registrar *Registrar) error {
	kv := registrar.FirstClient.KV()

	// PUT a new KV pair
	p := &api.KVPair{Key: KEY_WAS_SETUP, Value: []byte("---\nops\ncommon")}
	_, err := kv.Put(p, nil)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
		return err
	}

	// Lookup the pair
	pair, _, err := kv.Get(KEY_WAS_SETUP, nil)
	if err != nil || !strings.Contains(util.CToGoString(pair.Value), VALUE_WAS_SETUP) {
		logrus.Fatalf("Error: %v", err)
		return err
	}

	fmt.Printf("KV: %v\n", pair)

	fmt.Println("Querying datacenters...")
	catalog := registrar.FirstClient.Catalog()

	WaitForResult(func() (bool, error) {
		datacenters, err := catalog.Datacenters()
		if err != nil {
			return false, err
		}

		if len(datacenters) == 0 {
			return false, fmt.Errorf("Bad: %v\n", datacenters)
		} else {
			fmt.Printf("Datacenters: %v\n", datacenters)
		}

		return true, nil
	}, func(err error) {
		logrus.Fatalf("Error: %v", err)
		panic(err)
	})

	fmt.Printf("Querying nodes...")
	WaitForResult(func() (bool, error) {
		// meta := map[string]string{"somekey": "somevalue"}
		// catalogrus.Nodes(&QueryOptions{NodeMeta: meta})
		nodes, meta, err := catalog.Nodes(nil)
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
			fmt.Printf("    Nodes[i]: %v\n", node)
		}

		return true, nil
	}, func(err error) {
		logrus.Fatalf("Error: %v", err)
		panic(err)
	})

	fmt.Printf("Querying services...")
	WaitForResult(func() (bool, error) {
		services, meta, err := catalog.Services(nil)
		if err != nil {
			return false, err
		}

		if meta.LastIndex == 0 {
			return false, fmt.Errorf("Bad: %v", meta)
		}

		if len(services) == 0 {
			return false, fmt.Errorf("Bad: %v", services)
		}

		for i, service := range services {
			fmt.Printf("    Services[%s]: %v\n", i, service)
		}

		return true, nil
	}, func(err error) {
		logrus.Fatalf("Error: %v", err)
		panic(err)
	})

	return err
}

func testConsulapi(registrar *Registrar) {
	logrus.Debugf("Querying %s services...", SERVICE_CONSUL_API)
	theServices, err := QueryService(SERVICE_CONSUL_API, registrar.FirstClient.Catalog())
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	} else {
		// registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
		// fmt.Printf("    Using '%s', %s:%d\n", userServices[0].ServiceID, userServices[0].Address, userServices[0].ServicePort)
		for i, s := range theServices {
			fmt.Printf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
				i, s.ServiceID, s.ServiceAddress, s.ServicePort,
				s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			// NodeToAgent(registrar, s.Node).Node.Address
		}
	}
}

func testRedis(registrar *Registrar) {
	logrus.Debugf("Querying %s services...", SERVICE_CACHE)
	theServices, err := QueryService(SERVICE_CACHE, registrar.FirstClient.Catalog())
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	} else {
		// registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
		// fmt.Printf("    Using '%s', %s:%d\n", userServices[0].ServiceID, userServices[0].Address, userServices[0].ServicePort)
		for i, s := range theServices {
			fmt.Printf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
				i, s.ServiceID, s.ServiceAddress, s.ServicePort,
				s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			// NodeToAgent(registrar, s.Node).Node.Address
		}
	}
}

func testSqs(registrar *Registrar) {
	logrus.Debugf("Querying %s services...", SERVICE_MQ)
	theServices, err := QueryService(SERVICE_MQ, registrar.FirstClient.Catalog())
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	} else {
		// registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
		// fmt.Printf("    Using '%s', %s:%d\n", userServices[0].ServiceID, userServices[0].Address, userServices[0].ServicePort)
		for i, s := range theServices {
			// 服务 s 所在的 Node
			cn := NodeToAgent(registrar, s.Node)
			// 节点 cn 的服务表中名为 "consulapi" 的服务
			as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
			// 从 consulapi 指示Agent（也即服务 s 所对应的 Agent），建立一个临时的 Client
			client := getClientImpl(as.Address, as.Port, DEFAULT_CONSUL_SCHEME)
			agentService := cn.Services[s.ServiceID]
			tags := append(s.ServiceTags, "DEMO-DEMO")
			client.Agent().ServiceRegister(&api.AgentServiceRegistration{
				ID:                s.ServiceID,
				Name:              s.ServiceName,
				Tags:              tags,
				Port:              s.ServicePort,
				Address:           agentService.Address,
				EnableTagOverride: s.ServiceEnableTagOverride,
			})

			// 重新载入s的等价物，才能得到新的tags集合，s.ServiceTags并不会自动更新为新集合
			sNew, _ := QueryServiceByID(s.ServiceID, client)

			logrus.Infof("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s:%d\n",
				i, s.ServiceID, s.ServiceAddress, s.ServicePort,
				sNew.Tags, s.NodeMeta, s.Node, s.Address, as.Port)
		}
	}
}

func TestConsulRedisAndSqs(t *testing.T) {
	registrar := getTestRegistrar()
	// client, err := api.NewClient(conf)
	// if err != nil {
	// 	panic(err)
	// }

	cc := GetConsulApiEntryPoint(registrar)
	logrus.Debugf("GetConsulApiEntryPoint (via %s:%d): %v\n", DEFAULT_CONSUL_HOST, DEFAULT_CONSUL_PORT, cc)
	// for i, s := range cc.ServiceAddress {
	// 	fmt.Printf("    #%d. '%s', %s:%d, %v, %v, Node: %s,%s\n", i, s.ServiceID, s.ServiceAddress, s.ServicePort, s.ServiceTags, s.NodeMeta, s.Node, s.Address)
	// 	NodeToAgent(registrar, s.Node).Node.Address
	// }

	// fmt.Printf("Querying DB services...")
	// db, err := QueryService(SERVICE_DB, registrar.FirstClient.Catalog())
	// if err != nil {
	// 	fmt.Errorf("Error: %v\n", err)
	// } else {
	// 	//registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
	// 	fmt.Printf("    Using '%s', %s:%d\n", db[0].ServiceID, db[0].Address, db[0].ServicePort)
	// }

	testConsulapi(registrar)
	// testRedis(registrar)
	// testSqs(registrar, c)
}
