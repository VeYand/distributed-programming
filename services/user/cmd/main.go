package main

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"user/pkg/app/query"
	"user/pkg/app/service"
	"user/pkg/app/session"
	repo "user/pkg/infrastructure/redis/repository"
	"user/pkg/transport"
)

func createHandler() *transport.Handler {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_USER_URL"),
		Password: os.Getenv("REDIS_USER_PASSWORD"),
	})
	userRepository := repo.NewUserRepository(rdb)
	userService := service.NewUserService(userRepository)
	userQueryService := query.NewUserQueryService(userRepository)
	userSession := session.NewUserSession(userQueryService)

	authKey, _ := base64.StdEncoding.DecodeString(os.Getenv("SESSION_AUTH_KEY"))
	encKey, _ := base64.StdEncoding.DecodeString(os.Getenv("SESSION_ENC_KEY"))
	cookieStore := sessions.NewCookieStore(authKey, encKey)

	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false, // true для HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	return transport.NewHandler(userService, userQueryService, userSession, cookieStore)
}

func setupPublicRoutes(handler *transport.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/user/signin", handler.GetSignInPage).Methods("GET")
	router.HandleFunc("/user/signin", handler.SignIn).Methods("POST")
	router.HandleFunc("/user/signup", handler.GetSignUpPage).Methods("GET")
	router.HandleFunc("/user/signup", handler.SignUp).Methods("POST")
	router.HandleFunc("/user/signout", handler.SignOut).Methods("GET", "POST")
	return router
}

func setupInternalRoutes(handler *transport.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/internal/auth/check", handler.CheckAuthHandler).Methods("GET")
	return router
}

func main() {
	handler := createHandler()

	publicRouter := setupPublicRoutes(handler)
	internalRouter := setupInternalRoutes(handler)

	go func() {
		log.Println("Main server started at :8082")
		if err := http.ListenAndServe(":8082", publicRouter); err != nil {
			log.Fatalf("Main server failed: %v", err)
		}
	}()

	log.Println("Internal auth server started at :8081")
	if err := http.ListenAndServe(":8081", internalRouter); err != nil {
		log.Fatalf("Internal server failed: %v", err)
	}
}
