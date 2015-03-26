package main

import (
	"os"
	"time"
	"fmt"
	"runtime"
	"net/http"
)


func DebugHandler(w http.ResponseWriter, r *http.Request) {
	debug = !debug
	w.Write([]byte("OK\n"))
}

func DumpStat() {
	plot := `===
{ "Name" : "line", "Height" : 600, "Width" : 1900, "ItemName" : ["HeapSys : bytes obtained from system", "HeapAlloc : bytes allocated and still in use", "HeapIdle : bytes in idle spans", "NumGC", "ReqTotal", "ReqOK"] }
---
`
	path := "memory.chart"
	tty, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Printf("Open file %v failed %v\n", err.Error())
		return
	}
	defer func() {
		tty.Close()
	}()

	tty.Write([]byte(plot))
	var m runtime.MemStats
	j := 0
	for {
		time.Sleep(time.Second)
		runtime.ReadMemStats(&m)
		fmt.Fprintf(tty, "%d %d %d %d %d %d %d\n", j, m.HeapSys, m.HeapAlloc, m.HeapIdle, m.NumGC, ReqTotal, ReqOK)
		j++
	}
}

