package tagsugar

import (
	"log"
	"reflect"
	"testing"
)

type Model struct {
	Id     int
	Name   string `ts:"-"`
	Sex    int8   `ts:"assign_to(IsMan);assign_type(bool)"`
	IsMan  bool
	Image  string `ts:"url(http)"`
	Avatar string `ts:"host(cdn)"`

	Json   string `ts:"assign_to(Object);assign_type(unmarshal)"`
	Object interface{}

	PostJson string `ts:"assign_to(Post);assign_type(unmarshal)"`
	Post     Post

	ArrayJson string `ts:"assign_to(Array);assign_type(unmarshal)"`
	Array     []interface{}
}

type Post struct {
	Id   int
	Post int
}

func TestParseTag(t *testing.T) {
	model := Model{Id: 1, Name: "test"}
	v := reflect.ValueOf(model)

	p := v.Type()
	l := p.NumField()
	for i := 0; i < l; i++ {
		sf := p.Field(i)
		tagOptions := parseTag(sf.Tag.Get("ts"))
		log.Print(tagOptions)
	}

}
