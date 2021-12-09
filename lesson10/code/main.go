package main

import (
	"context"
	"flag"
	"io"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/golang/glog"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"ccamp/lesson10/metrics"
)

func main() {

	loglevel, isset := os.LookupEnv("LOGLEVEL")
	if !isset {
		loglevel = "2"
	}
	flag.Set("v", loglevel)
	metrics.Register()

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	port, isset := os.LookupEnv("PORT")
	if !isset {
		port = "9900"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//gracefully shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.V(1).Info("Starting Myhttpserver...")
	log.V(2).Info("Log level is " + loglevel)

	<-done
	log.V(1).Info("Myhttpserver Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Myhttpserver Shutdown Failed:%+v", err)
	}
	log.V(1).Info("Myhttpserver Exited Properly")

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(0, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	ip := GetIP(r)
	log.V(2).Info("entering root handler: ip [" + ip + "] , status code: 200")

	for k, v := range r.Header {
		for _, v1 := range v {
			w.Header().Add(k, v1)
		}
	}
	w.Header().Add("VERSION", os.Getenv("VERSION"))
	io.WriteString(w, "Congratulation! You have hit the page successfully\n ")
	log.V(4).Infof("Respond in %d ms", delay)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "200\n")
	log.V(1).Info("entering healthz handler: ip [" + GetIP(r) + "] , status code: 200")
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
