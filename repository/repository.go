package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nimaibhat/GoCrudBookManagement/model"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{db: db}
}

func (r repository) GetBook(ctx context.Context, id string) (model.Book, error) {
	var out book
	err := r.db.Collection("books").FindOne(ctx, bson.M{"_id": id}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Book{}, ErrBookNotFound
		}
		return model.Book{}, err
	}
	return toModel(out), nil
}

func (r repository) CreateBook(ctx context.Context, book model.Book) (model.Book, error) {
	out, err := r.db.Collection("books").InsertOne(ctx, fromModel(book))
	if err != nil {
		return model.Book{}, err
	}
	book.ID = out.InsertedID.(primitive.ObjectID).String()
	return book, nil
}

func (r repository) UpdateBook(ctx context.Context, book model.Book) (model.Book, error) {
	in := bson.M{}
	if book.Title != "" {
		in["title"] = book.Title
	}
	if book.Author != "" {
		in["author"] = book.Author
	}
	if book.Price != 0 {
		in["price"] = book.Price
	}
	out, err := r.db.Collection("books").UpdateOne(ctx, bson.M{"_id": book.ID}, bson.M{"$set": in})
	if err != nil {
		return model.Book{}, err
	}
	if out.MatchedCount == 0 {
		return model.Book{}, ErrBookNotFound
	}
	return book, nil
}

func (r repository) DeleteBook(ctx context.Context, id string) error {
	out, err := r.db.Collection("books").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if out.DeletedCount == 0 {
		return ErrBookNotFound
	}
	return nil
}

type book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Price  float64            `bson:"price,omitempty"`
}

func fromModel(in model.Book) book {
	return book{
		Title:  in.Title,
		Author: in.Author,
		Price:  in.Price,
	}
}

func toModel(in book) model.Book {
	return model.Book{
		ID:     in.ID.String(),
		Title:  in.Title,
		Author: in.Author,
		Price:  in.Price,
	}
}
