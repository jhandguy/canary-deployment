package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type response struct {
	Node      string `json:"node"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
}

func handleRequest(response response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(response)
	}
}

func main() {
	port, _ := os.LookupEnv("HTTP_PORT")
	node, _ := os.LookupEnv("KUBERNETES_NODE")
	namespace, _ := os.LookupEnv("KUBERNETES_NAMESPACE")
	pod, _ := os.LookupEnv("KUBERNETES_POD")
	resp := response{
		Node:      node,
		Namespace: namespace,
		Pod:       pod,
	}
	http.HandleFunc("/", handleRequest(resp))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
