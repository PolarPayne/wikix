package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func middlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()

		log.Printf("%v %v", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("request took %v", time.Now().Sub(t0))
	})
}

func middlewareTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orig := r.URL.Path
		trimmed := strings.TrimSuffix(orig, "/")

		if orig != trimmed && len(trimmed) > 0 {
			log.Printf("redirecting request to %v to %v", orig, trimmed)
			http.RedirectHandler(trimmed, http.StatusMovedPermanently).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
