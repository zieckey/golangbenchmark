package main

import (
	"flag"
	"github.com/kklis/gomemcache"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("maddr", "127.0.0.1:11211", "The host and listening port of the memcached server.")

var pool *Pool // The memcached connection pool

func handler(w http.ResponseWriter, r *http.Request) {
	post, err := ioutil.ReadAll(r.Body) //Read the http body
	if err != nil {
		w.WriteHeader(403)
		return
	}

	// URI : /memcached?key=a
	//log.Printf("path=[%v] uri=[%v] query=[%v] method=[%s]\n", r.URL.Path, r.URL.String(), r.URL.RawQuery, r.Method)
	if len(r.URL.RawQuery) <= 4 {
		w.WriteHeader(403)
		return
	}

	key := r.URL.RawQuery[4:]
	conn := pool.Get()
	if conn == nil {
		goto Handler403
	}
	defer pool.Put(conn)

	// HTTP GET
	if r.Method == "GET" {
		val, _, err := conn.Get(key)
		if err == gomemcache.NotFoundError {
			goto Handler404
		}

		if err != nil {
			goto Handler403
		}

		w.Write(val)
		return
	}

	// HTTP POST
	if len(post) == 0 {
		goto Handler403
	}
	err = conn.Set(key, post, 0, 0)
	if err != nil {
		goto Handler403
	}
	w.Write([]byte("STORED\r\n"))
	return

Handler403:
	w.WriteHeader(403)
	return
Handler404:
	w.WriteHeader(404)
	return
}

func main() {
	flag.Parse()
	pool = New(*addr)
	if pool == nil {
		log.Fatal("Connect memcached failed : %v\n", *addr)
	}
	http.HandleFunc("/memcached", handler)
	log.Fatal(http.ListenAndServe(":8091", nil))
}
