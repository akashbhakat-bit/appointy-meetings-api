package main

import (
	"context"
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Meeting struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	hotst_name string        				`json:"name,omitempty" bson:"name,omitempty"`
	email string             `json:"email,omitempty" bson:"email,omitempty"`
	rsvp  string             `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
	title string             `json:"title,omitempty" bson:"title,omitempty"`
	participant_name  string             `json:"participant,omitempty" bson:"participant,omitempty"`
	starttime  string             `json:"starttime,omitempty" bson:"starttime,omitempty"`
	endtime  string             `json:"endtime,omitempty" bson:"endtime,omitempty"`
	creation  string             `json:"creation,omitempty" bson:"creation,omitempty"`
}

func CreateMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Meeting
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("meeting").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}
func GetMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Meeting
	collection := client.Database("meeting").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}
func GetInfoEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Meeting
	collection := client.Database("meeting").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Meeting
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/create_meeting", CreateMeetingEndpoint).Methods("POST")
	router.HandleFunc("/meet", GetInfoEndpoint).Methods("GET")
	router.HandleFunc("/info/{id}", GetMeetingEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}