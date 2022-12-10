package doauth

import (
	"reflect"
	"testing"
)

func TestDoACL(t *testing.T) {
	reqs := []aclReq{
		{sub: "alice", obj: "data1", act: "read"},
		{sub: "alice", obj: "data2", act: "read"},
		{sub: "alice", obj: "data2", act: "write"},
		{sub: "bob", obj: "data2", act: "read"},
		{sub: "bob", obj: "data2", act: "write"},
		{sub: "root", obj: "data2", act: "write"},
		{sub: "root", obj: "data1", act: "read"},
	}
	if !reflect.DeepEqual([]bool{true, true, false, false, true, true, true}, doACL(reqs)) {
		t.Error("与预期不符, 请检查")
	}
}

func TestDoRBAC(t *testing.T) {
	reqs := []reqRBAC{
		{sub: "alice", obj: "data1", act: "read"},
		{sub: "alice", obj: "data2", act: "read"},
		{sub: "alice", obj: "data2", act: "write"},
		{sub: "bob", obj: "data2", act: "read"},
	}
	if !reflect.DeepEqual([]bool{true, true, true, false}, doRBAC(reqs)) {
		t.Error("与预期不符, 请检查")
	}
}
