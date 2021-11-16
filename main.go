package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/obukhov/redis-inventory/cmd/app"
)

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()
	app.Execute()
}
