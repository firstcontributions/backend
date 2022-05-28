package gateway

import "net/http"

type UserRequest struct {
	Name        *string `json:"name"`
	AccessToken *string `json:"access_token"`
	Handle      *string `json:"handle"`
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

}
