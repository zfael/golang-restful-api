package main

import (
	"encoding/json"
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID      bson.ObjectId `json:"_id" bson:"_id"`
	Name    string        `json:"name"`
	Address Address       `json:"address"`
}

type Address struct {
	Street string `json:"street"`
	Apt    string `json:"apt"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user/add", saveUserRoute).Methods("POST")
	r.HandleFunc("/user/{id}", getUserRoute).Methods("GET")
	r.HandleFunc("/user/delete/{id}", deleteUserRoute).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9999", r))
}

//routes handlers
func saveUserRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Received POST on /user/ADD")

	user := new(User)
	decoder := json.NewDecoder(r.Body)
	error := decoder.Decode(&user)

	if error != nil {
		log.Println(error.Error())
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	newID := addUser(user)

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, newID)
}

func getUserRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Received GET on /user/{id}")

	vars := mux.Vars(r)
	id := vars["id"]

	user := findUserByID(id)

	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(userJSON))
}

func deleteUserRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Received DELETE on /user/delete")

	vars := mux.Vars(r)
	id := vars["id"]

	result := deleteUserByID(id)

	if result {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

/*mgo stuff*/
func addUser(user *User) bson.ObjectId {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//creating new Object Id
	id := bson.NewObjectId()
	user.ID = id

	c := session.DB("golab").C("user")
	err = c.Insert(user)
	if err != nil {
		log.Fatal(err)
	}

	return id
}

func findUserByID(id string) User {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("golab").C("user")

	user := User{}
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)
	if err != nil {
		log.Fatal(err)
		return User{}
	}

	return user
}

func deleteUserByID(id string) bool {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("golab").C("user")

	err = c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}
