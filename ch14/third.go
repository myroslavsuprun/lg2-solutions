package main

import (
	"context"
	"fmt"
	"net/http"
)

type Level string
type levelKey struct{}

const (
	Debug Level = "debug"
	Info  Level = "info"
)

func contextWithLevel(ctx context.Context, level string) context.Context {
	return context.WithValue(ctx, levelKey{}, Level(level))
}

func levelFromContext(ctx context.Context) (Level, bool) {
	v, ok := ctx.Value(levelKey{}).(Level)
	return v, ok
}

func Log(ctx context.Context, level Level, message string) {
	inLevel, ok := levelFromContext(ctx)
	if !ok {
		return
	}

	if level == Debug && inLevel == Debug {
		fmt.Println(message)
	}
	if level == Info && (inLevel == Debug || inLevel == Info) {
		fmt.Println(message)
	}
}

func LoggerLevelMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		ok := q.Has("log_level")
		if !ok {
			return
		}

		ll := q.Get("log_level")
		if lvl := Level(ll); lvl != Debug && lvl != Info {
			http.Error(w, "wrong log level", http.StatusBadRequest)
		}

		ctx := r.Context()
		ctx = contextWithLevel(ctx, ll)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})

}

func thirdTask() {}
