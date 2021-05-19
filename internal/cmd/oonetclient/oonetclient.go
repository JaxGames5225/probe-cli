package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apex/log"
	"github.com/ooni/probe-cli/v3/internal/netplumbing"
)

/*
	Bottom line: this client can leverage existing Go functionality
	to easily connect to either socks5 or http proxy w/ auth.
*/

func main() {
	log.SetLevel(log.DebugLevel)
	//clnt := &http.Client{Transport: netplumbing.DefaultTransport}
	config := &netplumbing.Config{
		Logger: log.Log,
		/*
			Proxy: &url.URL{
				Scheme: "socks5",
				Host:   "127.0.0.1:9050",
			},
		*/
	}
	ctx := netplumbing.WithConfig(context.Background(), config)
	theader := &netplumbing.TraceHeader{}
	ctx = netplumbing.WithTraceHeader(ctx, theader)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.google.com", nil)
	if err != nil {
		log.WithError(err).Fatal("http.NewRequest failed")
	}
	resp, err := netplumbing.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.WithError(err).Fatal("clnt.Get failed")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Fatal("ioutil.ReadAll failed")
	}
	log.Infof("got %d bytes", len(data))
	for _, ev := range theader.MoveOut() {
		data, _ := json.Marshal(ev)
		fmt.Printf("%s\n", string(data))
	}
}
