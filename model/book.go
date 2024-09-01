package model

// import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
    ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
    Title  string  `json:"title"`
    Author string  `json:"author"`
    Price  float64 `json:"price"`
}
