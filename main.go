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
	db, err := database.DB()
	if err != nil {
		panic(err)
	}
	ur := repo.NewUserRepo(db)
	pr := &repo.GormPostRepo{DB: db}
	fr := &repo.GormFollowRepo{DB: db} // Assuming you have a follow repository

	instahandler := handler.NewInstaHandler(ur, pr, fr)

	router := mux.NewRouter()

	router.Use(middleware.CORS)
	router.Use(middleware.ContextUpdater)
	router.Use(middleware.Logger)

	router.HandleFunc("/health", instahandler.Health).Methods("GET")
	router.HandleFunc("/createAccount", instahandler.CreateAcc).Methods("POST")
	router.HandleFunc("/login", instahandler.Login).Methods("POST")

	protected := router.NewRoute().Subrouter()

	protected.Use(middleware.JWTAuthMiddleware)
	protected.HandleFunc("/getUser", instahandler.GetAllUser).Methods("GET")
	protected.HandleFunc("/follow", instahandler.Follow).Methods("POST")
	protected.HandleFunc("/followers", instahandler.GetFollowers).Methods("GET")
	protected.HandleFunc("/addComment/post", instahandler.AddComment).Methods("POST")
	protected.HandleFunc("/addComment/comment", instahandler.AddComment).Methods("POST")
	protected.HandleFunc("/createPost", instahandler.CreatePost).Methods("POST")
	protected.HandleFunc("/likePost", instahandler.LikePost).Methods("POST")
	protected.HandleFunc("/home", instahandler.ShowHomePage).Methods("GET")
	protected.HandleFunc("/post/comments", instahandler.PostGetComments).Methods("GET")

	http.ListenAndServe(":8080", router)
}
