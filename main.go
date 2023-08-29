package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Response is used to marshal the JSON response, now including the origin (client's IP address).
type Response struct {
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
}

// Extracts the IP address from the incoming request.
func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map to store headers.
	headerMap := make(map[string]string)

	// Iterate over all headers and populate the map.
	for k, v := range r.Header {
		headerMap[k] = v[0]
	}

	// Prepare the response struct, including the origin (client's IP address).
	resp := Response{Headers: headerMap, Origin: getIPAddress(r)}
	// print to console
	respJSON, respErr := json.Marshal(resp)
	if respErr != nil {
		fmt.Print(respErr)
	} else {
		fmt.Printf("respJSON: %v\n", string(respJSON))
	}

	// Set Content-Type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Convert the struct to JSON and write it to the response.
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Register the headersHandler function at the "/headers" path.
	http.HandleFunc("/headers", headersHandler)

	// Start the server on port 8080.
	fmt.Println("Server running at http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
