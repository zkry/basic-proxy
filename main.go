package main

import (
	"flag"
	"io/ioutil"
	"net/http"

	"log"
)

func getEndpoint(scrapeRoot, path, username, password string) ([]byte, error) {
	req, err := http.NewRequest("GET", scrapeRoot+path[1:], nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error in performing request: ", err)
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func metricsProxy(root, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := getEndpoint(root, r.URL.Path, username, password)
		if err != nil {
			w.Write([]byte("# Error occured"))
			log.Println(err)
		}
		log.Println("Request for path ", r.URL.Path)
		w.Write(data)
	}
}

func main() {
	root := flag.String("root", "", "The root for the proxy to scan")
	username := flag.String("uname", "", "The username for basic auth")
	password := flag.String("pswd", "", "The password for basic auth")
	port := flag.String("port", ":8080", "The port for the server to run on")
	flag.Parse()

	// Params check
	if *root == "" {
		flag.Usage()
		return
	}
	if (*port)[0] != ':' {
		*port = ":" + *port
	}
	if (*root)[len(*root)-1] != '/' {
		*root = *root + "/"
	}

	http.HandleFunc("/", metricsProxy(*root, *username, *password))
	log.Println("Starting proxy on port", *port)
	http.ListenAndServe(*port, nil)
}
