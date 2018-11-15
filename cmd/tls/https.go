package main

import (
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

// curl -v -i 'https://yang.io:1443/hello'
// https://cloud.tencent.com/developer/article/1064243
// https://tls.ulfheim.net/
func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":1443", "tls-gen/basic/server/cert.pem", "tls-gen/basic/server/key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
