package main

import "github.com/gorilla/mux"

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/register", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	r.HandleFunc("/api/posts/{category}", handlers.ListByCategory).Methods("GET")
	r.HandleFunc("/api/posts/", handlers.List).Methods("GET")
	r.HandleFunc("/api/posts/", handlers.Add).Methods("POST")
	r.HandleFunc("/api/post/{post_id}", handlers.PostById).Methods("GET")
	r.HandleFunc("/api/post/{post_id}", handlers.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{post_id}/{comment_id}", handlers.DeleteComment).Methods("DELETE")
	r.HandleFunc("/api/post/{post_id}/upvote", handlers.Upvote).Methods("GET")
	r.HandleFunc("/api/post/{post_id}/downvote", handlers.Downvote).Methods("GET")
	r.HandleFunc("/api/post/{post_id}", handlers.DeletePost).Methods("DELETE")
	r.HandleFunc("/api/user/{user_id}", handlers.ListByAuthor).Methods("GET")
}
