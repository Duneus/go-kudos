package gokudos

import "net/http"

type KudosService interface {
	Handler(rw http.ResponseWriter, r *http.Request)
}
