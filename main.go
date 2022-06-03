package main

import (
	"fmt"

	"github.com/NarsilWorks-Inc/servicebase"
)

var (
	sb *servicebase.ServiceBase
)

func main() {

	var (
		ok bool
	)

	sb, ok = servicebase.CreateService()
	if !ok {
		for _, m := range sb.GetMessages() {
			fmt.Println(m)
		}
	}

	sb.Name = "PlainKV"
	sb.Version = "1.0"

	sb.AddMime(".txt", "text/plain")
	sb.Router.PathPrefix("/api/").Handler(PlainKVRequestHandler())

	sb.Serve()
}
