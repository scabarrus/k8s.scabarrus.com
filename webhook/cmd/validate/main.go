package main

import(
	"net/http"
	"log"
	"flag"
	. "k8s.scabarrus.com/webhook/pkg/validate"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
)


func main(){
	log.Println("Validate webhook")
	cert := flag.String("cert","server.crt","Certificate SSL path")
	key := flag.String("key","server.key","Private key")
	protocol := flag.String("protocol","https","protocol used: [http|https]")
	port := flag.Int("port",8080,"http port")
	flag.Parse()
	log.Println("certificate:  ",*cert," private key : ",*key," protocol : ",*protocol," port : ",*port)
	r:= mux.NewRouter()	
	validate := Validate{}
	r.HandleFunc("/",validate.Request).Methods("POST")
	log.Println("Start webserver")
	if (len(strings.TrimSpace(*protocol)) > 0){//Check if protocol is defined
		log.Println("Check which protocol")
		if(*protocol == "http"){ //if http
			log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port),r)) //load http listner
		}else if (*protocol == "https"){// if https
			log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(*port),*cert,*key,r)) //load http listner with TLS
		}else{
			log.Fatal("Protocol set is incorrect")
		}

	}else{
		log.Fatal("No protocol is set")
	}
}