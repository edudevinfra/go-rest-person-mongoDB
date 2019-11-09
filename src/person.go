package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const ( 
	DATABASE = "senai"
	COLLECTION = "people"
)


type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Contact   *Contact           `json:"contact,omitempty"`
}

//type allPersons []Person

type Contact struct {
	Address *Address `json:"address,omitempty"`
	Phone   *Phone   `json:"phone,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type Phone struct {
	Ddd   string `json:"ddd,omitempty"`
	Number string `json:"number,omitempty"`
}

func createPerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
	response.WriteHeader(201)//eu add
}

func readPerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	personID := mux.Vars(request)["id"]
	if len(personID) == 0 {
		retrivePerson(ctx, collection, response, request)
	} else {
		retriveOnePerson(personID, response, request)
	}
	

	json.NewEncoder(response).Encode(people)
}

func retriveOnePerson(personID string, response http.ResponseWriter, request *http.Request) {

	id, _ := primitive.ObjectIDFromHex(personID)
	var person Person
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	
}

func retrivePerson(ctx context.Context, collection *mongo.Collection,
	response http.ResponseWriter, request *http.Request) {
	var people []Person
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
}

func readContact(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	personID := mux.Vars(request)["id"]
	if len(personID) == 0 {
		retrivePerson(ctx, collection, response, request)
	} else {
		retriveOnePerson(personID, response, request)
	}

	json.NewEncoder(response).Encode(people)
}

func createContact(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var contact Contact
	_ = json.NewDecoder(request.Body).Decode(&contact)
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, contact)
	json.NewEncoder(response).Encode(result)
	response.WriteHeader(201)//eu add

}
func atualizaContact(contact Contact, personID string) {
//func atualizaContact(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	doc := db.Collection(COLLECTION).FindOneAndUpdate(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("id", personID),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
			bson.EC.String("address.city", contact.City),
			bson.EC.String("address.state", contact.State),
			bson.EC.String("phone.ddd", contact.Ddd),
			bson.EC.String("phone.number", contact.Number),
		),
	nil)
		
	fmt.Println(doc)
}

func updateContact(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")



}
func deleteContact(response http.ResponseWriter, request *http.Request) {
	var contact Contact
	_ = json.NewDecoder(r.Body).Decode(&contact)

	_, err := db.Collection(COLLECTION).DeleteOne(context.Background(), contact, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/person", createPerson).Methods("POST")
	router.HandleFunc("/person", readPerson).Methods("GET")
	router.HandleFunc("/person/{id}", readPerson).Methods("GET")
	
	router.HandleFunc("/person/{id}/contact", readContact).Methods("GET")
	router.HandleFunc("/person/{id}/contact", createContact).Methods("POST")
	router.HandleFunc("/person/{id}/contact", atualizaContact).Methods("PUT")
	router.HandleFunc("/person/{id}/contact", updateContact).Methods("PATCH")
	router.HandleFunc("/person/{id}/contact", deleteContact).Methods("DELETE")

	http.ListenAndServe(":12345", router)
}