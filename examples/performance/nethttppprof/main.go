// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package main

// 通过net/http/pprof包采集
import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	log.Printf("Listen at port: 6060")
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	for {
		_ = fmt.Sprint("test sprint")
		time.Sleep(time.Millisecond)
	}
}
