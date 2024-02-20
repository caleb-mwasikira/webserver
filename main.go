package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// parse command-line flags
	var (
		host       string
		port       int = 8080
		static_dir string
	)

	flag.StringVar(&host, "host", host, "HTTP network address")
	flag.IntVar(&port, "port", port, "Port number to run the web server")
	flag.StringVar(&static_dir, "static-dir", static_dir, "Path to static assets")
	flag.Parse()

	server, err := NewServer(host, uint16(port), static_dir, true)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(*server.LogWriter)

	router := httprouter.New()
	router.GET("/", homePage)
	router.GET("/about", aboutPage)
	router.ServeFiles("/static/*filepath", server.StaticDir)

	log.Printf("Starting server on %v ...\n", server.Addr())

	err = http.ListenAndServe(server.Addr(), router)
	if err != nil {
		log.Fatalf("Failed to start server due to error: %v", err)
	}
}

func homePage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("%v %v", req.Method, req.URL)
	res.Write([]byte("Welcome to the Home page"))
}

func aboutPage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("%v %v", req.Method, req.URL)
	res.Write([]byte("This is the About page"))
}
