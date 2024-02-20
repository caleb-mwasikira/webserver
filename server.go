package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/caleb-mwasikira/webservers/projectpath"
)

const (
	INFO_LOG = "info.log"
)

type Server struct {
	Host      string
	Port      uint16
	StaticDir http.FileSystem
	LogWriter *io.Writer
}

/*
checks if path exists and is a directory
*/
func dirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if !stat.IsDir() {
		return false, fmt.Errorf("path %v is a file not a directory", path)
	}
	return true, nil
}

func NewServer(host string, port uint16, static_dir string, enable_log bool) (*Server, error) {
	// if invalid host IP provided, set host IP address to anycast
	if addr := net.ParseIP(host); addr == nil {
		host = "0.0.0.0"
	}

	if ok, err := dirExists(static_dir); !ok {
		log.Println(err)
		static_dir = filepath.Join(projectpath.Root, "public")
	}

	var log_wrt io.Writer = os.Stdout

	// enable_log saves logs to the <projectpath.Root>/.logs directory
	if enable_log {
		log_dir := filepath.Join(projectpath.Root, ".logs/")

		// create log directory
		err := os.MkdirAll(log_dir, 0766)
		if err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}

		// create info.log file
		filename := filepath.Join(log_dir, INFO_LOG)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			return nil, err
		}

		log_wrt = io.MultiWriter(os.Stdout, file)
	}

	return &Server{
		Host:      "127.0.0.1",
		Port:      8080,
		StaticDir: http.Dir(static_dir),
		LogWriter: &log_wrt,
	}, nil
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}
