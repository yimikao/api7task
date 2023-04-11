package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const GPTAPI = ""

type client struct {
	url *url.URL
	c   *http.Client
}

type options func(*client)

func withURL(u string) options {
	return func(c *client) {
		pu, err := url.Parse(u)
		if err != nil {
			log.Fatal("client url could not be parsed")
		}
		c.url = pu
	}
}

func newclient(opts ...options) *client {
	c := new(client)
	c.c = http.DefaultClient
	for _, o := range opts {
		o(c)
	}
	return c
}

type response struct {
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

func main() {

	clt := newclient(
		withURL(GPTAPI),
	)

	buf := bytes.NewBuffer(nil)
	req, err := http.NewRequest(
		http.MethodGet,
		clt.url.Host,
		buf,
	)
	if err != nil {
		log.Fatalf("couldnt create request, %s", err)
	}
	resp, err := clt.c.Do(req)
	if err != nil {
		log.Fatalf("error while sending request, %s", err)
	}
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Fatalf("error while reading response body, %s", err)
	} else {
		if _, err := buf.Write(b); err != nil {
			log.Fatalf("error appending response bytes to buffer, %s", err)
		}
	}
	r := new(response)
	err = json.Unmarshal(buf.Bytes(), r)
	if err != nil {
		log.Fatalf("couldnt parse response body, %s", err)
	}

}
