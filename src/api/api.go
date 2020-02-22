package dinoportal

import (
	"fmt"
	"net/http"
	"strconv"
)

// RunWebPortalAPI runs the API that serves data into the portal
func RunWebPortalAPI(hostname string, port int) error {
	http.HandleFunc("/", rootHandler)
	return http.ListenAndServe(hostname+":"+strconv.Itoa(port), nil)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Welcome to the dino portal running on %s", request.RemoteAddr)
}
