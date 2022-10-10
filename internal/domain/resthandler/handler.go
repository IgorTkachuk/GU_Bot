package resthandler

import "github.com/julienschmidt/httprouter"

type Handler interface {
	Register(r *httprouter.Router)
}
