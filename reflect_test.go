package mr_smart

import (
	"reflect"
	"testing"
)

type msg struct {
	code int
}

func TestReflect(t *testing.T) {
	tp := reflect.TypeOf(&msg{})
	instance := reflect.New(tp.Elem()).Interface()
	t.Logf("%v", instance)
}
