package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chatgp/gpt3"
)

const DOGSAPI = "https://dog.ceo/api/breeds/image/random"
const GPTURI = "/v1/images/generations"

type client struct {
	// url *url.URL
	c *http.Client
}

type options func(*client)

//	func withURL(u string) options {
//		return func(c *client) {
//			pu, err := url.Parse(u)
//			if err != nil {
//				log.Fatal("client url could not be parsed")
//			}
//			c.url = pu
//		}
//	}

func withRoundTripper(rt http.RoundTripper) options {
	return func(c *client) {
		c.c.Transport = rt
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

type test string

func (c test) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	byts, err := json.Marshal("hello world")
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(byts))
	return resp, nil
}

type gpt struct {
	c *gpt3.Client
}

func newGptClient(k string) *gpt {
	cli, _ := gpt3.NewClient(&gpt3.Options{
		ApiKey:  k,
		Timeout: 30 * time.Second,
		Debug:   true,
	})
	return &gpt{c: cli}
}
func (c *gpt) RoundTrip(req *http.Request) (*http.Response, error) {
	params := map[string]interface{}{
		"prompt":          "a beautiful girl with big eyes",
		"n":               1,
		"size":            "256x256",
		"response_format": "url",
	}
	rest, err := c.c.Post(GPTURI, params)
	if err != nil {
		return nil, err
	}
	bts, err := json.Marshal(rest.Get("data.0.url"))
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(bts)
	res := new(http.Response)
	res.Body = io.NopCloser(b)
	return res, nil
}

type router func(w http.ResponseWriter, r *http.Request)

func (rtr router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr(w, r)
}

func main() {
	var APIKEY = os.Getenv("API_KEY")
	if len(APIKEY) == 0 {
		APIKEY = "xxxxxxxxxxxxxxxxxxxxxxx"
	}

	clt := newclient(
		withRoundTripper(test("")),
	)

	buf := bytes.NewBuffer(nil)
	req, err := http.NewRequest(
		http.MethodGet,
		"",
		buf,
	)
	if err != nil {
		log.Fatalf("couldnt create request, %s", err)
	}
	resp, err := clt.c.Do(req)
	if err != nil {
		log.Fatalf("error while sending request, %s", err)
	}
	if b, err := io.ReadAll(resp.Body); err != nil {
		log.Fatalf("error while reading response body, %s", err)
	} else {
		if _, err := buf.Write(b); err != nil {
			log.Fatalf("error appending response bytes to buffer, %s", err)
		}
	}
	// fmt.Printf("response: %s", buf.String())

	h := router(func(w http.ResponseWriter, r *http.Request) {
		w.Write(buf.Bytes())
	})

	r := http.NewServeMux()
	r.Handle("/home", h)

	fmt.Println("...............")
	log.Fatal(http.ListenAndServe(":8080", r))
	// err = json.Unmarshal(buf.Bytes(), r)
	// if err != nil {
	// 	log.Fatalf("couldnt parse response body, %s", err)
	// }

}
