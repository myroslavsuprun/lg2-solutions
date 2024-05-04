package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func logger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return logger
}

type Time struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth int    `json:"day_of_month,int"`
	Month      string `json:"month"`
	Year       int    `json:"year,int"`
	Hour       int    `json:"hour,int"`
	Minute     int    `json:"minute,int"`
	Second     int    `json:"second,int"`
}

func rfc3339Server() {
	log := logger()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /time", HandleTime)

	wrapped := RequestLogger(log.WithGroup("request"))(mux)

	server := &http.Server{
		Addr:         ":3031",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      wrapped,
	}

	log.Info("started HTTP server on port 3031")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func HandleTime(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	accept := r.Header.Get("Accept")
	if accept == "application/json" {
		tS := Time{
			DayOfWeek:  t.Weekday().String(),
			DayOfMonth: t.Day(),
			Month:      t.Month().String(),
			Year:       t.Year(),
			Hour:       t.Hour(),
			Minute:     t.Minute(),
			Second:     t.Second(),
		}

		b, err := json.Marshal(tS)
		if err != nil {
			http.Error(w, "failed converting time", http.StatusInternalServerError)
			return
		}

		w.Write(b)
		return
	}

	w.Write([]byte(t.Format(time.RFC3339)))
}

func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("HTTP request", slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.String("ip", r.RemoteAddr))
			h.ServeHTTP(w, r)
		})
	}
}
