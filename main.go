package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "sync/atomic"
	"time"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
	"github.com/mkling/GoWebServ/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	dbQueries		*database.Queries
	platform		string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {

	godotenv.Load()
	const filepathRoot = "."
	const port = "8080"

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to open the Database")
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{}
	apiCfg.dbQueries = dbQueries
	apiCfg.platform = os.Getenv("PLATFORM")
	
	serveMux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serveMux.Handle("/app/", fsHandler)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerAddChirp)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerAddUsers)
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	serveMux.HandleFunc("GET /api/healthz", handleHealthz)

	server := &http.Server{
		Addr:		":" + port,
		Handler:	serveMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}


