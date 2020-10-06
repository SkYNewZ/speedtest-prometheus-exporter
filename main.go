package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	speedtestPath string
)

func init() {
	log.SetLevel(log.DebugLevel)

	// Parse the speetest path flag
	flag.StringVar(&speedtestPath, "speedtest-path", "/usr/bin/speedtest", "the speedtest cli path")
	flag.Parse()
}

func main() {
	// Register our collector
	speedtest := newspeedTestCollector()
	prometheus.MustRegister(speedtest)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/health", healthCheckHandler)
	r.Handle("/metrics", promhttp.Handler())

	// Set port to listen to
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up server
	srv := http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Start server
	go func() {
		log.Printf("Listening on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Start the speedtest scrape process
	startSpeedTestTicker()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("Shutting down server")
	os.Exit(0)
}
