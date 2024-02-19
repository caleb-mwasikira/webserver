package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/caleb-mwasikira/webservers/projectpath"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Host      string
	Port      uint16
	StaticDir http.FileSystem
}

func getServerConfiguration() *Server {
	host := os.Getenv("HOST")
	static_dir := filepath.Join(projectpath.Root, "public")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	if ok := net.ParseIP(host); ok == nil {
		host = "127.0.0.1"
	}

	return &Server{
		Host:      host,
		Port:      uint16(port),
		StaticDir: http.Dir(static_dir),
	}
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

func main() {
	server := getServerConfiguration()
	router := httprouter.New()

	router.GET("/", homePage)
	router.GET("/about", aboutPage)
	router.ServeFiles("/static/*filepath", server.StaticDir)

	log.Printf("Starting server on %v ...\n", server.Addr())
	err := http.ListenAndServe(server.Addr(), router)
	if err != nil {
		log.Fatalf("Failed to start server due to error: %v", err)
	}
}

func homePage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res.Write([]byte("Welcome to the Home page"))
}

func aboutPage(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res.Write([]byte("This is the About page"))
}
