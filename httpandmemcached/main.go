package main

import (
	"flag"
	"github.com/kklis/gomemcache"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	_ "net/http/pprof"
)

var addr = flag.String("maddr", "127.0.0.1:11211", "The host and listening port of the memcached server.")
var poolSize = flag.Int("poolsize", 4096, "The connection pool size to the memcached server")
var pool *Pool // The memcached connection pool
var debug bool
var ReqTotal int32
var ReqOK int32

func handler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&ReqTotal, 1)
	post, err := ioutil.ReadAll(r.Body) //Read the http body
	if err != nil {
		w.WriteHeader(403)
		return
	}

	// URI : /memcached?key=a
	//log.Printf("path=[%v] uri=[%v] query=[%v] method=[%s]\n", r.URL.Path, r.URL.String(), r.URL.RawQuery, r.Method)
	//The output of above : path=[/memcached] uri=[/memcached?key=a] query=[key=a] method=[GET]
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

		atomic.AddInt32(&ReqOK, 1)
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
	atomic.AddInt32(&ReqOK, 1)
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
	go DumpStat()
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	flag.Parse()
	pool = New(*addr, *poolSize)
	if pool == nil {
		log.Fatalf("Connect memcached failed : %v\n", *addr)
	}
	http.HandleFunc("/memcached", handler)
	http.HandleFunc("/debug", DebugHandler)
	log.Fatal(http.ListenAndServe(":8091", nil))
}
