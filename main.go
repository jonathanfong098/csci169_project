package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ayush6624/go-chatgpt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"gopkg.in/mail.v2"

	"github.com/jonathanfong098/csci169project/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB         *database.Queries
	SmtpServer *smtpServer
	client     *chatgptClient
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	portNumber, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		// If there's an error in conversion, use the default port number
		portNumber = 587
	}

	smtpConfig := &smtpConfig{
		server:   os.Getenv("SMTP_SERVER"),
		port:     portNumber,
		user:     os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASS"),
	}

	smtp := &smtpServer{
		dialer: mail.NewDialer(smtpConfig.server, smtpConfig.port, smtpConfig.user, smtpConfig.password),
		config: smtpConfig,
	}

	key := os.Getenv("OPENAI_KEY")
	if key == "" {
		log.Fatal("OPENAI_KEY environment variable is not set")
	}
	client, err := chatgpt.NewClient(key)
	if err != nil {
		log.Fatal(err)
	}
	chatgptClient := &chatgptClient{
		client: client,
	}

	apiCfg := apiConfig{
		DB:         dbQueries,
		SmtpServer: smtp,
		client:     chatgptClient,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/users", apiCfg.handlerUsersCreate)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedCreate))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowCreate))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerPostsGet))
	v1Router.Get("/summarized_posts", apiCfg.middlewareAuth(apiCfg.handlerSummarizedPostsGet))

	v1Router.Put("/subscribe", apiCfg.middlewareAuth(apiCfg.handlerSubscribeUser))
	v1Router.Put("/unsubscribe", apiCfg.middlewareAuth(apiCfg.handlerUnsubscribeUser))

	v1Router.Put("/toggle_summarize_posts", apiCfg.middlewareAuth(apiCfg.handlerToggleSummarizePosts))

	v1Router.Get("/recommend_feeds", apiCfg.recommendFeeds)

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	smtp.startDailyEmails(dbQueries, chatgptClient)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
