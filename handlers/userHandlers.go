package handlers

import (
	"gopkg.in/mgo.v2/bson"
	"GoAPI/user"
	"errors"
	"io/ioutil"
	"net/http"
	"encoding/json"
)


func usersGetAll(w http.ResponseWriter, r *http.Request) {
	users, err := user.All()
	if err != nil {
		postError(w, http.StatusInternalServerError)
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"users" : users})
}

func bodyToUser(r *http.Request, u *user.User) error {
	if r.Body == nil {
		return errors.New("Request body is empty.")
	}

	if u == nil {
		return errors.New("A user is required.")
	}

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bd, u)
}

func usersPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Location", "/users/" + u.ID.Hex())
	w.WriteHeader(http.StatusCreated)
}