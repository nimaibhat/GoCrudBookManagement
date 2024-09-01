package http

import (
	"context"
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nimaibhat/GoCrudBookManagement/model"
	"github.com/nimaibhat/GoCrudBookManagement/repository"
)

type Server struct {
	repository repository.Repository
}

func NewServer(repository repository.Repository) *Server {
	return &Server{repository: repository}
}

func (s *Server) GetBook(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        respondWithError(c, http.StatusBadRequest, "invalid argument: id")
        return
    }

    // Convert the string ID to ObjectID
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        respondWithError(c, http.StatusBadRequest, "invalid ID format")
        return
    }

    log.Println("Attempting to retrieve document with ID:", id)

    // Create a filter to query by ObjectID
    filter := bson.M{"_id": objectID}

    var book model.Book
    err = s.repository.GetBook(context.Background(), filter).Decode(&book)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            respondWithError(c, http.StatusNotFound, "book not found")
            return
        }
        respondWithError(c, http.StatusInternalServerError, "error retrieving book")
        return
    }

    respondWithSuccess(c, http.StatusOK, book)
}

func (s *Server) CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	createdBook, err := s.repository.CreateBook(context.Background(), book)
	if err != nil {
		handleError(c, err)
		return
	}

	respondWithSuccess(c, http.StatusCreated, createdBook)
}

func (s *Server) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "invalid argument id")
		return
	}

	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	book.ID = id

	updatedBook, err := s.repository.UpdateBook(context.Background(), book)
	if err != nil {
		handleError(c, err)
		return
	}

	respondWithSuccess(c, http.StatusOK, updatedBook)
}

func (s *Server) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "invalid argument id")
		return
	}

	err := s.repository.DeleteBook(context.Background(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	respondWithSuccess(c, http.StatusNoContent, nil)
}

func respondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"status": statusCode, "result": data})
}

func respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"status": statusCode, "error": message})
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrBookNotFound) {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}
	respondWithError(c, http.StatusInternalServerError, "internal server error")
}
