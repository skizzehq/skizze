package main

import (
	"counts/counters"
	"counts/server"
)

func main() {
	server := server.New(counters.ManagerProxy)
	server.Run()
}
