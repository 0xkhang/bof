package main

import "net/http"

func (cfg *apiConfig) HandlerResetDB(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`You do not have the permisson. imposter!`))
		return
	}

	if err := cfg.db.Reset(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`DB is clean!`))
}
