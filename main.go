package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dewey/go-radiooooo/scrape"
	"github.com/dewey/go-radiooooo/store"
	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	var t = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var c = &http.Client{
		Timeout:   time.Second * 10,
		Transport: t,
	}
	archive := store.NewArchive(logger, "/Users/philipp/Desktop/stuff")
	scrape := scrape.API{
		Endpoint: "http://radiooooo.com",
		Client:   c,
		Storage:  archive,
		Log:      logger,
	}

	success, err := scrape.Start()
	if err != nil {
		logger.Log("error", err)
	}
	ai, err := archive.GetArchiveInfo()
	if err != nil {
		logger.Log("msg", err)
	}
	fmt.Println(success)
	for _, c := range ai.Countries {
		fmt.Println(c.String())
		for _, d := range c.Decades {
			fmt.Println(d.String())
		}
	}
}
