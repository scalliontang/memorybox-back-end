package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
)

var mySigningKey = []byte("secret")

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for upload")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]
	p := Post{
		User:    username.(string),
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	err = savePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post from Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to save post from Elasticsearch %v\n", err)
		return
	}
	fmt.Println("Post is saved successfully")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []Post
	var err error
	if user != "" {
		searchPostsByUser(user)
	} else {
		searchPostsByKeywords(keywords)
	}
	if err != nil {
		http.Error(w, "Failed to read post from Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from Elasticsearch %v\n", err)
		return
	}
	js, err := json.Marshal(posts)

	if err != nil {
		http.Error(w, "Failed to parse post into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse post into JSON format %v\n", err)
		return
	}
	w.Write(js)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for signin")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	exists, err := checkUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "User doesn't exists", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Cannot generate token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for signup")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v\n", err)
		return
	}

	if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		fmt.Printf("Invalid username or password\n")
		return
	}

	success, err := addUser(&user)
	if err != nil {
		http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
		return
	}

	if !success {
		http.Error(w, "User already exists", http.StatusBadRequest)
		fmt.Println("User already exists")
		return
	}
	fmt.Printf("User added successfully: %s.\n", user.Username)
}
