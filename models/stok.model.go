package models

import (
	"fmt"
	"net/http"
	"onboarding/helpers"
	"strings"

	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	Qty      uint `gorm:"column:qty" json:"qty"`
	MaxStock uint `gorm:"column:maxStok" json:"max_stock"`
	BookID   uint `gorm:"column:book_id" json:"book_id"`
	Book     Book `gorm:"foreignKey:BookID"`
}

type RequestStock struct {
	Qty uint `json:"qty"`
}

func (stock *Stock) Validate(w http.ResponseWriter, bookID uint) (*Stock, error) {
	fmt.Println(stock.BookID)
	if bookID == 0 || stock.Qty == 0 {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	return stock, nil
}

func (stock *Stock) Create(conn *gorm.DB, w http.ResponseWriter, idUser uint, idBook uint) (*Stock, error) {
	var petugas User
	var book Book
	resp, err := stock.Validate(w, idBook)
	if err != nil {
		return nil, err
	}
	if err := conn.Model(&petugas).Preload("Role").Find(&petugas).First(&petugas, idUser).Error; err != nil {
		return nil, err
	}
	if err := conn.First(&book, idBook).Error; err != nil {
		return nil, err
	}
	var addStock Stock
	if strings.ToLower(petugas.Role.Role) == "petugas" {
		addStock = Stock{
			Qty:      resp.Qty,
			MaxStock: resp.Qty,
			BookID:   book.ID,
		}
		err := conn.Debug().Create(&addStock).Error
		if err != nil {
			return nil, err
		}
	}
	if err := conn.Preload("Book").Find(&addStock).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Stock Not Found")
	}
	return &addStock, nil
}
