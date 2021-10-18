package main

import (
	"flag"
	"log"

	"io"

	"net/http"

	_ "net/http/pprof"

	"os"

	"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ip := GetIP(r)

	for k, v := range r.Header {
		for _, v1 := range v {
			w.Header().Add(k, v1)
		}
	}
	w.Header().Add("VERSION", os.Getenv("VERSION"))
	io.WriteString(w, "Congratulation! You have hit the page successfully\n")
	log.Println("entering root handler: ip [" + ip + "] , status code: 200")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	log.Println("entering root handler")
	w.WriteHeader(200)
	io.WriteString(w, "200\n")
	log.Println("entering healthz handler: ip [" + GetIP(r) + "] , status code: 200")
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
