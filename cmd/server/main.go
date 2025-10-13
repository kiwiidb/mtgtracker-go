package main

import (
	"context"
	"log"
	"mtgtracker/internal/core"
	"mtgtracker/internal/events"
	"mtgtracker/internal/feed"
	"mtgtracker/internal/follows"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/notification"
	"mtgtracker/pkg/moxfield"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/kiwiidb/utils/pkg/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("starting program")

	// Initialize Firebase app
	ctx := context.Background()
	var authClient *auth.Client
	if os.Getenv("FIREBASE_CONFIG") != "" {
		log.Println("initializing firebase")
		app, err := firebase.NewApp(ctx, nil)
		if err != nil {
			log.Fatal("failed to initialize Firebase app", err)
		}

		authClient, err = app.Auth(ctx)
		if err != nil {
			log.Fatal("failed to initialize Firebase auth client", err)
		}
	}

	// Initialize the postgres database connection
	// a local postgres dsn mtgtracker, the dsn is:
	// export POSTGRES_DSN="host=localhost user=postgres password=postgres dbname=mtgtracker port=5432 sslmode=disable"
	// docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=mtgtracker -p 5432:5432 postgres
	// docker exec -it postgres psql -U postgres -d mtgtracker
	log.Println("initializing database")
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	// Initialize event bus
	log.Println("initializing event bus")
	eventBus := events.NewEventBus()

	// // Initialize the repositories
	notificationsRepo := notification.NewRepository(db)
	coreRepo := core.NewRepository(db)
	followRepo := follows.NewRepository(db)

	// // Initialize the S3 storage
	log.Println("initializing storage")
	storage := storage.InitStorage()

	// // Initialize the services
	coreService := core.NewService(coreRepo, storage, eventBus)
	notificationsSvc := notification.NewService(notificationsRepo, coreService)
	followService := follows.NewService(followRepo, coreService)
	feedService := feed.NewService(followRepo, coreRepo, coreService)
	moxfieldService := moxfield.NewService()

	// Register event handlers
	log.Println("registering event handlers")
	notificationHandlers := notification.NewEventHandlers(notificationsRepo, coreService)
	notificationHandlers.RegisterHandlers(eventBus)

	followHandlers := follows.NewEventHandlers(followRepo, coreService)
	followHandlers.RegisterHandlers(eventBus)

	// // Create a new HTTP server
	mux := http.NewServeMux()

	coreService.RegisterRoutes(mux)
	moxfieldService.RegisterRoutes(mux)
	notificationsSvc.RegisterRoutes(mux)
	followService.RegisterRoutes(mux)
	feedService.RegisterRoutes(mux)

	// add middleware chain
	handler := middleware.ApacheLogMw(mux)
	handler = middleware.CorsMw(handler)
	handler = middleware.JsonMw(handler)
	if authClient != nil {
		// Use Firebase Auth middleware if authClient is available
		handler = middleware.FirebaseAuthMw(authClient, handler)
	} else {
		// Use mock Firebase Auth middleware if authClient is not available
		handler = middleware.MockFirebaseAuthMw(handler)
	}

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
