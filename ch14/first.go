package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func firstTask() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 1)
		ctx := r.Context()

		_, ok := ctx.Deadline()
		if ok {
			http.Error(w, "timeout", http.StatusGatewayTimeout)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hi there"))

	})

	wrapped := TimeoutMiddleware(500)(mux)

	server := &http.Server{
		Addr:         ":3031",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      wrapped,
	}

	slog.Info("starting server on port 3031")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func TimeoutMiddleware(ms int64) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(ms))
			defer cancel()

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
			slog.Error(ctx.Err().Error())
		})
	}
}
