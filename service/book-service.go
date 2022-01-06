package service

import (
	"fmt"
	"gogin/dto"
	"gogin/entity"
	"gogin/repository"
	"log"

	"github.com/mashingan/smapping"
)

type BookService interface {
	Insert(book dto.BookCreatedDTO) entity.Book
	Update(book dto.BookUpdateDTO) entity.Book
	Delete(book entity.Book)
	All() []entity.Book
	FindByID(bookID uint64) entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRep repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRep,
	}
}

func (service *bookService) Insert(book dto.BookCreatedDTO) entity.Book {
	bookToInsert := entity.Book{}
	err := smapping.FillStruct(&bookToInsert, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed map %v : ", err)
	}
	return service.bookRepository.InsertBook(bookToInsert)
}

func (service *bookService) Update(book dto.BookUpdateDTO) entity.Book {
	bookToUpdate := entity.Book{}
	err := smapping.FillStruct(&bookToUpdate, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed map %v : ", err)
	}
	return service.bookRepository.UpdateBook(bookToUpdate)

}

func (service *bookService) Delete(book entity.Book) {
	service.bookRepository.DeleteBook(book)
}

func (service *bookService) All() []entity.Book {
	return service.bookRepository.AllBook()
}

func (service *bookService) FindByID(bookID uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookID)
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	book := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", book.UserID)
	return userID == id
}
