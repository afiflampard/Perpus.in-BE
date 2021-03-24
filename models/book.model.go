package models

import (
	"fmt"
	"net/http"
	"onboarding/helpers"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name         string         `gorm:"column:name; type:varchar(255); not null" json:"name"`
	OrderDetails []OrderDetail  `gorm:"foreignKey:IDBuku"`
	Author       string         `gorm:"column:author; type:varchar(255); not null" json:"author"`
	Isbn         string         `gorm:"column:isbn; type: varchar(255); not null" json:"isbn"`
	Isbn13       string         `gorm:"column:isbn13; type:varchar(255) not null" json:"isbn13"`
	Genre        pq.StringArray `gorm:"column:genre; type:text[]" json:"genre"`
	Language     string         `gorm:"column:language; type:varchar(255)" json:"language"`
	DatePub      string         `gorm:"column:date_pub;" json:"date_pub"`
	Pages        uint           `gorm:"column:pages;" json:"pages"`
	Sinopsis     string         `gorm:"column:sinopsis; type:text" json:"sinopsis"`
	Price        uint64         `gorm:"column:price" json:"price"`
	Fineamt      uint64         `gorm:"column:denda" json:"denda"`
}

func (book *Book) Create(conn *gorm.DB, id uint) (helpers.MessageResponse, *Book) {
	var petugas User
	if err := conn.Model(&petugas).Preload("Role").Find(&petugas).First(&petugas, id).Error; err != nil {
		return *helpers.MessageResponses(false, http.StatusUnprocessableEntity, "Tidak bisa"), nil
	}
	var addBuku Book
	fmt.Println(petugas)
	if strings.ToLower(petugas.Role.Role) == "petugas" {
		addBuku = Book{
			Name:     book.Name,
			Author:   book.Author,
			Isbn:     book.Isbn,
			Isbn13:   book.Isbn13,
			Genre:    book.Genre,
			Language: book.Language,
			DatePub:  book.DatePub,
			Pages:    book.Pages,
			Sinopsis: book.Sinopsis,
			Price:    book.Price,
			Fineamt:  book.Fineamt,
		}
		err := conn.Debug().Create(&addBuku).Error
		if err != nil {
			return *helpers.MessageResponses(false, http.StatusBadRequest, "Cannot add Buku"), nil
		}
	} else {
		return *helpers.MessageResponses(false, http.StatusBadRequest, "User cannot add Buku"), nil
	}
	return *helpers.MessageResponses(true, http.StatusAccepted, "Successfully"), &addBuku
}

func (book *Book) GetBookById(conn *gorm.DB, id uint) (*Book, error) {
	if err := conn.First(&book, id).Error; err != nil {
		return nil, err
	}
	return book, nil
}
func (book *Book) GetAllBook(conn *gorm.DB) ([]Book, error) {
	var books []Book
	if err := conn.Find(&books).Error; err != nil {
		helpers.MessageResponses(false, http.StatusBadRequest, "Book not Found")
	}
	fmt.Println(book)
	return books, nil
}
func (book *Book) UpdateBook(conn *gorm.DB, idBook uint, idUser uint) (*Book, error) {
	var user User
	var updateBook Book
	if err := conn.Model(&user).Preload("Role").Find(&user).First(&user, idUser).Error; err != nil {
		return nil, err
	}
	if err := conn.First(&updateBook, idBook).Error; err != nil {
		return nil, err
	}
	const layoutFormat = "2006-01-02"
	t, _ := time.Parse(layoutFormat, book.DatePub)
	if strings.ToLower(user.Role.Role) == "petugas" {
		updateBook.Name = book.Name
		updateBook.Author = book.Author
		updateBook.Isbn = book.Isbn
		updateBook.Isbn13 = book.Isbn13
		updateBook.Language = book.Language
		updateBook.DatePub = t.String()
		updateBook.Pages = book.Pages
		updateBook.Genre = book.Genre
		updateBook.Sinopsis = book.Sinopsis
		updateBook.Price = book.Price
		updateBook.Fineamt = book.Fineamt
		conn.Save(&updateBook)
	}
	return &updateBook, nil
}

func (book *Book) NewestBook(conn *gorm.DB) ([]Book, error) {
	var books []Book
	if err := conn.Order("created_at desc").Find(&books).Limit(3).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (stock *Stock) PopulerBook(conn *gorm.DB) ([]Stock, error) {
	var stocks []Stock
	if err := conn.Find(&stocks).Order("qty desc").Limit(3).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}
