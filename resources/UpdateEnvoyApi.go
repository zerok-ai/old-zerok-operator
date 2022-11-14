package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var port = ":9093"

func RegisterUpdateEnvoyAPI() {
	router := mux.NewRouter()

	router.HandleFunc("/updenvoy/", UpdateEnvoyConfig).Methods("POST")

	fmt.Printf("Server at %v.\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func UpdateEnvoyConfig(w http.ResponseWriter, r *http.Request) {
	var updateEnvoyConfigRequest UpdateEnvoyConfRequest
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error while reading request body in update envoy request.")
	}

	json.Unmarshal(reqBody, &updateEnvoyConfigRequest)

}
