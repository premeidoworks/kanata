package kanata_discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

const (
	_DefaultTag = "$default"
	_BufferSize = 1024
)

// path format:
// characters: A-Za-z . _ -
// service path format: /{namespace}/service/{service}/{version}/{tag}/{node}
// configuration format: /{namespace}/configuration/{group}/{config}
type discovery struct {
	serviceStoreLock     sync.RWMutex
	serviceBufferChannel chan struct {
		Key           string
		Value         string
		ServiceOpType ServiceOpType
	}

	// service/version -> tag -> struct 1. session -> Node  2. node list
	serviceMapping map[string]map[string]*struct {
		mapping map[string]*Node
		list    []*Node
	}

	servicePrefix       string
	configurationPrefix string

	quorumManager QuorumManager
}

func serviceMappingFirstKey(service, version string) string {
	return service + "/" + version
}

func NewDiscovery(endpoints []string, namespace string) (*discovery, error) {
	d := new(discovery)
	d.quorumManager = NewEtcdManager(endpoints)

	d.servicePrefix = fmt.Sprint("/", namespace, "/service/")
	d.configurationPrefix = fmt.Sprint("/", namespace, "/configuration/")
	d.serviceMapping = make(map[string]map[string]*struct {
		mapping map[string]*Node
		list    []*Node
	})
	d.serviceBufferChannel = make(chan struct {
		Key           string
		Value         string
		ServiceOpType ServiceOpType
	}, _BufferSize)

	return d, nil
}

func (this *discovery) Start() error {
	err := this.quorumManager.Start()
	if err != nil {
		return err
	}

	// service data init
	this.serviceStoreLock.Lock()
	// service change
	err = this.quorumManager.WatchServiceChange(this.servicePrefix, this.serviceChange)
	if err != nil {
		this.serviceStoreLock.Unlock()
		return err
	}
	//FIXME need to optimize memory use in order to avoid out of memory
	list, err := this.quorumManager.GetAllServices(this.servicePrefix)
	if err != nil {
		this.serviceStoreLock.Unlock()
		return err
	}
	err = initServices(this, list)
	if err != nil {
		this.serviceStoreLock.Unlock()
		return err
	}
	this.serviceStoreLock.Unlock()

	go this.processServiceChange()

	return nil
}

func initServices(this *discovery, list []struct {
	Key   string
	Value string
}) error {
	for _, v := range list {
		service, version, tag, node, err := splitServicePath(v.Key)
		if err != nil {
			return errors.New("split service path error")
		}
		tagMap, ok := this.serviceMapping[serviceMappingFirstKey(service, version)]
		if !ok {
			tagMap = make(map[string]*struct {
				mapping map[string]*Node
				list    []*Node
			})
			this.serviceMapping[serviceMappingFirstKey(service, version)] = tagMap
		}
		tagServices, ok := tagMap[tag]
		if !ok {
			tagServices = &struct {
				mapping map[string]*Node
				list    []*Node
			}{mapping: make(map[string]*Node), list: []*Node{}}
			tagMap[tag] = tagServices
		}

		nodeStruct := new(Node)
		err = json.Unmarshal([]byte(v.Value), nodeStruct)
		if err != nil {
			logErr("unmarshal error when initService, key:", v.Key, " value:"+v.Value)
			return errors.New("unmarshal service data error")
		}
		nodeStruct.nodeId = node
		tagServices.mapping[node] = nodeStruct
		tagServices.list = append(tagServices.list, nodeStruct)
	}
	return nil
}

func (this *discovery) processServiceChange() {
	bufferChannel := this.serviceBufferChannel
	mapping := this.serviceMapping
	for {
		//FIXME need add a short delay to batch update for performance
		nodeChange := <-bufferChannel
		service, version, tag, node, err := splitServicePath(nodeChange.Key)
		if err != nil {
			continue
		}
		this.serviceStoreLock.Lock()
		switch nodeChange.ServiceOpType {
		case OpServiceTypeCreate:
			tagMap, ok := mapping[serviceMappingFirstKey(service, version)]
			if !ok {
				tagMap = make(map[string]*struct {
					mapping map[string]*Node
					list    []*Node
				})
				mapping[serviceMappingFirstKey(service, version)] = tagMap
			}
			tagServices, ok := tagMap[tag]
			if !ok {
				tagServices = &struct {
					mapping map[string]*Node
					list    []*Node
				}{mapping: make(map[string]*Node), list: []*Node{}}
				tagMap[tag] = tagServices
			}

			nodeStruct := new(Node)
			err = json.Unmarshal([]byte(nodeChange.Value), nodeStruct)
			if err != nil {
				logErr("unmarshal service create error:", err)
				this.serviceStoreLock.Unlock()
				continue
			}
			nodeStruct.nodeId = node
			tagServices.mapping[node] = nodeStruct
			//FIXME need optimize
			newList := make([]*Node, 0, len(tagServices.list)+1)
			for _, v := range tagServices.list {
				if v.nodeId == node {
					continue
				}
				newList = append(newList, v)
			}
			newList = append(newList, nodeStruct)
			tagServices.list = newList
		case OpServiceTypeDelete:
			tagMap, ok := mapping[serviceMappingFirstKey(service, version)]
			if !ok {
				this.serviceStoreLock.Unlock()
				continue
			}
			tagServices, ok := tagMap[tag]
			if !ok {
				this.serviceStoreLock.Unlock()
				continue
			}
			delete(tagServices.mapping, node)
			//FIXME need optimize
			newList := make([]*Node, 0, len(tagServices.list)+1)
			for _, v := range tagServices.list {
				if v.nodeId == node {
					continue
				}
				newList = append(newList, v)
			}
			tagServices.list = newList
		}
		this.serviceStoreLock.Unlock()
	}
}

func (this *discovery) serviceChange(eventType WatchEvent, k, v string) {
	switch eventType {
	case WatchEventCreate:
		this.serviceBufferChannel <- struct {
			Key           string
			Value         string
			ServiceOpType ServiceOpType
		}{Key: k, Value: v, ServiceOpType: OpServiceTypeCreate}
	case WatchEventDelete:
		this.serviceBufferChannel <- struct {
			Key           string
			Value         string
			ServiceOpType ServiceOpType
		}{Key: k, Value: v, ServiceOpType: OpServiceTypeDelete}
	}
}

func (this *discovery) GracefulShutdown() error {
	return this.quorumManager.Shutdown()
}

func (this *discovery) AcquireSession() (int64, error) {
	return this.quorumManager.AcquireSession()
}

func (this *discovery) SessionHeatbeat(session int64) error {
	return this.quorumManager.SessionHeatbeat(session)
}

func (this *discovery) PublishService(node *Node, service, version string, tags []string, session int64) error {
	if err := checkPathNode(service); err != nil {
		return err
	}
	if err := checkPathNode(version); err != nil {
		return err
	}
	tagsNonDup := make(map[string]struct{})
	for _, v := range tags {
		if err := checkPathNode(v); err != nil {
			return err
		}
		tagsNonDup[v] = struct{}{}
	}
	tagsNonDup[_DefaultTag] = struct{}{}

	data, err := json.Marshal(node)
	if err != nil {
		return errors.New("marshal node data error")
	}

	nodeString := fmt.Sprint(session)
	for k := range tagsNonDup {
		path := makeServicePath(this.servicePrefix, service, version, k, nodeString)
		err := this.quorumManager.PutService(path, string(data), session)
		//FIXME need revert success operations?
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *discovery) DisableService(session int64) error {
	//TODO
	return nil
}

func (this *discovery) PickOneFromTag(service, version, tag string) (*Node, error) {
	if len(tag) == 0 {
		tag = _DefaultTag
	}
	this.serviceStoreLock.RLock()
	tagMap, ok := this.serviceMapping[serviceMappingFirstKey(service, version)]
	if !ok {
		this.serviceStoreLock.RUnlock()
		return nil, ErrNoServiceFound
	}
	tagService, ok := tagMap[tag]
	if !ok {
		this.serviceStoreLock.RUnlock()
		return nil, ErrNoServiceFound
	}
	cnt := len(tagService.list)
	if cnt <= 0 {
		this.serviceStoreLock.RUnlock()
		return nil, ErrNoServiceFound
	}
	node := tagService.list[rand.Intn(cnt)]
	this.serviceStoreLock.RUnlock()
	return node, nil
}

func (this *discovery) PickAllFromTag(service, version, tag string) ([]*Node, error) {
	if len(tag) == 0 {
		tag = _DefaultTag
	}
	this.serviceStoreLock.RLock()
	tagMap, ok := this.serviceMapping[serviceMappingFirstKey(service, version)]
	if !ok {
		this.serviceStoreLock.RUnlock()
		return nil, ErrNoServiceFound
	}
	tagService, ok := tagMap[tag]
	if !ok {
		this.serviceStoreLock.RUnlock()
		return nil, ErrNoServiceFound
	}
	cnt := len(tagService.list)
	if cnt <= 0 {
		this.serviceStoreLock.RUnlock()
		return nil, ErrNoServiceFound
	}
	nodes := make([]*Node, cnt)
	// make copy of services
	for i, v := range tagService.list {
		nodes[i] = v
	}
	this.serviceStoreLock.RUnlock()
	return nodes, nil
}

func (this *discovery) ListenChange(service, version, tag string) (<-chan []*Node, error) {
	if len(tag) == 0 {
		tag = _DefaultTag
	}
	//TODO
	return nil, nil
}
