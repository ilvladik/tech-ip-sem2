package httpserver

import (
	"log"
	"net/http"
)

func StartRedirectServer(addr, httpsTarget string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		target := httpsTarget + r.URL.RequestURI()
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})

	go func() {
		log.Println("HTTP redirect server started on http://localhost" + addr)
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Fatal("redirect server failed:", err)
		}
	}()
}
