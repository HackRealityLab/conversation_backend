package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Body struct {
	Message string `json:"message"`
}

type IDResponse struct {
	ID string `json:"id"`
}

func WriteResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err := w.Write(body)
	if err != nil {
		log.Printf("error occurred when write body msg: %v", body)
		return
	}
}

func WriteMessage(w http.ResponseWriter, status int, msg string) {
	respBody := Body{Message: msg}
	body, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error occurred when marshalling response body: %v\n", respBody)
		return
	}

	WriteResponse(w, status, body)
}

func BadRequest(w http.ResponseWriter, text string) {
	WriteMessage(w, http.StatusBadRequest, text)
}

func NotFound(w http.ResponseWriter, text string) {
	WriteMessage(w, http.StatusNotFound, text)
}

func OKMessage(w http.ResponseWriter, text string) {
	WriteMessage(w, http.StatusOK, text)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Conflict(w http.ResponseWriter, text string) {
	WriteMessage(w, http.StatusConflict, text)
}

func InternalServerError(w http.ResponseWriter) {
	WriteMessage(w, http.StatusInternalServerError, "Internal server error")
}

func Unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	WriteMessage(w, http.StatusUnauthorized, "Unauthorized")
}

func IdResponse(w http.ResponseWriter, id string) {
	respBody := IDResponse{ID: id}
	body, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error occurred when marshalling response body: %v\n", respBody)
		return
	}

	WriteResponse(w, http.StatusOK, body)
}

func NotConfirmedError(w http.ResponseWriter) {
	WriteMessage(w, http.StatusForbidden, "You are not confirmed")
}
