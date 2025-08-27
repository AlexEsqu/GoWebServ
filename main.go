package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "sync/atomic"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
	"/home/mkling/GoWebServ/internal/database"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	dbQueries		*database.Queries
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
	
	serveMux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serveMux.Handle("/app/", fsHandler)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	serveMux.HandleFunc("GET /api/healthz", handleHealthz)


	server := &http.Server{
		Addr:		":" + port,
		Handler:	serveMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}


