package main

import server "counts/server"

func main() {
	server := server.New()
	server.Run()
}
