package authn
//https://pkg.go.dev/k8s.io/api/authorization/v1#SubjectAccessReview
//https://pkg.go.dev/k8s.io/api/authentication/v1#TokenReview
import(
	. "k8s.io/api/authentication/v1"
	"net/http"
	"encoding/json"
	"log"
	"strings"
)

// Authn defines one methods for validating user authentication on kubernetes by 
// authentication webhook. It's a proof of concept and should not be used in production.
// A user token should be in format <username>@<secret>:<group1>,<group2>,...
// Any value for username is allowed, but secret should be set to value expected by the webhook.
// Any value for groups is accepted
type Authn struct {
	Secret string
}

// CheckAuthn return a TokenReview with Authenticated status set to true or false
func (a Authn) CheckAuthn(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-type","applicaiton/json")
	request := TokenReview{}
	response :=TokenReview{}
	_ = json.NewDecoder(r.Body).Decode(&request)//Retrieve body
	token := request.Spec.Token 
	userInfo := strings.Split(token,":")
	userDetails :=strings.Split(userInfo[0],"@")
	if(len(userDetails) == 2){ //if userinfo contains username and secret
		username := userDetails[0]
		log.Println("username : ",username)
		secret := userDetails[1]
		response.Status.User.Username = username
		if secret ==  a.Secret { //if secret is correct then user is authenticated
			userGroups := strings.Split(userInfo[1],",")
			groups := make([]string,0)
			for _,v := range userGroups { //Parse userGroups to set then in slice
				groups = append(groups,v)
			}
			response.Status.Authenticated = true
			response.Status.User.Groups = groups
		
		}else{ // if secret not match
			log.Println("Secret is incorrect !")
			response.Status.Authenticated = false
			response.Status.Error= "Secret is incorrect for username "+username
		}
	}else{ //if userInfo doesn't contain username and secret
		log.Println("Secret is missing !")
		response.Status.Authenticated = false
		response.Status.Error= "Secret is missing !"

	}
	w.Header().Set("Content-type","application/json")
	json.NewEncoder(w).Encode(response)

}