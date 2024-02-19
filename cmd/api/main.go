package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/caleb-mwasikira/webservers/projectpath"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Host      string
	Port      uint16
	StaticDir http.FileSystem
}

func NewServer(host string, port uint16, static_dir string) *Server {
	// if invalid host IP provided, set host IP address to anycast
	if addr := net.ParseIP(host); addr == nil {
		host = "0.0.0.0"
	}

	// verify static directory
	stat, err := os.Stat(static_dir)
	if err != nil {
		log.Printf("static directory path %v does not exist", static_dir)
		static_dir = filepath.Join(projectpath.Root, "public")
	} else {
		if !stat.IsDir() {
			log.Printf("static directory path %v is not a valid directory", static_dir)
			static_dir = filepath.Join(projectpath.Root, "public")
		}
	}

	return &Server{
		Host:      "127.0.0.1",
		Port:      8080,
		StaticDir: http.Dir(static_dir),
	}
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

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

	server := NewServer(host, uint16(port), static_dir)

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
