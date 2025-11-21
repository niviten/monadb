package main

import (
	"log"

	"github.com/niviten/monadb"
)

func main() {
	monadb.Start()
	conn := monadb.GetDB()
	if conn == nil {
		log.Fatalln("db not init")
	}
}
