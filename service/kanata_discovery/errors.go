package kanata_discovery

import "errors"

var (
	ErrSessionExpired           = errors.New("session expired")
	ErrIllegalPathNodeCharacter = errors.New("illegal path node character")
	ErrPathNodeEmpty            = errors.New("path node empty")
	ErrNoServiceFound           = errors.New("no service found")
)
