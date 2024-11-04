package case1

import (
	"github.com/go-kod/kod"
)

type errorConfig struct {
	A int `default:"sss"`
}

type test1ComponentDefaultErrorImpl struct {
	kod.Implements[test1ComponentDefaultError]
	kod.WithConfig[*errorConfig]
}

type test1ComponentGlobalDefaultErrorImpl struct {
	kod.Implements[test1ComponentGlobalDefaultError]
	kod.WithGlobalConfig[*errorConfig]
}
