package repository

import (
	"context"

	"github.com/nimaibhat/GoCrudBookManagement/model"
)

type Repository interface {
	GetBook(ctx context.Context, id string) (model.Book, error)      // Updated method names and parameters
	CreateBook(ctx context.Context, in model.Book) (model.Book, error)
	UpdateBook(ctx context.Context, in model.Book) (model.Book, error)
	DeleteBook(ctx context.Context, id string) error
}
