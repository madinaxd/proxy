package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type response struct {
	ID      int               `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int64             `json:"length"`
}

var (
	requests  = make(map[int]request)
	responses = make(map[int]response)
	mu        sync.Mutex
	id        = 0
)

func main() {
	http.HandleFunc("/curl", handleRequest)
	http.ListenAndServe(":3000", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := parseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := storeRequest(request)

	response, err := sendRequest(id, *request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storeResponse(id, response)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func parseRequest(r *http.Request) (*request, error) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}

	var req request
	err = json.Unmarshal(requestBody, &req)
	if err != nil {
		return nil, fmt.Errorf("error parsing request body: %v", err)
	}

	return &req, nil
}

func sendRequest(id int, r request) (*response, error) {
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to the service: %v", err)
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to the service: %v", err)
	}

	response := &response{
		ID:      id,
		Status:  resp.StatusCode,
		Headers: make(map[string]string),
		Length:  resp.ContentLength,
	}

	for k, v := range resp.Header {
		response.Headers[k] = v[0]
	}

	return response, nil
}

func storeRequest(req *request) int {
	mu.Lock()
	defer mu.Unlock()
	id++
	requests[id] = *req
	return id
}

func storeResponse(id int, res *response) {
	mu.Lock()
	defer mu.Unlock()
	responses[id] = *res
}
