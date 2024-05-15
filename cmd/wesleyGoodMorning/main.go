package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {

	port := 80
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("Listening on port %d\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've playing %s\n", r.URL.Path)
	})

	// Start serving, using the listener we created
	log.Fatal(http.Serve(listener, nil))
}
