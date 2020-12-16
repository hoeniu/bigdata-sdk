package main

import (
	"net/http"

	"github.com/hoeniu/bigdata-sdk/monitor"
)

func main() {
	ab := monitor.Ambaris{
		IP:   "192.168.2.144",
		Port: "50070",
		Path: "/",
		//Username: "admin",
		//Password: "admin",
	}

	s := &http.Server{
		Addr:    ":8081",
		Handler: ab.Proxy(),
	}
	s.ListenAndServe()
}
