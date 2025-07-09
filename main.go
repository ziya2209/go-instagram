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
	ur := repo.NewGormUserRepo(db)
	pr := repo.NewGormPostRepo(db)
	fr := repo.NewGormFollowRepo(db)
	cr := repo.NewGormCommentRepo(db)

	instahandler := handler.NewInstaHandler(ur, pr, fr, cr)

	ctrl := &controller{}
	ctrl.Start(instahandler)
}

type controller struct{}

func (c *controller) Start(instahandler handler.InstaHandler) {
	router := mux.NewRouter()

	// Use standard CORS middleware (recommended)
	// router.Use(middleware.CORS)

	// Alternative: Use custom CORS middleware if you have issues with file:// protocol
	// Uncomment the line below and comment out the line above if needed
	router.Use(middleware.CustomCORS)

	router.Use(middleware.ContextUpdater)
	router.Use(middleware.Logger)

	// Global OPTIONS handler for CORS preflight
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers are already set by the CustomCORS middleware
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/health", instahandler.Health).Methods("GET")
	router.HandleFunc("/createAccount", instahandler.CreateAcc).Methods("POST")
	router.HandleFunc("/login", instahandler.Login).Methods("POST")

	protected := router.NewRoute().Subrouter()

	protected.Use(middleware.JWTAuthMiddleware)
	protected.HandleFunc("/users", instahandler.GetAllUser).Methods("GET")
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
