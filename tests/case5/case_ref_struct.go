package case5

import "github.com/go-kod/kod"

type refStructImpl struct {
	kod.Implements[kod.Main]

	_ kod.Ref[testRefStruct1]
}

func (t *refStructImpl) Hello() string {
	return "Hello, World!"
}

type testRefStruct1 struct {
	kod.Implements[TestRefStruct1]
}
