package main

import (
	_ "locksim/boot"
	_ "locksim/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
