/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	cmdrv2 "github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"gopkg.in/hedzr/errors.v3"
)

const (
	// MS_PREFIX   = "app.ms"
	// TAGS_PREFIX = "app.ms.tags"

	MS_PREFIX_V2   = "ms"
	TAGS_PREFIX_V2 = "ms.tags"
)

func TagsList() (err error) {
	registrar := getRegistrarV2()
	err = listServiceTagsV2(registrar)
	return
}

func TagsToggle() error {
	registrar := getRegistrarV2()
	toggleServiceTagsV2(registrar)
	return nil
}

func Tags() error {
	registrar := getRegistrarV2()

	if true {
		cc := GetConsulApiEntryPoint(registrar)
		// cmdr.Logger.Debugf("GetConsulApiEntryPoint (via %s): %v\n", cmdr.GetStringP(TAGS_PREFIX, "addr"), cc)
		logz.Debug("GetConsulApiEntryPoint", "addr", registrar.Clients[0].Address, "cc", cc)

		// for i, n := range viper.GetFlagNames() {
		// 	cmdr.Logger.Debugf("    - flag name %d: %s, value: %v\n", i, n, viper.Get(n))
		// }
	}

	return modifyServiceTagsV2(registrar)
}

func listServiceTagsV2(registrar *Registrar) (err error) {
	cs2 := cmdrv2.Store().WithPrefix(MS_PREFIX_V2)
	name := cs2.MustString("name")
	if name != "" {
		listServiceTagsByNameV2(registrar, name)
		return
	}
	id := cs2.MustString("id")
	if id != "" {
		listServiceTagsByIDV2(registrar, id)
		return
	}

	return errors.New("--name ServiceName or --id ServiceID should be specified.")
}

func listServiceTagsByNameV2(registrar *Registrar, serviceName string) {
	// cmdr.Logger.Debugf("List the tags of service '%s' at '%s'...", serviceName, cmdr.GetStringP(TAGS_PREFIX, "addr"))
	logz.Debug("List service tags by name...", "name", serviceName, "registrar", registrar)

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
		logz.Fatal("QueryService failed", "err", err)
	})
}

func listServiceTagsByIDV2(registrar *Registrar, id string) {
	// cmdr.Logger.Debugf("List the tags of service by id '%s'...", id)
	logz.Debug("List service tags by id...", "id", id, "registrar", registrar)

	as0, err := QueryServiceByID(id, registrar.FirstClient)
	if err != nil {
		// cmdr.Logger.Fatalf("Error: %v", err)
		logz.Fatal("QueryServiceByID failed", "err", err)
	} else {
		s, err1 := AgentServiceToCatalogService(as0, registrar.FirstClient)
		if err1 != nil {
			// cmdr.Logger.Fatalf("Error: %v", err)
			logz.Fatal("AgentServiceToCatalogService failed", "err", err)
			return
		}

		// cmdr.Logger.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
		// 		1, s.ServiceID, s.ServiceAddress, s.ServicePort,
		// 		s.ServiceTags, s.NodeMeta, s.Node, s.Address)

		fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(s.ServiceTags, ","))
	}
}

func modifyServiceTagsV2(registrar *Registrar) error {
	cs2 := cmdrv2.Store().WithPrefix(MS_PREFIX_V2)
	name := cs2.MustString("name")
	if name != "" {
		return modifyServiceTagsByNameV2(registrar, name)
	}
	id := cs2.MustString("id")
	if id != "" {
		return modifyServiceTagsByIDV2(registrar, id)
	}
	return errors.New("--name ServiceName or --id ServiceID should be specified.")
}

func modifyServiceTagsByNameV2(registrar *Registrar, serviceName string) (err error) {
	// cmdr.Logger.Debugf("Modifying the tags of service '%s'...", serviceName)
	logz.Debug("Modifying the tags of service '" + serviceName + "'...")

	var (
		catalogServices []*api.CatalogService
		cs3             = cmdrv2.Store().WithPrefix(TAGS_PREFIX_V2)
		bothMode        = cs3.MustBool("modify.both")
		metaMode        = cs3.MustBool("modify.meta")
		plainMode       = cs3.MustBool("modify.plain")
		stringMode      = cs3.MustBool("modify.string")
		addList         = cs3.MustStringSlice("modify.add")
		rmList          = cs3.MustStringSlice("modify.remove")
		delim           = cs3.MustString("modify.delim")
		clearFlag       = cs3.MustBool("modify.clear")
	)

	catalogServices, err = QueryService(serviceName, registrar.FirstClient.Catalog())
	if err != nil {
		// cmdr.Logger.Fatalf("Error: %v", err)
		logz.Fatal("QueryService failed", "err", err)
		return
	}

	// registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
	// fmt.Printf("    Using '%catalogService', %catalogService:%d\n", userServices[0].ServiceID, userServices[0].Address, userServices[0].ServicePort)
	for _, catalogService := range catalogServices {
		// 服务 catalogService 所在的 Node
		cn := NodeToAgent(registrar, catalogService.Node)
		// 节点 cn 的服务表中名为 "consulapi" 的服务
		as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
		// 从 consulapi 指示Agent（也即服务 catalogService 所对应的 Agent），建立一个临时的 Client
		client := getClientV2(as.Address, as.Port)
		agentService := cn.Services[catalogService.ServiceID]

		if bothMode || metaMode == false {

			// cmdr.Logger.Debugf("    %s: tags: %v", catalogService.ServiceID, catalogService.ServiceTags)

			tags := ModifyTags(catalogService.ServiceTags, addList, rmList, delim, clearFlag, plainMode, stringMode)

			if err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
				ID:                catalogService.ServiceID,
				Name:              catalogService.ServiceName,
				Tags:              tags,
				Port:              catalogService.ServicePort,
				Address:           agentService.Address,
				EnableTagOverride: catalogService.ServiceEnableTagOverride,
			}); err != nil {
				// cmdr.Logger.Errorf("Error: %v", err)
				logz.Fatal("ServiceRegister failed", "err", err)
				return
			}

			// 重新载入s的等价物，才能得到新的tags集合，catalogService.ServiceTags并不会自动更新为新集合
			sNew, _ := QueryServiceByID(catalogService.ServiceID, client)

			// fmt.Printf("    #%d. id='%catalogService'[%catalogService:%d], tags=%v, meta=%v, Node: %catalogService,%catalogService:%d\n",
			// 	i, catalogService.ServiceID, catalogService.ServiceAddress, catalogService.ServicePort,
			// 	sNew.Tags, catalogService.NodeMeta, catalogService.Node, catalogService.Address, as.Port)
			// cmdr.Logger.Debugf("    #%d. id='%catalogService'[%catalogService:%d], tags=%v, meta=%v, Node: %catalogService,%catalogService\n",
			// 	i, catalogService.ServiceID, catalogService.ServiceAddress, catalogService.ServicePort,
			// 	catalogService.ServiceTags, catalogService.NodeMeta, catalogService.Node, catalogService.Address)
			fmt.Printf("%s: %s\n", catalogService.ServiceID, strings.Join(sNew.Tags, ","))
			fmt.Printf("\tmeta: %v\n", catalogService.NodeMeta)
		}

		if bothMode || metaMode {

			// cmdr.Logger.Debugf("    %s: meta: %v", catalogService.ServiceID, catalogService.NodeMeta)

			ModifyNodeMeta(catalogService.NodeMeta, addList, rmList, delim, clearFlag, false, stringMode)

			// cmdr.Logger.Debugf("    %s: meta: %v, modified.", catalogService.ServiceID, catalogService.NodeMeta)

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
				// cmdr.Logger.Errorf("Error: %v", err)
				logz.Fatal("Register failed", "err", err)
				return
			}

			// cmdr.Logger.Debugf("\twriteMeta: %v", writeMeta)
			logz.Debug("Meta:", "writeMeta", writeMeta)
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
			// cmdr.Logger.Errorf("err: %v", err)
			logz.Fatal("QueryService failed", "err", err)
		})
	}
	return
}

func modifyServiceTagsByIDV2(registrar *Registrar, id string) (err error) {
	// cmdr.Logger.Debugf("Modifying the tags of service by id '%s'...", id)
	logz.Debug("Modifying the tags of service by id '" + id + "'...")

	var (
		as0, sNew  *api.AgentService
		s          *api.CatalogService
		cs3        = cmdrv2.Store().WithPrefix(TAGS_PREFIX_V2)
		addList    = cs3.MustStringSlice("modify.add")
		rmList     = cs3.MustStringSlice("modify.remove")
		delim      = cs3.MustString("modify.delim")
		clearFlag  = cs3.MustBool("modify.clear")
		plainMode  = cs3.MustBool("modify.plain")
		stringMode = cs3.MustBool("modify.string")
		// bothMode   = cmdr.GetBoolP(TAGS_PREFIX, "modify.both")
		// metaMode   = cmdr.GetBoolP(TAGS_PREFIX, "modify.meta")
	)

	as0, err = QueryServiceByID(id, registrar.FirstClient)
	if err != nil {
		// cmdr.Logger.Errorf("Error: %v", err)
		logz.Fatal("QueryServiceByID failed", "err", err)
		return
	}

	s, err = AgentServiceToCatalogService(as0, registrar.FirstClient)
	if err != nil {
		// cmdr.Logger.Errorf("Error: %v", err)
		logz.Fatal("AgentServiceToCatalogService failed", "err", err)
		return
	}

	// 服务 s 所在的 Node
	cn := NodeToAgent(registrar, s.Node)
	// 节点 cn 的服务表中名为 "consulapi" 的服务
	as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
	// 从 consulapi 指示Agent（也即服务 s 所对应的 Agent），建立一个临时的 Client
	client := getClientV2(as.Address, as.Port)
	agentService := cn.Services[id]

	tags := ModifyTags(s.ServiceTags, addList, rmList, delim, clearFlag, plainMode, stringMode)

	// for _, t = range tags {
	// 	cmdr.Logger.Debugf("    *** Tags: %v", tags)
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
		logz.Error("Error:", "err", err)
		return
	}

	// 重新载入s的等价物，才能得到新的tags集合，s.ServiceTags并不会自动更新为新集合
	sNew, err = QueryServiceByID(as0.ID, client)
	if err != nil {
		// cmdr.Logger.Errorf("Error: %v", err)
		logz.Fatal("QueryServiceByID failed", "err", err)
		return
	}

	// fmt.Printf("    id='%s'[%s:%d], tags=%v, Node: %s,%s:%d\n",
	// 	as0.ID, as0.Address, as0.Port,
	// 	sNew.Tags, s.Node, as0.Address, as.Port)
	// cmdr.Logger.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
	// 	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
	// 	s.ServiceTags, s.NodeMeta, s.Node, s.Address)
	fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(sNew.Tags, ","))
	return
}

func toggleServiceTagsV2(registrar *Registrar) {
	cs2 := cmdrv2.Store().WithPrefix(MS_PREFIX_V2)
	name := cs2.MustString("name")
	if name != "" {
		toggleServiceTagsByNameV2(registrar, name)
		return
	}
	id := cs2.MustString("id")
	if id != "" {
		logz.Fatal("toggle tags can be applied with --name, rather --id")
		return
	}
}

func toggleServiceTagsByNameV2(registrar *Registrar, name string) {
	var (
		theServices []*api.CatalogService
		err         error
		cs3         = cmdrv2.Store().WithPrefix(TAGS_PREFIX_V2)
		masterTag   = cs3.MustStringSlice("toggle.set")
		slaveTag    = cs3.MustStringSlice("toggle.unset")
		addresses   = cs3.MustString("toggle.service-addr")
		delim       = cs3.MustString("toogle.delim")
		clearFlag   = cs3.MustBool("toggle.clear")
		plainMode   = cs3.MustBool("modify.plain")
		stringMode  = cs3.MustBool("modify.string")
		// bothMode   = cmdr.GetBoolP(TAGS_PREFIX, "modify.both")
		// metaMode   = cmdr.GetBoolP(TAGS_PREFIX, "modify.meta")
	)

	// cmdr.Logger.Debugf("Toggle the tags of service '%s'...", name)
	logz.Debug("Toggle the tags of service '" + name + "'...")

	theServices, err = QueryService(name, registrar.FirstClient.Catalog())
	if err != nil {
		// cmdr.Logger.Fatalf("Error: %v", err)
		logz.Fatal("QueryService failed", "err", err)
	} else {
		newMaster := strings.Split(addresses, ":")
		newMasterPort := 0
		if len(newMaster) > 1 {
			newMasterPort, err = strconv.Atoi(newMaster[1])
			if err != nil {
				// cmdr.Logger.Fatalf("Error: %v", err)
				logz.Fatal("Atoi failed", "err", err)
				return
			}
		}
		if len(newMaster) == 0 {
			logz.Fatal("--address to specify the master ip:port, it's NOT optional.")
			return
		}
		// if len(masterTag) == 1 {
		if len(masterTag) == 0 {
			logz.Fatal("--set to specify the master tag, it's NOT optional.")
			return
		}
		if len(slaveTag) == 0 {
			logz.Fatal("--reset to specify the slave tag, it's NOT optional.")
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
				tags = ModifyTags(tags, masterTag, slaveTag, delim, clearFlag, plainMode, stringMode)
			} else {
				tags = ModifyTags(tags, slaveTag, masterTag, delim, clearFlag, plainMode, stringMode)
			}
			// }

			// for _, t = range tags {
			// cmdr.Logger.Debugf("    *** Tags: %v\n", tags)
			// }

			cn := NodeToAgent(registrar, s.Node)
			as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
			// cmdr.Logger.Debugf("    %s=%v\n", SERVICE_CONSUL_API, as)
			logz.Debug("CatalogNodeGetService:", "API", SERVICE_CONSUL_API, "as", as)
			client := getClientV2(as.Address, as.Port)
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
				// cmdr.Logger.Fatalf("Error: %v", err)
				logz.Fatal("ServiceRegister failed", "err", err)
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
			// 	cmdr.Logger.Fatal(fmt.Errorf("Error: %v, %v", err, cs))
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
				// cmdr.Logger.Fatalf("Error: %v", err)
				logz.Fatal("QueryServiceByID failed", "err", err)
				return
			}

			// fmt.Printf("    TAGS=%v\n\n", sNew.Tags)
			// cmdr.Logger.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
			// 	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
			// 	s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(sNew.Tags, ","))
		}

		//
	}
}

func getRegistrarV2() *Registrar {
	cs2 := cmdrv2.Store().WithPrefix(MS_PREFIX_V2)
	addr := cs2.MustString("addr")
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf("%v:%v", addr, cs2.MustInt("port"))
	}
	scheme := cs2.MustString("scheme")
	return getRegistrarImplV2(addr, scheme)
}

func getRegistrarImplV2(addr, scheme string) *Registrar {
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

func getClientV2(host string, port int) *api.Client {
	cs2 := cmdrv2.Store().WithPrefix(MS_PREFIX_V2)
	scheme := cs2.MustString("scheme")
	return getClientImplV2(host, port, scheme)
}

func getClientImplV2(host string, port int, scheme string) *api.Client {
	return MakeClientWithConfig(func(clientConfig *api.Config) {
		clientConfig.Address = host + ":" + strconv.Itoa(port)
		clientConfig.Scheme = scheme
	})
}

//

// func getRegistrar() *Registrar {
// 	addr := cmdr.GetStringP(TAGS_PREFIX, "addr")
// 	if !strings.Contains(addr, ":") {
// 		addr = fmt.Sprintf("%v:%v", addr, cmdr.GetIntP(TAGS_PREFIX, "port"))
// 	}
// 	return getRegistrarImpl(addr, cmdr.GetStringP(TAGS_PREFIX, "scheme"))
// }
//
// func getRegistrarImpl(addr, scheme string) *Registrar {
// 	return &Registrar{
// 		Base: Base{
// 			FirstClient: MakeClientWithConfig(func(clientConfig *api.Config) {
// 				clientConfig.Address = addr
// 				clientConfig.Scheme = scheme
// 			}),
// 		},
// 		Clients:       nil,
// 		CurrentClient: nil,
// 	}
// }
//
// func getClient(host string, port int) *api.Client {
// 	return getClientImpl(host, port, cmdr.GetStringP(TAGS_PREFIX, "scheme"))
// }
//
// func getClientImpl(host string, port int, scheme string) *api.Client {
// 	return MakeClientWithConfig(func(clientConfig *api.Config) {
// 		clientConfig.Address = host + ":" + strconv.Itoa(port)
// 		clientConfig.Scheme = scheme
// 	})
// }
