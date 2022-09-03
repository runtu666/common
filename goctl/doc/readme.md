# goctl-doc

Improvement of `goctl api doc`.

## Usage

example

```shell
goctl-doc -api ttt.api
```

> or use goctl-doc as a plugin of goctl:
>
> ```shell
> goctl api plugin -p goctl-doc="-o ~/ttt.md"  -dir . -api ttt.api
> ```

You can also use your own document template:

```shell
goctl-doc -api ttt.api -mainTemplate xxx.template -routesTemplate yyy.template
```

## Output Examp
ttt.api
```api
syntax = "v1"

info(
	title: "ttt"
	desc: " "
	author: "zrd"
	email: ""
	version: "v1"
)

type (
	GetNodesReq {
		Vendor *Vendor `form:"vendor"`
	}
	GetNodesResp {
		Items []NodeBase `json:"items"`
	}
	Vendor {
		ID string `form:"id"`
	}
	NodeBase {
		ID     string  `json:"id"`     // node ID
	}
)

@server(
	middleware: TTTAuth
	group: node
)
service ttt-api {
	@doc "get nod list"
	@handler GetNodeList
	get /v1/nodes (GetNodesReq) returns (GetNodesResp)
}
```

the output will be(ttt.md):

````markdown

# ttt
author: zrd

version: v1

## ttt-api

group: node

middleware: TTTAuth

### 1. get nod list

##### GET /v1/nodes

Request:

```golang
type GetNodesReq struct {
	Vendor *Vendor `form:"vendor"`
}
type Vendor struct {
	ID string `form:"id"`
}
```


Response:

```golang
type GetNodesResp struct {
	Items []NodeBase `json:"items"`
}
type NodeBase struct {
	ID string `json:"id"` // node ID
}
```
````
