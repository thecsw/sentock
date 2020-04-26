package main

import (
	"net"
	"net/http"
	"os"
)

var (
	privateIPBlocks []*net.IPNet
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the request is a get, let them through
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		// Check the secret keyword
		if r.Header.Get("webkey") != os.Getenv("WEB_KEY") {
			http.Error(w, "You are not authorized to perform this action.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
