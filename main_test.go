package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error("GET /ping error, your code is " + strconv.Itoa(w.Code))
	}
}

func TestPdfAPIRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	// 目前只能測試傳送無效連結
	var jsonStr = []byte(`{"apiKey":"123456789","url":"/ping","size":"A4"}`)
	req, _ := http.NewRequest("POST", "/pdf", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	// t.Error(w.Body.String())
	if w.Code != 400 {
		t.Error("POST /pdf error, your code is " + strconv.Itoa(w.Code))
	}
}
