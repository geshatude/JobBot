package server

import ("log"
"net/http")

// The Server

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.HandleWebhook(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Request error:", err)
		return
	}
}

