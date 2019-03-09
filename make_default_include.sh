#!/bin/sh

cat << EOF > include/init.go
package include

import (
	_ "github.com/premeidoworks/components/kanatasupport"
)
EOF
