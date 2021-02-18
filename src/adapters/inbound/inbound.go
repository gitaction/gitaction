package inbound

import (
	"github.com/gitaction/martini"
	"github.com/sunwei/gitaction/src/adapters/inbound/rpc"
)

var m = martini.Classic()

func Run() {
	m.Action(rpc.NewGitRouter().Handle)
	
	m.Run()
}
