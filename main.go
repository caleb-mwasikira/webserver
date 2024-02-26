package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	c "github.com/caleb-mwasikira/dfs/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func connectToDatabase() (*sql.DB, error) {
	var db *sql.DB

	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to load mysql environment configuration: %v", err)
	}

	// Initialise a new sql.DB object (which is not a database connection but
	// a pool of connections) based on a DSN(data source name) in the format
	// <username>:<password>@[protocol(address)]/<db-name>?[...parameters]
	// The parameter ?parseTime=True is a driver-specific parameter that informs
	// the mysql driver to convert SQL TIME and DATE fields to Go time.Time objects
	dsn := fmt.Sprintf("%v:%v@/%v?parseTime=True",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// The sql.Open function doesn’t actually create any connections, all
	// it does is initialize the pool for future use. Actual connections to the
	// database are established lazily, as and when needed for the first time.
	// So to verify that everything is set up correctly we need to use the
	// db.Ping method to create a connection and check for any errors.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
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

	server, err := NewServer(host, uint16(port), static_dir, true)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(*server.LogWriter)

	db, err := connectToDatabase()
	if err != nil {
		log.Fatalf("failed to connect to mysql database: %v", err)
	}
	defer db.Close()
	log.Println("connected to database...")

	// controllers/handlers
	app := c.NewApplication(db)

	router := httprouter.New()
	router.GET("/", app.HomePage)
	router.GET("/about", app.AboutPage)
	router.GET("/notes", app.GetAllNotes)
	router.POST("/notes", app.CreateNewNote)
	router.GET("/notes/:id", app.GetNote)
	router.ServeFiles("/static/*filepath", server.StaticDir)

	log.Printf("Starting server on %v ...\n", server.Addr())

	err = http.ListenAndServe(server.Addr(), router)
	if err != nil {
		log.Fatalf("Failed to start server due to error: %v", err)
	}
}
