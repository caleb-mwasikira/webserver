package controllers

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) HomePage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("%v %v", req.Method, req.URL)
	res.Write([]byte("Welcome to the Home page"))
}

func (app *Application) AboutPage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("%v %v", req.Method, req.URL)
	res.Write([]byte("This is the About page"))
}
