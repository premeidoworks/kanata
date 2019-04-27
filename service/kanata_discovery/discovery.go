package kanata_discovery

type Node struct {
	nodeId string

	Weight int `json:"weight"` // 1-65535, 0 default

	Address []byte `json:"address"`
	Port    uint16 `json:"port"`
}

type WatchEvent int

const (
	WatchEventCreate WatchEvent = iota + 1
	WatchEventDelete
)

type ServiceOpType int

const (
	OpServiceTypeCreate = iota + 1
	OpServiceTypeDelete
)

type ServiceDescription struct {
	Service string
	Version string
	Group   string
}
