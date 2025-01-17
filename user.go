package main

import (
	"reflect"

	"github.com/olivere/elastic"
)

const (
	USER_INDEX = "user"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int64  `json:"age"`
	Gender   string `json:"gender"`
}

func checkUser(username, password string) (bool, error) {
	query := elastic.NewTermQuery("username", username)
	searchResult, err := readFromES(query, USER_INDEX)
	if err != nil {
		return false, err
	}

	var utype User
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		if u, ok := item.(User); ok {
			if u.Password == password {
				return true, nil
			}
		}
	}
	return false, nil
}

func addUser(user *User) (bool, error) {
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := readFromES(query, USER_INDEX)
	if err != nil {
		return false, err
	}

	if searchResult.TotalHits() != 0 {
		return false, nil
	}

	err = saveToES(user, USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	return true, nil
}
