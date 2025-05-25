package main

import (
	"fmt"
	"net/http"
)

func (config *Config) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html> 
  <body> 
	  <h1>Welcome, Chirpy Admin</h1> 
		<p>Chirpy has been visited %d times!</p> 
	</body> 
</html>`, config.hits.Load())))
}

func (config *Config) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.hits.Add(1)
		next.ServeHTTP(w, r)
	})
}
