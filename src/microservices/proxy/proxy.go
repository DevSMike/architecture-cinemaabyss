package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Route struct {
	Prefix string
	Target *url.URL
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}


func main() {
	moviesURL := env("MOVIES_SERVICE_URL", "")
	monolithURL := env("MONOLITH_URL", "")
	gradual := strconv.ParseBool(env("GRADUAL_MIGRATION", "false"))
	percent := strconv.Atoi(env("MOVIES_MIGRATION_PERCENT", "0"))

    moviesProxy = httputil.NewSingleHostReverseProxy(moviesURL)
    monolithProxy = httputil.NewSingleHostReverseProxy(monolithURL)

    handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    		req.Header.Set("X-Forwarded-Host", req.Host)
      		var p *httputil.ReverseProxy

    		switch {
    		case strings.HasPrefix(req.URL.Path, "/api/movies"):
    			switch {
    			case !gradual || percent <= 0:
    				p = monolithProxy
    			case percent >= 100:
    				p = moviesProxy
    			case rand.Intn(100) < percent:
    				p = moviesProxy
    			default:
    				p = monolithProxy
    			}
    		default:
    			p = monolithProxy
    		}

    		p.ServeHTTP(w, req)
    	})

    	addr := env("PORT", "8000")
    	log.Printf("strangler proxy on %s (gradual=%v, percent=%d)", addr, gradual, percent)
    	log.Fatal(http.ListenAndServe(addr, handler))
}

