package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/SkYNewZ/speedtest-prometheus-exporter/internal/speedtest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, port uint16, logger *logrus.Logger) {
	// Register our collector
	c := newCollector()
	prometheus.MustRegister(c)

	r := http.NewServeMux()
	r.HandleFunc("/health", HandleHealth)
	r.Handle("/results", HandleResults(c))
	r.Handle("/metrics", promhttp.Handler())

	if port == 0 {
		port = 3100
	}

	// start the server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("listening: %s", err)
		}
	}()

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalln("server forced to shutdown: ", err)
	}

	logger.Println("server exiting")
}

func HandleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "OK")
}

func HandleResults(c *collector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var data speedtest.Result
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// register our data
		c.data = &data
	})
}
