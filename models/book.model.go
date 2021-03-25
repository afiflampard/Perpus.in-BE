package models

import (
	"fmt"
	"net/http"
	"onboarding/helpers"
	"sort"
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

func (book *Book) Validasi(w http.ResponseWriter) (*Book, error) {
	if book.Name == "" || book.Author == "" || book.Isbn == "" || book.Isbn13 == "" || book.Language == "" || book.DatePub == "" || book.Pages == 0 || book.Sinopsis == "" || book.Price <= 0 || book.Fineamt <= 0 {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Reques")
	}
	return book, nil
}

func (book *Book) Create(conn *gorm.DB, id uint, w http.ResponseWriter) (helpers.MessageResponse, *Book) {
	var petugas User

	resp, err := book.Validasi(w)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}

	if err := conn.Model(&petugas).Preload("Role").Find(&petugas).First(&petugas, id).Error; err != nil {
		return *helpers.MessageResponses(false, http.StatusUnprocessableEntity, "Tidak bisa"), nil
	}
	var addBuku Book
	if strings.ToLower(petugas.Role.Role) == "petugas" {
		addBuku = Book{
			Name:     resp.Name,
			Author:   resp.Author,
			Isbn:     resp.Isbn,
			Isbn13:   resp.Isbn13,
			Genre:    resp.Genre,
			Language: resp.Language,
			DatePub:  resp.DatePub,
			Pages:    resp.Pages,
			Sinopsis: resp.Sinopsis,
			Price:    resp.Price,
			Fineamt:  resp.Fineamt,
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

func (history *History) PopulerBook(conn *gorm.DB) ([]Book, error) {

	var tempHistories []History
	if err := conn.Where("state_no = ?", 1).Preload("Buku").Preload("Order").Find(&tempHistories).Error; err != nil {
		return nil, err
	}
	var books []Book
	if len(tempHistories) < 3 {
		if err := conn.Find(&books).Error; err != nil {
			return nil, err
		}
		return books, nil
	}
	var tempID []uint
	for index := 0; index < len(tempHistories); index++ {
		tempID = append(tempID, tempHistories[index].IDBuku)
	}
	countElement := CountElement(tempID)

	keys := make([]uint, 0, len(countElement))
	for key := range countElement {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return countElement[keys[i]] > countElement[keys[j]]
	})

	var IdBuku []uint
	var IdHistory []uint
	var book []Book
	for _, key := range keys {
		IdBuku = append(IdBuku, key)
	}
	fmt.Println(IdBuku)
	if len(IdBuku) < 3 {
		if err := conn.Where("id IN ?", IdBuku).Find(&book).Error; err != nil {
			return nil, err
		}
	} else {
		for index := 0; index < 3; index++ {
			IdHistory = append(IdHistory, IdBuku[index])
		}
		if err := conn.Where("id IN ?", IdHistory).Find(&book).Error; err != nil {
			return nil, err
		}
	}
	return book, nil
}

func CountElement(history []uint) map[uint]uint {
	var dict = make(map[uint]uint)
	for _, num := range history {
		dict[num] = dict[num] + 1
	}
	return dict
}
