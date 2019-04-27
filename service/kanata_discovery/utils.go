package kanata_discovery

import (
	"errors"
	"fmt"
	"strings"
)

func checkPathNode(pathNode string) error {
	if len(pathNode) == 0 {
		return ErrPathNodeEmpty
	}
	data := []byte(pathNode)
	for i := 0; i < len(data); i++ {
		c := data[i]
		allow := c >= 'A' && c <= 'Z' ||
			c >= 'a' && c <= 'z' ||
			c >= '0' && c <= '9' ||
			c == '.' || c == '_' || c == '-'
		if !allow {
			return ErrIllegalPathNodeCharacter
		}
	}
	return nil
}

func makeServicePath(prefix, service, version, tag, node string) string {
	return fmt.Sprint(prefix, service, "/", version, "/", tag, "/", node)
}

func splitServicePath(path string) (service, version, tag, node string, err error) {
	// /{namespace}/service/{service}/{version}/{tag}/{node}
	splits := strings.Split(path, "/")
	if len(splits) != 7 {
		err = errors.New("path is invalid")
		return
	}
	service = splits[3]
	version = splits[4]
	tag = splits[5]
	return
}
