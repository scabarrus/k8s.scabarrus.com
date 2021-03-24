package validate

import (
	"net/http"
	"encoding/json"
	. "k8s.io/api/admission/v1" 
	"log"
)

type Validate struct{}

func (v Validate)Request(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-type","applicaiton/json")
	adm := AdmissionReview{}
	_ = json.NewDecoder(r.Body).Decode(&adm)//Retrieve body
	log.Println("full Request : ",adm)
	log.Println("Request : ",adm.Request)

}