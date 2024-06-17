package main

import (
	"fmt"
	"sync"
)

type BookKeeper struct {
	name        string
	id          string
	phoneNumber string
}

func NewBookKeeper(name, id, phoneNumber string) *BookKeeper {
	return &BookKeeper{
		name:        name,
		id:          id,
		phoneNumber: phoneNumber,
	}
}

func (bk *BookKeeper) getNameAndPhoneNumber() (string, string) {
	return bk.name, bk.phoneNumber
}

func main() {
	bk := NewBookKeeper("John", "123456789", "1234567890")

	onceValues := sync.OnceValues[string, string](bk.getNameAndPhoneNumber)
	for i := 0; i < 5; i++ {
		fmt.Println(onceValues())
	}
}
