package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	coursephone "example.com/course-phone-access"
)

func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		log.Fatal("INFRAI_API_KEY is required")
	}
	client := coursephone.NewClient(coursephone.DefaultBaseURL, key, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup/code", func(w http.ResponseWriter, r *http.Request) {
		var input coursephone.CodeRequest
		if !decode(w, r, &input) {
			return
		}
		if err := client.SendSignupCode(r.Context(), input); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "code_sent"})
	})
	mux.HandleFunc("POST /signup/verify", func(w http.ResponseWriter, r *http.Request) {
		var input coursephone.VerifyRequest
		if !decode(w, r, &input) {
			return
		}
		result, err := client.VerifyCourseAccess(r.Context(), input)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	})
	log.Printf("course phone access listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, err error) {
	var apiErr *coursephone.APIError
	if errors.As(err, &apiErr) && apiErr.HTTPStatus >= 400 && apiErr.HTTPStatus < 500 {
		writeJSON(w, apiErr.HTTPStatus, map[string]string{"error": apiErr.Code, "message": apiErr.Message})
		return
	}
	writeJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream_request_failed"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
