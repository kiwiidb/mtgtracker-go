package main

import (
	"log"
	"mtgtracker/internal/mtgtracker"
	"mtgtracker/internal/repository"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("starting program")
	// Initialize the postgres database connection
	// a local postgres dsn mtgtracker, the dsn is:
	// export POSTGRES_DSN="host=localhost user=postgres password=postgres dbname=mtgtracker port=5432 sslmode=disable"
	// docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=mtgtracker -p 5432:5432 postgres
	// docker ex
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	// // Initialize the repository
	repo := repository.NewRepository(db)

	// // Initialize the service
	service := mtgtracker.NewService(repo)
	// // Create a new HTTP server
	mux := http.NewServeMux()

	service.RegisterRoutes(mux)

	//serve static files
	mux.Handle("/", http.FileServer(http.Dir("static")))

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
