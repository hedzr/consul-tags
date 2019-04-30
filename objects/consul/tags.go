/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

func TagsList() error {
	registrar := getRegistrar()
	listServiceTags(registrar)
	return nil
}

func TagsToggle() error {
	registrar := getRegistrar()
	toggleServiceTags(registrar)
	return nil
}

func Tags() error {
	registrar := getRegistrar()

	if true {
		cc := GetConsulApiEntryPoint(registrar)
		logrus.Debugf("GetConsulApiEntryPoint (via %s): %v\n", viper.GetString("app.ms.addr"), cc)

		// for i, n := range viper.GetFlagNames() {
		// 	logrus.Debugf("    - flag name %d: %s, value: %v\n", i, n, viper.Get(n))
		// }
	}

	return modifyServiceTags(registrar)
}

func listServiceTags(registrar *Registrar) (err error) {
	name := viper.GetString("app.ms.name")
	if name != "" {
		listServiceTagsByName(registrar, name)
		return
	}
	id := viper.GetString("app.ms.id")
	if id != "" {
		listServiceTagsByID(registrar, id)
		return
	}

	return errors.New("--name ServiceName or --id ServiceID should be specified.")
}

func listServiceTagsByName(registrar *Registrar, serviceName string) {
	logrus.Debugf("List the tags of service '%s' at '%s'...", serviceName, viper.GetString("app.ms.addr"))

	WaitForResult(func() (bool, error) {
		catalogServices, err := QueryService(serviceName, registrar.FirstClient.Catalog())
		if err != nil {
			return false, err
		}

		if _, ok := catalogServices[0].TaggedAddresses["wan"]; !ok {
			return false, fmt.Errorf("Bad: %v\n", catalogServices[0])
		}

		for _, catalogService := range catalogServices {
			fmt.Printf("%s:\n", catalogService.ServiceID)
			fmt.Printf("\tname: %s\n", catalogService.ServiceName)
			fmt.Printf("\tnode: %s\n", catalogService.Node)
			fmt.Printf("\taddr: %s, tagged: %v\n", catalogService.Address, catalogService.TaggedAddresses)
			fmt.Printf("\tendp: %s:%d\n", catalogService.ServiceAddress, catalogService.ServicePort)
			fmt.Printf("\ttags: %v\n", strings.Join(catalogService.ServiceTags, ","))
			fmt.Printf("\tmeta: %v\n", catalogService.NodeMeta)
			fmt.Printf("\tenableTagOveerride: %v\n", catalogService.ServiceEnableTagOverride)
		}

		return true, nil
	}, func(err error) {
		logrus.Fatalf("err: %v", err)
	})
}

func listServiceTagsByID(registrar *Registrar, id string) {
	logrus.Debugf("List the tags of service by id '%s'...", id)
	as0, err := QueryServiceByID(id, registrar.FirstClient)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	} else {
		s, err1 := AgentServiceToCatalogService(as0, registrar.FirstClient)
		if err1 != nil {
			logrus.Fatalf("Error: %v", err)
			return
		}

		// logrus.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
		// 		1, s.ServiceID, s.ServiceAddress, s.ServicePort,
		// 		s.ServiceTags, s.NodeMeta, s.Node, s.Address)

		fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(s.ServiceTags, ","))
	}
}

func modifyServiceTags(registrar *Registrar) error {
	name := viper.GetString("app.ms.name")
	if name != "" {
		return modifyServiceTagsByName(registrar, name)
	}
	id := viper.GetString("app.ms.id")
	if id != "" {
		return modifyServiceTagsByID(registrar, id)
	}
	return errors.New("--name ServiceName or --id ServiceID should be specified.")
}

func modifyServiceTagsByName(registrar *Registrar, serviceName string) (err error) {
	logrus.Debugf("Modifying the tags of service '%s'...", serviceName)

	var (
		catalogServices []*api.CatalogService
	)

	catalogServices, err = QueryService(serviceName, registrar.FirstClient.Catalog())
	if err != nil {
		logrus.Fatalf("Error: %v", err)
		return
	}

	bothMode := viper.GetBool("app.ms.tags.both-mode")
	metaMode := viper.GetBool("app.ms.tags.meta-mode")

	// registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
	// fmt.Printf("    Using '%catalogService', %catalogService:%d\n", userServices[0].ServiceID, userServices[0].Address, userServices[0].ServicePort)
	for _, catalogService := range catalogServices {
		// 服务 catalogService 所在的 Node
		cn := NodeToAgent(registrar, catalogService.Node)
		// 节点 cn 的服务表中名为 "consulapi" 的服务
		as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
		// 从 consulapi 指示Agent（也即服务 catalogService 所对应的 Agent），建立一个临时的 Client
		client := getClient(as.Address, as.Port)
		agentService := cn.Services[catalogService.ServiceID]

		if bothMode || metaMode == false {

			logrus.Debugf("    %s: tags: %v", catalogService.ServiceID, catalogService.ServiceTags)

			tags := ModifyTags(catalogService.ServiceTags,
				viper.GetStringSlice("app.ms.tags.add-list"),
				viper.GetStringSlice("app.ms.tags.rm-list"),
				viper.GetString("app.ms.tags.delim"),
				viper.GetBool("app.ms.tags.clear"),
				viper.GetBool("app.ms.tags.plain-mode"),
				viper.GetBool("app.ms.tags.string-mode"))

			if err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
				ID:                catalogService.ServiceID,
				Name:              catalogService.ServiceName,
				Tags:              tags,
				Port:              catalogService.ServicePort,
				Address:           agentService.Address,
				EnableTagOverride: catalogService.ServiceEnableTagOverride,
			}); err != nil {
				logrus.Errorf("Error: %v", err)
				return
			}

			// 重新载入s的等价物，才能得到新的tags集合，catalogService.ServiceTags并不会自动更新为新集合
			sNew, _ := QueryServiceByID(catalogService.ServiceID, client)

			// fmt.Printf("    #%d. id='%catalogService'[%catalogService:%d], tags=%v, meta=%v, Node: %catalogService,%catalogService:%d\n",
			// 	i, catalogService.ServiceID, catalogService.ServiceAddress, catalogService.ServicePort,
			// 	sNew.Tags, catalogService.NodeMeta, catalogService.Node, catalogService.Address, as.Port)
			// logrus.Debugf("    #%d. id='%catalogService'[%catalogService:%d], tags=%v, meta=%v, Node: %catalogService,%catalogService\n",
			// 	i, catalogService.ServiceID, catalogService.ServiceAddress, catalogService.ServicePort,
			// 	catalogService.ServiceTags, catalogService.NodeMeta, catalogService.Node, catalogService.Address)
			fmt.Printf("%s: %s\n", catalogService.ServiceID, strings.Join(sNew.Tags, ","))
			fmt.Printf("\tmeta: %v\n", catalogService.NodeMeta)
		}

		if bothMode || metaMode {

			logrus.Debugf("    %s: meta: %v", catalogService.ServiceID, catalogService.NodeMeta)

			ModifyNodeMeta(catalogService.NodeMeta,
				viper.GetStringSlice("app.ms.tags.add-list"),
				viper.GetStringSlice("app.ms.tags.rm-list"),
				viper.GetString("app.ms.tags.delim"),
				viper.GetBool("app.ms.tags.clear"),
				false, // viper.GetBool("app.ms.tags.plain-mode"),
				viper.GetBool("app.ms.tags.string-mode"))

			logrus.Debugf("    %s: meta: %v, modified.", catalogService.ServiceID, catalogService.NodeMeta)

			// catalogService.NodeMeta["id"] = catalogService.ServiceID
			// catalogService.NodeMeta["addr"] = catalogService.Address
			// catalogService.NodeMeta["s-addr"] = catalogService.ServiceAddress
			// catalogService.NodeMeta["s-port"] = strconv.Itoa(catalogService.ServicePort)

			var writeMeta *api.WriteMeta
			writeMeta, err = client.Catalog().Register(&api.CatalogRegistration{
				ID:              catalogService.ID,
				Node:            catalogService.Node,
				Address:         catalogService.Address,
				TaggedAddresses: catalogService.TaggedAddresses,
				NodeMeta:        catalogService.NodeMeta,
				Service:         agentService,
				// Datacenter      : registrar.FirstClient.Catalog().Datacenters()[0],
			}, nil)
			if err != nil {
				logrus.Errorf("Error: %v", err)
				return
			}

			logrus.Debugf("\twriteMeta: %v", writeMeta)
		}
	}

	if bothMode || metaMode {
		fmt.Printf("**** Results of service '%s':\n", serviceName)
		WaitForResult(func() (bool, error) {
			catalogServicesNew, err := QueryService(serviceName, registrar.FirstClient.Catalog())
			if err != nil {
				return false, err
			}
			for _, catalogService := range catalogServicesNew {
				fmt.Printf("    %s: meta: %v.\n", catalogService.ServiceID, catalogService.NodeMeta)
			}
			return true, err
		}, func(err error) {
			logrus.Errorf("err: %v", err)
		})
	}
	return
}

func modifyServiceTagsByID(registrar *Registrar, id string) (err error) {
	logrus.Debugf("Modifying the tags of service by id '%s'...", id)

	var (
		as0, sNew *api.AgentService
		s         *api.CatalogService
	)

	as0, err = QueryServiceByID(id, registrar.FirstClient)
	if err != nil {
		logrus.Errorf("Error: %v", err)
		return
	}

	s, err = AgentServiceToCatalogService(as0, registrar.FirstClient)
	if err != nil {
		logrus.Errorf("Error: %v", err)
		return
	}

	// 服务 s 所在的 Node
	cn := NodeToAgent(registrar, s.Node)
	// 节点 cn 的服务表中名为 "consulapi" 的服务
	as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
	// 从 consulapi 指示Agent（也即服务 s 所对应的 Agent），建立一个临时的 Client
	client := getClient(as.Address, as.Port)
	agentService := cn.Services[id]

	tags := ModifyTags(s.ServiceTags,
		viper.GetStringSlice("app.ms.tags.add-list"),
		viper.GetStringSlice("app.ms.tags.rm-list"),
		viper.GetString("app.ms.tags.delim"),
		viper.GetBool("app.ms.tags.clear"),
		viper.GetBool("app.ms.tags.plain-mode"),
		viper.GetBool("app.ms.tags.string-mode"))

	// for _, t = range tags {
	// 	logrus.Debugf("    *** Tags: %v", tags)
	// }

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:                as0.ID,
		Name:              as0.Service,
		Tags:              tags,
		Port:              as0.Port,
		Address:           agentService.Address,
		EnableTagOverride: as0.EnableTagOverride,
	})
	if err != nil {
		logrus.Errorf("Error: %v", err)
		return
	}

	// 重新载入s的等价物，才能得到新的tags集合，s.ServiceTags并不会自动更新为新集合
	sNew, err = QueryServiceByID(as0.ID, client)
	if err != nil {
		logrus.Errorf("Error: %v", err)
		return
	}

	// fmt.Printf("    id='%s'[%s:%d], tags=%v, Node: %s,%s:%d\n",
	// 	as0.ID, as0.Address, as0.Port,
	// 	sNew.Tags, s.Node, as0.Address, as.Port)
	// logrus.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
	// 	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
	// 	s.ServiceTags, s.NodeMeta, s.Node, s.Address)
	fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(sNew.Tags, ","))
	return
}

func toggleServiceTags(registrar *Registrar) {
	name := viper.GetString("app.ms.name")
	if name != "" {
		toggleServiceTagsByName(registrar, name)
		return
	}
	id := viper.GetString("app.ms.id")
	if id != "" {
		logrus.Fatalf("toggle tags can be applied with --name but --id")
		return
	}
}

func toggleServiceTagsByName(registrar *Registrar, name string) {
	logrus.Debugf("Toggle the tags of service '%s'...", name)
	theServices, err := QueryService(name, registrar.FirstClient.Catalog())
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	} else {
		newMaster := strings.Split(viper.GetString("app.ms.tags.toggle.address"), ":")
		newMasterPort := 0
		if len(newMaster) > 1 {
			newMasterPort, err = strconv.Atoi(newMaster[1])
			if err != nil {
				logrus.Fatalf("Error: %v", err)
				return
			}
		}
		masterTag := viper.GetStringSlice("app.ms.tags.toggle.set-list")
		slaveTag := viper.GetStringSlice("app.ms.tags.toggle.reset-list")
		if len(newMaster) == 0 {
			logrus.Fatalf("--address to specify the master ip:port, it's NOT optional.")
			return
		}
		// if len(masterTag) == 1 {
		if len(masterTag) == 0 {
			logrus.Fatalf("--set to specify the master tag, it's NOT optional.")
			return
		}
		if len(slaveTag) == 0 {
			logrus.Fatalf("--reset to specify the slave tag, it's NOT optional.")
			return
		}

		for i, s := range theServices {
			fmt.Printf("    #%d. id='%s'[%s:%d #%s], tags=%v, meta=%v, Node: %s,%s\n",
				i, s.ServiceID, s.ServiceAddress, s.ServicePort, s.Address,
				s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			// for _, t := range s.ServiceTags {
			matched := strings.EqualFold(s.ServiceAddress, newMaster[0])
			if matched && len(newMaster) > 1 {
				matched = s.ServicePort == newMasterPort
			}
			tags := s.ServiceTags
			if matched {
				tags = ModifyTags(tags, masterTag, slaveTag,
					viper.GetString("app.ms.tags.delim"),
					viper.GetBool("app.ms.tags.clear"),
					viper.GetBool("app.ms.tags.plain-mode"),
					viper.GetBool("app.ms.tags.string-mode"))
			} else {
				tags = ModifyTags(tags, slaveTag, masterTag,
					viper.GetString("app.ms.tags.delim"),
					viper.GetBool("app.ms.tags.clear"),
					viper.GetBool("app.ms.tags.plain-mode"),
					viper.GetBool("app.ms.tags.string-mode"))
			}
			// }

			// for _, t = range tags {
			// logrus.Debugf("    *** Tags: %v\n", tags)
			// }

			cn := NodeToAgent(registrar, s.Node)
			as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
			logrus.Debugf("    %s=%v\n", SERVICE_CONSUL_API, as)
			client := getClient(as.Address, as.Port)
			agentService := cn.Services[s.ServiceID]

			err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
				ID:                s.ServiceID,
				Name:              s.ServiceName,
				Tags:              tags,
				Port:              s.ServicePort,
				Address:           agentService.Address,
				EnableTagOverride: s.ServiceEnableTagOverride,
			})
			if err != nil {
				logrus.Fatalf("Error: %v", err)
				return
			}

			// csa, _, err := client.Catalog().Service(s.ServiceName, "", nil)
			// var cs *api.CatalogService = nil
			// for _, cs0 := range csa {
			// 	if strings.EqualFold(cs0.ServiceID, s.ServiceID) {
			// 		cs = cs0
			// 		break
			// 	}
			// }
			// if err != nil {
			// 	logrus.Fatal(fmt.Errorf("Error: %v, %v", err, cs))
			// }
			// client.Catalog().Register(&api.CatalogRegistration{
			// 	ID: s.ServiceID,
			// 	Node: s.Node,
			// 	Address: s.ServiceAddress,
			// 	NodeMeta: s.NodeMeta,
			// 	Service: agentService,
			// }, &api.WriteOptions{})

			// 重新载入s的等价物，才能得到新的tags集合，s.ServiceTags并不会自动更新为新集合
			sNew, err := QueryServiceByID(s.ServiceID, client)
			if err != nil {
				logrus.Fatalf("Error: %v", err)
				return
			}

			// fmt.Printf("    TAGS=%v\n\n", sNew.Tags)
			// logrus.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
			// 	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
			// 	s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(sNew.Tags, ","))
		}

		//
	}
}

func getRegistrar() *Registrar {
	return getRegistrarImpl(viper.GetString("app.ms.addr"), viper.GetString("app.ms.scheme"))
}

func getRegistrarImpl(addr, scheme string) *Registrar {
	return &Registrar{
		Base: Base{
			FirstClient: MakeClientWithConfig(func(clientConfig *api.Config) {
				clientConfig.Address = addr
				clientConfig.Scheme = scheme
			}),
		},
		Clients:       nil,
		CurrentClient: nil,
	}
}

func getClient(host string, port int) *api.Client {
	return getClientImpl(host, port, viper.GetString("app.ms.scheme"))
}

func getClientImpl(host string, port int, scheme string) *api.Client {
	return MakeClientWithConfig(func(clientConfig *api.Config) {
		clientConfig.Address = host + ":" + strconv.Itoa(port)
		clientConfig.Scheme = scheme
	})
}
