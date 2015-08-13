package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header)
	buf, err := ioutil.ReadAll(r.Body) //Read the http body
	log.Printf("recv n=%v\n", len(buf))
	if err == nil {
		w.Write(buf)
		w.Write([]byte("<---backend--->"))
		return
	}

	w.WriteHeader(403)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header)
	buf, err := ioutil.ReadAll(r.Body) //Read the http body
	if err == nil {
		w.Write(buf)
		return
	}

	w.WriteHeader(403)
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	http.HandleFunc("/echo", handler)
	http.HandleFunc("/proxyecho", proxyHandler)
	log.Fatal(http.ListenAndServe(":8091", nil))
}
