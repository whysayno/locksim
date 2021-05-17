package main

import (
	_ "locksim/boot"
	_ "locksim/router"

	"github.com/gogf/gf/frame/g"
)

// @title Swagger Locksim API
// @version 0.1.1
// @license.name Apache 2.0

func main() {
	g.Server().Run()
}
