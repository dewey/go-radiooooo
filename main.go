package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/dewey/go-radiooooo/scrape"
	"github.com/dewey/go-radiooooo/store"
)

func main() {
	c := http.DefaultClient
	archive := store.NewArchive("/Users/philipp/Desktop/stuff")
	scrape := scrape.API{
		Endpoint: "http://radiooooo.com",
		Client:   c,
		Storage:  archive,
	}

	success, err := scrape.Start()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(success)
}
