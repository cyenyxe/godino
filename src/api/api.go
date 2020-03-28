package dinoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
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
	apirouter.Methods("POST").HandlerFunc(api.addHandler)

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

func (api *Api) addHandler(writer http.ResponseWriter, request *http.Request) {
	var animal Animal
	err := json.NewDecoder(request.Body).Decode(&animal)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Provided animal is not correctly formatted")
		return
	}

	id, err := api.dbHandler.AddNewAnimal(&animal)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqlerr.ER_DUPLICATED_VALUE_IN_TYPE {
				writer.WriteHeader(http.StatusBadRequest)
				io.WriteString(writer, "The specimen already exists")
				return
			}
		}

		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, "An internal error occurred")
		fmt.Println(err)
	}

	fmt.Fprintf(writer, "New specimen created with ID %d", id)
}
