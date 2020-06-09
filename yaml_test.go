package config

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGetFromYaml(t *testing.T) {
	buf, _ := ioutil.ReadFile("testdata/2.yaml")
	yml := NewYaml()
	err := yml.SetFromBytes(buf)
	if err != nil {
		t.Fatal(err)
	}
	v, err := yml.Get("weixin.mp.appid")
	if err != nil {
		t.Fatal(err)
	}
	if hobby, ok := v.(string); ok {
		fmt.Println(hobby)
	} else {
		fmt.Println(reflect.TypeOf(v))
	}
}
