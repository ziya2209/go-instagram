package handler

import (
	"net/http"
)

type InstaHandler interface {
	// PostGetComments will return all the comments for the given post id
	// in nested form
	Health(http.ResponseWriter, *http.Request)
	PostGetComments(http.ResponseWriter, *http.Request)
	ShowHomePage(http.ResponseWriter, *http.Request)
	CreatePost(http.ResponseWriter, *http.Request)
	AddComment(http.ResponseWriter, *http.Request)
	LikePost(http.ResponseWriter, *http.Request)
	CreateAcc(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	GetAllUser(http.ResponseWriter, *http.Request)
	Follow(http.ResponseWriter, *http.Request)
	GetFollowers(http.ResponseWriter, *http.Request)
}
