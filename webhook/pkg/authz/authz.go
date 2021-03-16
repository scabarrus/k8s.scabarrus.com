package authz
//https://pkg.go.dev/k8s.io/api/authorization/v1#SubjectAccessReview
//https://pkg.go.dev/k8s.io/api/authentication/v1#TokenReview
import(
	. "k8s.io/api/authorization/v1" //https://pkg.go.dev/k8s.io/api/authorization/v1#SubjectAccessReview
	"net/http"
	"encoding/json"
	"log"
	//"strings"
)

// Authz defines one methods for validating user authorization on kubernetes by 
// authorization webhook. It's a proof of concept and should not be used in production.
// A simple rules is set, user in admin group can do anything and user in dev group can read data only.
type Authz struct {
	
}

// contains is an internal method that just return true/false if a str string is in s slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// CheckAuthn return a SubjectAccessReview with Allowed status set to true or false
func (a Authz) CheckAuthz(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-type","applicaiton/json")
	
	request := SubjectAccessReview{}
	response :=SubjectAccessReview{}
	_ = json.NewDecoder(r.Body).Decode(&request)//Retrieve body
	username := request.Spec.User 
	groups := request.Spec.Groups
	action := request.Spec.ResourceAttributes.Verb
	log.Println("REQUEST : ",request)

	if (contains(groups,"admin")){// admin can do anything
		response.Status.Allowed = true
	}else if (contains(groups,"dev")){
		if (action != "get" && action != "list"){ //dev can read data only
			response.Status.Allowed = false
		}else{
			response.Status.Allowed = true
		}
	}else{
		response.Status.Allowed = false
	}
	log.Printf("%s in group %s perform %s on resource %s is allowed == ",username,groups,action,request.Spec.ResourceAttributes.Resource,response.Status.Allowed)
	
	w.Header().Set("Content-type","application/json")
	json.NewEncoder(w).Encode(response)

}