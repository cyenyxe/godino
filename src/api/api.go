package dinoportal

import (
	"fmt"
	"net/http"
)

// RunWebPortalAPI runs the API that serves data into the portal
func RunWebPortalAPI(address string) error {
	http.HandleFunc("/", rootHandler)
	return http.ListenAndServe(address, nil)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Welcome to the dino portal running on %s", request.RemoteAddr)
}
