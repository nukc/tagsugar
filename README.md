# TagSugar

> 根据字段里的 StructField 设置进行相应的处理，使用方法就好像 json 和 beego.orm 那样。

<p align="center">
    <a href="https://travis-ci.org/nukc/tagsugar"><img src="https://img.shields.io/travis/nukc/tagsugar.svg?style=flat-square"/></a>
    <a href="https://github.com/nukc/tagsugar/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-Apache-brightgreen.svg?style=flat-square"/></a>
</p>

## Intention

数据库存储头像或者其它文件的时候保存的是一个路径或者一个 key，查询出来的时候，往往并没有拼接域名地址；如果每次都要自己手动拼接一次的话，真的感觉好累，为了不想多做不必要的逻辑判断和代码

还有些情况就是数据库字段保存的是一个 json 字符串，想要在输出之前转换成一个 json 对象或数组。



## Usage

### Download and install

```
go get github.com/nukc/tagsugar
```

```go
package main

import "github.com/nukc/tagsugar"

func main(){
    // if use url(http) options
    tagsugar.Http = "https://cdn.github.com/"

    tagsugar.Lick(&model)
}

```

### Tag options support list

- url(http) ``// += tagsugar.Http``
- initial ``// set a initial value ``
- assign_to(FiledB) ``// FiledB is other filed``
- assign_type(mode)  ``// mode list:``
  > - raw: default mode, assign_to FiledB a raw value
  > - bool: assign_to FiledB a bool value
  > - unmarshal: set the FiledB a json.Unmarshal(str, &obj) value

### Tag options sample

> url(http)

```go
type Model struct {
    Id    int
    Image string `ts:"url(http)"`
}

model := Model{Id: 1, Image: "test.png"}
tagsugar.Lick(&model)

log.Print(model.Image)
// https://cdn.github.com/test.png
```

> assign_type(bool)

```go
type Model struct {
	Id    int
	Sex   int8   `ts:"assign_to(IsMan);assign_type(bool)"`
	IsMan bool
}

model := Model{Id: 2, Sex: 1}
// IsMan: false
tagsugar.Lick(&model)
// IsMan : true
```

> assign_type(unmarshal)

```go
type Model struct {
	Id     int
	Json   string `ts:"assign_to(Object);assign_type(unmarshal)" json:"-"`
	Object interface{}
}

json := "{\"id\": 1, \"post\": 2}"
model := Model{Id: 3, Json: json}
// Object: <nil>
tagsugar.Lick(&model)
// Object: map[id:1 post:2]
```

```go
type Model struct {
	Id     int
	Json   string `ts:"assign_to(Post);assign_type(unmarshal)" json:"-"`
	Post   Post
}

type Post struct {
	Id   int
	Post int
}

json := "{\"id\": 2, \"post\": 6}"
model := Model{Id: 4, Json: json}
// Post: {0 0}
tagsugar.Lick(&model)
// Post: {2 6}
```

```go
type Model struct {
	Id     int
	Json   string `ts:"assign_to(Array);assign_type(unmarshal)" json:"-"`
	Array  []interface{}
}

json := "[{\"id\": 1, \"post\": 3},{\"id\": 2, \"post\": 66}]"
model := Model{Id: 4, Json: json}
// Array: []
tagsugar.Lick(&model)
// Array: [map[id:1 post:3] map[id:2 post:66]]
```



## License

Apache License, Version 2.0