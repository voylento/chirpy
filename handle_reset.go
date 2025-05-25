package main

import (
	"net/http"
)
func (config *Config) HandleReset(w http.ResponseWriter, r *http.Request) {
	if config.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	config.hits.Store(0)
	config.db.DeleteAllUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state\n"))
}
