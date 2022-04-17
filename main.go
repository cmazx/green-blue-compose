package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var healthy = true
var healthyMutex = sync.RWMutex{}

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong\n")
}
func health(w http.ResponseWriter, req *http.Request) {
	healthyMutex.RLock()
	if !healthy {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	healthyMutex.RUnlock()
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.Println("start")
	healthCheckInterval, err := time.ParseDuration(os.Getenv("TRAEFIK_HEALTH_INTERVAL"))
	if err != nil {
		panic(err)
	}
	wg, srv := startWebserver()

	//wait os signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop

	log.Println("stopping")

	healthyMutex.Lock()
	healthy = false
	healthyMutex.Unlock()

	//give traefik chance to detect unhealthy container
	time.Sleep(healthCheckInterval)

	log.Println("server shutdown")
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
	//wait for http requests done
	wg.Wait()
	log.Println("stopped")
}

func startWebserver() (*sync.WaitGroup, *http.Server) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	//start server
	srv := &http.Server{Addr: ":80"}
	http.HandleFunc("/health", health)
	http.HandleFunc("/ping", ping)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return wg, srv
}
