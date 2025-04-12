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
	router.Use(middleware.Logger)

	router.HandleFunc("/health", instahandler.Health).Methods("GET")

	router.HandleFunc("/addComment/post", instahandler.AddComment).Methods("POST")
	router.HandleFunc("/addComment/comment", instahandler.AddComment).Methods("POST")
	router.HandleFunc("/createAccount", instahandler.CreateAcc).Methods("POST")
	router.HandleFunc("/createPost", instahandler.CreatePost).Methods("POST")
	router.HandleFunc("/login", instahandler.Login).Methods("POST")
	router.HandleFunc("/likePost", instahandler.LikePost).Methods("POST")
	router.HandleFunc("/home", instahandler.ShowHomePage).Methods("GET")
	router.HandleFunc("/post/comments", instahandler.PostGetComments).Methods("GET")

	http.ListenAndServe(":8080", router)
}
