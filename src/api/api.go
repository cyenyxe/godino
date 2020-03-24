package dinoapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Api struct {
	dbHandler *DinoDatabaseHandler
}

// RunWebPortalAPI runs the API that serves data into the portal
func RunWebPortalAPI(hostname string, port int) error {
	handler, err := NewDinoDatabaseHandler("dinoadmin:dinoadmin@/dino?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Close()

	api := Api{handler}

	r := mux.NewRouter()
	apirouter := r.PathPrefix("/api/dinos").Subrouter()

	apirouter.Methods("GET").Path("/age/{age}").HandlerFunc(api.findByAgeHandler)
	apirouter.Methods("GET").HandlerFunc(api.findAllHandler)

	return http.ListenAndServe(hostname+":"+strconv.Itoa(port), r)
}

func (api *Api) findAllHandler(writer http.ResponseWriter, request *http.Request) {
	animals := api.dbHandler.FindAll()
	json.NewEncoder(writer).Encode(animals)
}

func (api *Api) findByAgeHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	ageVar, ok := vars["age"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Please provide a minimum age\n")
		return
	}

	age, err := strconv.Atoi(ageVar)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Age must be an integer number")
		return
	}

	animals := api.dbHandler.FindByAge(age)
	json.NewEncoder(writer).Encode(animals)
}
