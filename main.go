package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Booking struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Room      string             `json:"room" bson:"room"`
	GuestName string             `json:"guestName" bson:"guestName"`
	CheckIn   string             `json:"checkIn" bson:"checkIn"`
	CheckOut  string             `json:"checkOut" bson:"checkOut"`
}

// Room struct updated to include an ImageUrl
type Room struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
	ImageUrl  string  `json:"imageUrl"` // New field for the image
}

var client *mongo.Client
var rooms []Room

// --- Helper function to update room availability ---
func updateRoomAvailability(roomID string, available bool) {
    for i := range rooms {
        if rooms[i].ID == roomID {
            rooms[i].Available = available
            break
        }
    }
}

// --- API Handlers ---

func getRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func getBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bookings []Booking
	collection := client.Database("hoteldb").Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var booking Booking
		cursor.Decode(&booking)
		bookings = append(bookings, booking)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(bookings)
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var booking Booking
	_ = json.NewDecoder(r.Body).Decode(&booking)

    // **IMPROVED LOGIC**: Check if room is available before booking
    for _, room := range rooms {
        if room.ID == booking.Room && !room.Available {
            http.Error(w, "Room is already booked", http.StatusConflict)
            return
        }
    }

	collection := client.Database("hoteldb").Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, booking)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

    // **IMPROVED LOGIC**: Update room availability after booking
    updateRoomAvailability(booking.Room, false)
	json.NewEncoder(w).Encode(result)
}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

    // Find the booking to get the room ID before deleting
    var bookingToDelete Booking
    collection := client.Database("hoteldb").Collection("bookings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
    
    err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bookingToDelete)
    if err != nil {
        http.Error(w, "Booking not found", http.StatusNotFound)
        return
    }

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

    // **IMPROVED LOGIC**: Update room availability after cancellation
    updateRoomAvailability(bookingToDelete.Room, true)
	json.NewEncoder(w).Encode(result)
}


func main() {
	fmt.Println("Starting the Go hotel booking server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Hardcoded rooms with real image URLs from Unsplash
	rooms = []Room{
		{ID: "101", Type: "Standard Room", Price: 120.00, Available: true, ImageUrl: "https://images.unsplash.com/photo-1566665797739-1674de7a421a?q=80&w=2874&auto=format&fit=crop"},
		{ID: "102", Type: "Standard Room", Price: 120.00, Available: true, ImageUrl: "https://images.unsplash.com/photo-1598605272254-16f0c0ecdfa5?q=80&w=2874&auto=format&fit=crop"},
		{ID: "201", Type: "Deluxe Room", Price: 180.00, Available: true, ImageUrl: "https://images.unsplash.com/photo-1590490360182-c33d57733427?q=80&w=2874&auto=format&fit=crop"},
		{ID: "202", Type: "Deluxe Room", Price: 180.00, Available: true, ImageUrl: "https://images.unsplash.com/photo-1568495248636-6432b97bd949?q=80&w=2874&auto=format&fit=crop"},
        {ID: "301", Type: "Executive Suite", Price: 250.00, Available: true, ImageUrl: "https://images.unsplash.com/photo-1611892440504-42a792e24d32?q=80&w=2940&auto=format&fit=crop"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/rooms", getRooms).Methods("GET")
	router.HandleFunc("/bookings", getBookings).Methods("GET")
	router.HandleFunc("/bookings", createBooking).Methods("POST")
	router.HandleFunc("/bookings/{id}", deleteBooking).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":8000", handler))
}
