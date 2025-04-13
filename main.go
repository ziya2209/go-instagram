package main

import (
	"instagram/database"
	"instagram/handler"
	"instagram/middleware"
	"instagram/repo"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db := database.DB()
	ur := repo.NewUserRepo(db)

	instahandler := handler.NewInstaHandler(ur)

	router := mux.NewRouter()

	router.Use(middleware.CORS)
	router.Use(middleware.ContextUpdater)
	router.Use(middleware.Logger)

	router.HandleFunc("/health", instahandler.Health).Methods("GET")

	protected := router.NewRoute().Subrouter()

	protected.Use(middleware.JWTAuthMiddleware)

	protected.HandleFunc("/addComment/post", instahandler.AddComment).Methods("POST")
	protected.HandleFunc("/addComment/comment", instahandler.AddComment).Methods("POST")
	protected.HandleFunc("/createAccount", instahandler.CreateAcc).Methods("POST")
	protected.HandleFunc("/createPost", instahandler.CreatePost).Methods("POST")
	protected.HandleFunc("/login", instahandler.Login).Methods("POST")
	protected.HandleFunc("/likePost", instahandler.LikePost).Methods("POST")
	protected.HandleFunc("/home", instahandler.ShowHomePage).Methods("GET")
	protected.HandleFunc("/post/comments", instahandler.PostGetComments).Methods("GET")

	http.ListenAndServe(":8080", router)
}
