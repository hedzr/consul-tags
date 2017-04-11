package consul

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/hashicorp/consul/api"
	"gopkg.in/urfave/cli.v2"
	"hedzr.com/consul-tags/util"
	"strconv"
	"strings"
)

func TagsList(c *cli.Context) error {
	registrar := getRegistrar(c)
	listServiceTags(registrar, c)
	return nil
}

func TagsToggle(c *cli.Context) error {
	registrar := getRegistrar(c)
	toggleServiceTags(registrar, c)
	return nil
}

func Tags(c *cli.Context) {
	registrar := getRegistrar(c)

	if true {
		cc := GetConsulApiEntryPoint(registrar)
		log.Debugf("GetConsulApiEntryPoint (via %s): %v\n", c.String("addr"), cc)

		for i, n := range c.FlagNames() {
			log.Debugf("    - flag name %d: %s, value: %v\n", i, n, c.Generic(n))
		}
	}

	modifyServiceTags(registrar, c)
}

func listServiceTags(registrar *Registrar, c *cli.Context) {
	name := c.String("name")
	if name != "" {
		listServiceTagsByName(registrar, c, name)
		return
	}
	id := c.String("id")
	if id != "" {
		listServiceTagsByID(registrar, c, id)
		return
	}
}

func listServiceTagsByName(registrar *Registrar, c *cli.Context, name string) {
	log.Debugf("List the tags of service '%s' at '%s'...", name, c.String("addr"))
	theServices, err := QueryService(name, registrar.FirstClient.Catalog())
	if err != nil {
		log.Critical(fmt.Errorf("Error: %v", err))
	} else {
		for _, s := range theServices {
			//log.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
			//	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
			//	s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(s.ServiceTags, ","))
		}
	}
}

func listServiceTagsByID(registrar *Registrar, c *cli.Context, id string) {
	log.Debugf("List the tags of service by id '%s'...", id)
	as0, err := QueryServiceByID(id, registrar.FirstClient)
	if err != nil {
		log.Critical(fmt.Errorf("Error: %v", err))
	} else {
		s, err1 := AgentServiceToCatalogService(as0, registrar.FirstClient)
		if err1 != nil {
			log.Critical(fmt.Errorf("Error: %v", err))
			return
		}

		//log.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
		//		1, s.ServiceID, s.ServiceAddress, s.ServicePort,
		//		s.ServiceTags, s.NodeMeta, s.Node, s.Address)

		fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(s.ServiceTags, ","))
	}
}

func modifyServiceTags(registrar *Registrar, c *cli.Context) {
	name := c.String("name")
	if name != "" {
		modifyServiceTagsByName(registrar, c, name)
		return
	}
	id := c.String("id")
	if id != "" {
		modifyServiceTagsByID(registrar, c, id)
		return
	}
}

func modifyServiceTagsByName(registrar *Registrar, c *cli.Context, name string) {
	log.Debugf("Modifying the tags of service '%s'...", name)
	theServices, err := QueryService(name, registrar.FirstClient.Catalog())
	if err != nil {
		log.Critical(fmt.Errorf("Error: %v", err))
	} else {
		//registrarId, registrarAddr, registrarPort := consulapi[0].ServiceID, consulapi[0].Address, consulapi[0].ServicePort
		//fmt.Printf("    Using '%s', %s:%d\n", userServices[0].ServiceID, userServices[0].Address, userServices[0].ServicePort)
		for _, s := range theServices {
			// 服务 s 所在的 Node
			cn := NodeToAgent(registrar, s.Node)
			// 节点 cn 的服务表中名为 "consulapi" 的服务
			as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
			// 从 consulapi 指示Agent（也即服务 s 所对应的 Agent），建立一个临时的 Client
			client := getClient(as.Address, as.Port, c)
			agentService := cn.Services[s.ServiceID]

			tags := modifyTags(s.ServiceTags, c.StringSlice("add"), c.StringSlice("rm"), c.String("delim"), c)

			//for _, t = range tags {
			//log.Debugf("    *** Tags: %v", tags)
			//}

			client.Agent().ServiceRegister(&api.AgentServiceRegistration{
				ID:                s.ServiceID,
				Name:              s.ServiceName,
				Tags:              tags,
				Port:              s.ServicePort,
				Address:           agentService.Address,
				EnableTagOverride: s.ServiceEnableTagOverride,
			})

			//重新载入s的等价物，才能得到新的tags集合，s.ServiceTags并不会自动更新为新集合
			sNew, _ := QueryServiceByID(s.ServiceID, client)

			//fmt.Printf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s:%d\n",
			//	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
			//	sNew.Tags, s.NodeMeta, s.Node, s.Address, as.Port)
			//log.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
			//	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
			//	s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(sNew.Tags, ","))
		}
	}
}

func modifyServiceTagsByID(registrar *Registrar, c *cli.Context, id string) {
	log.Debugf("Modifying the tags of service by id '%s'...", id)
	as0, err := QueryServiceByID(id, registrar.FirstClient)
	if err != nil {
		log.Critical(fmt.Errorf("Error: %v", err))
	} else {
		s, err1 := AgentServiceToCatalogService(as0, registrar.FirstClient)
		if err1 != nil {
			log.Critical(fmt.Errorf("Error: %v", err))
			return
		}
		// 服务 s 所在的 Node
		cn := NodeToAgent(registrar, s.Node)
		// 节点 cn 的服务表中名为 "consulapi" 的服务
		as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
		// 从 consulapi 指示Agent（也即服务 s 所对应的 Agent），建立一个临时的 Client
		client := getClient(as.Address, as.Port, c)
		agentService := cn.Services[id]

		tags := modifyTags(s.ServiceTags, c.StringSlice("add"), c.StringSlice("rm"), c.String("delim"), c)

		//for _, t = range tags {
		//	log.Debugf("    *** Tags: %v", tags)
		//}

		client.Agent().ServiceRegister(&api.AgentServiceRegistration{
			ID:                as0.ID,
			Name:              as0.Service,
			Tags:              tags,
			Port:              as0.Port,
			Address:           agentService.Address,
			EnableTagOverride: as0.EnableTagOverride,
		})

		//重新载入s的等价物，才能得到新的tags集合，s.ServiceTags并不会自动更新为新集合
		sNew, _ := QueryServiceByID(as0.ID, client)

		//fmt.Printf("    id='%s'[%s:%d], tags=%v, Node: %s,%s:%d\n",
		//	as0.ID, as0.Address, as0.Port,
		//	sNew.Tags, s.Node, as0.Address, as.Port)
		//log.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
		//	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
		//	s.ServiceTags, s.NodeMeta, s.Node, s.Address)
		fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(sNew.Tags, ","))
	}
}

func modifyTags(tags, addTags, removeTags []string, delim string, c *cli.Context) []string {
	if c.Bool("clear") {
		tags = make([]string, 0)
	}

	if !c.Bool("plain") {
		log.Debug("    --- in ext mode")
		list := make([]string, 0)
		for _, t := range removeTags {
			if c.Bool("string") {
				list = append(list, t)
			} else {
				for _, t1 := range strings.Split(t, ",") {
					list = append(list, t1)
				}
			}
		}
		for _, t := range list {
			log.Debugf("    --- slice: erasing %s", t)
			for {
				erased := false
				for i, v := range tags {
					va := strings.Split(v, delim)
					ta := strings.Split(t, delim)
					if len(ta) > 0 && strings.EqualFold(va[0], ta[0]) {
						tags = util.SliceEraseByIndex(tags, i)
						log.Debugf("      - slice: erased '%s%s%s'", va[0], delim, va[1])
						erased = true
						break
					}
				}
				if !erased {
					break
				}
			}
		}
		list = make([]string, 0)
		for _, t := range addTags {
			if c.Bool("string") {
				list = append(list, t)
			} else {
				for _, t1 := range strings.Split(t, ",") {
					list = append(list, t1)
				}
			}
		}
		for _, t := range list {
			log.Debugf("    --- slice: appending %s", t)
			matched := false
			for i, v := range tags {
				va := strings.Split(v, delim)
				ta := strings.Split(t, delim)
				if len(va) > 0 && strings.EqualFold(va[0], ta[0]) {
					tags[i] = t
					log.Debugf("      - slice: appended '%s%s%s'", va[0], delim, va[1])
					matched = true
				}
			}
			if !matched {
				tags = append(tags, t)
			}
		}

	} else {
		for _, t := range removeTags {
			tags = util.SliceErase(tags, t)
		}
		for _, t := range addTags {
			tags = append(tags, t)
		}
	}

	return tags
}

func toggleServiceTags(registrar *Registrar, c *cli.Context) {
	name := c.String("name")
	if name != "" {
		toggleServiceTagsByName(registrar, c, name)
		return
	}
	id := c.String("id")
	if id != "" {
		log.Critical(fmt.Errorf("toggle tags can be applied with --name but --id"))
		return
	}
}

func toggleServiceTagsByName(registrar *Registrar, c *cli.Context, name string) {
	log.Debugf("Toggle the tags of service '%s'...", name)
	theServices, err := QueryService(name, registrar.FirstClient.Catalog())
	if err != nil {
		log.Critical(fmt.Errorf("Error: %v", err))
	} else {
		newMaster := strings.Split(c.String("address"), ":")
		newMasterPort := 0
		if len(newMaster) > 1 {
			newMasterPort, err = strconv.Atoi(newMaster[1])
			if err != nil {
				log.Critical(fmt.Errorf("Error: %v", err))
				return
			}
		}
		masterTag := c.StringSlice("set")
		slaveTag := c.StringSlice("reset")
		if len(newMaster) == 0 {
			log.Critical(fmt.Errorf("--address to specify the master ip:port, it's NOT optional."))
			return
		}
		//if len(masterTag) == 1 {
		if len(masterTag) == 0 {
			log.Critical(fmt.Errorf("--set to specify the master tag, it's NOT optional."))
			return
		}
		if len(slaveTag) == 0 {
			log.Critical(fmt.Errorf("--reset to specify the slave tag, it's NOT optional."))
			return
		}

		for i, s := range theServices {
			fmt.Printf("    #%d. id='%s'[%s:%d #%s], tags=%v, meta=%v, Node: %s,%s\n",
				i, s.ServiceID, s.ServiceAddress, s.ServicePort, s.Address,
				s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			//for _, t := range s.ServiceTags {
			matched := strings.EqualFold(s.ServiceAddress, newMaster[0])
			if matched && len(newMaster) > 1 {
				matched = s.ServicePort == newMasterPort
			}
			tags := s.ServiceTags
			if matched {
				tags = modifyTags(tags, masterTag, slaveTag, c.String("delim"), c)
			} else {
				tags = modifyTags(tags, slaveTag, masterTag, c.String("delim"), c)
			}
			//}

			//for _, t = range tags {
			//log.Debugf("    *** Tags: %v\n", tags)
			//}

			cn := NodeToAgent(registrar, s.Node)
			as := CatalogNodeGetService(cn, SERVICE_CONSUL_API)
			log.Debugf("    %s=%v\n", SERVICE_CONSUL_API, as)
			client := getClient(as.Address, as.Port, c)
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
				log.Critical(fmt.Errorf("Error: %v", err))
				return
			}

			//csa, _, err := client.Catalog().Service(s.ServiceName, "", nil)
			//var cs *api.CatalogService = nil
			//for _, cs0 := range csa {
			//	if strings.EqualFold(cs0.ServiceID, s.ServiceID) {
			//		cs = cs0
			//		break
			//	}
			//}
			//if err != nil {
			//	log.Fatal(fmt.Errorf("Error: %v, %v", err, cs))
			//}
			//client.Catalog().Register(&api.CatalogRegistration{
			//	ID: s.ServiceID,
			//	Node: s.Node,
			//	Address: s.ServiceAddress,
			//	NodeMeta: s.NodeMeta,
			//	Service: agentService,
			//}, &api.WriteOptions{})

			//重新载入s的等价物，才能得到新的tags集合，s.ServiceTags并不会自动更新为新集合
			sNew, err := QueryServiceByID(s.ServiceID, client)
			if err != nil {
				log.Critical(fmt.Errorf("Error: %v", err))
				return
			}

			//fmt.Printf("    TAGS=%v\n\n", sNew.Tags)
			//log.Debugf("    #%d. id='%s'[%s:%d], tags=%v, meta=%v, Node: %s,%s\n",
			//	i, s.ServiceID, s.ServiceAddress, s.ServicePort,
			//	s.ServiceTags, s.NodeMeta, s.Node, s.Address)
			fmt.Printf("%s: %s\n", s.ServiceID, strings.Join(sNew.Tags, ","))
		}

		//
	}
}

func getRegistrar(c *cli.Context) *Registrar {
	return getRegistrarImpl(c.String("addr"), c.String("scheme"))
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

func getClient(host string, port int, c *cli.Context) *api.Client {
	return getClientImpl(host, port, c.String("scheme"))
}

func getClientImpl(host string, port int, scheme string) *api.Client {
	return MakeClientWithConfig(func(clientConfig *api.Config) {
		clientConfig.Address = host + ":" + strconv.Itoa(port)
		clientConfig.Scheme = scheme
	})
}
