package case1

import (
	"github.com/go-kod/kod"
)

type errorConfig struct {
	A int
}

type test1ComponentDefaultErrorImpl struct {
	kod.Implements[test1ComponentDefaultError]
}

type test1ComponentGlobalDefaultErrorImpl struct {
	kod.Implements[test1ComponentGlobalDefaultError]
}
