package lru

import (
	"fmt"
	"reflect"
	"testing"
)

type Integer int32

func (i Integer) Len() int {
	return 4
}

func TestGet(t *testing.T) {
	cache := New(0, nil)
	cache.Add("lij", Integer(32))
	lij, ok := cache.Get("lij")
	if !ok || !reflect.DeepEqual(lij.(Integer), Integer(2)) {
		fmt.Println("error")
		t.Fail()
	}
}
