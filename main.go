package main

import (
	"fmt"
	"log"

	"github.com/hanako-san/go-hanako/repo"
)

func main() {
	db, err := repo.FetchFromInternet("Kanto")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", db)
}
