package inbound

import (
	"github.com/gitaction/martini"
	"github.com/sunwei/gitaction/src/adapters/inbound/rpc"
)

var m = martini.Classic()

func Run() {
	m.Group("/git", func(router martini.Router) {
		router.Get("/hello", rpc.Hello)
	})
	
	m.Action(rpc.NewGitRouter().Handle)
	
	m.Run()
}
