package kanata_discovery

type Node struct {
	nodeId string

	Weight int `json:"weight"` // 1-65535, 0 default

	Address []byte `json:"address"`
	Port    uint16 `json:"port"`
}
