package main

import (
	"log"
	"net/http"
	"io/ioutil"
)

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
	http.HandleFunc("/echo", handler)
	http.HandleFunc("/proxyecho", handler)
	log.Fatal(http.ListenAndServe(":8091", nil))
}
