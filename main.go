package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

var (
	endpoint string
	port     int
	timeout  int
	version  bool
)

func initFlag() {
	flag.StringVar(&endpoint, "endpoint", "", "Target Endpoint (e.g: https://google.com)")
	flag.IntVar(&port, "port", 8080, "Local TCP port to listen on")
	flag.IntVar(&timeout, "timeout", 15, "Set a request timeout. Specify in seconds, defaults to 15")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.Parse()

	if version {
		log.Printf("Current version is: v%.1f", 1.0)
		os.Exit(0)
	}
}

func main() {
	initFlag()

	targetUrl, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("error: failure while parsing endpoint: %s. Error: %s",
			endpoint, err.Error())
	}

	revProxy := httputil.NewSingleHostReverseProxy(targetUrl)

	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("x-proxy", "true")
		request.Host = targetUrl.Host
		revProxy.Transport = &transport{http.DefaultTransport}
		revProxy.ServeHTTP(writer, request)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  time.Duration(timeout) * time.Second,
		WriteTimeout: time.Duration(timeout) * time.Second,
		IdleTimeout:  time.Duration(timeout) * time.Second,
		Handler:      serverMux,
	}

	log.Printf("server start at %v", port)
	fmt.Printf("server has stopped: %+v", server.ListenAndServe())
}

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		msg := fmt.Sprintf("ERROR: %q", err)
		log.Print(msg)
		opErr, _ := err.(*net.OpError)
		msg = fmt.Sprintf("ERROR: %q", opErr)
		log.Print(msg)
		body := ioutil.NopCloser(bytes.NewReader([]byte(msg)))
		return &http.Response{
			Body:       body,
			StatusCode: 500,
		}, nil
	}
	log.Printf(
		"proxied: %s %s %d",
		req.Method,
		endpoint+req.RequestURI,
		resp.StatusCode,
	)
	resp.Header.Del("X-Frame-Options")
	return resp, nil
}
