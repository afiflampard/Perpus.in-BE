package models

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"onboarding/helpers"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Borrow struct {
	gorm.Model
	TanggalPeminjaman time.Time  `gorm:"column:tanggal_peminjaman" json:"tanggal_peminjaman"`
	KodeBorrow        int        `gorm:"kode_borrow" json:"kode_borrow"`
	TanggalKembali    time.Time  `gorm:"column:tanggal_kembali" json:"tanggal_kembali"`
	IDUser            uint       `gorm:"column:user_id" json:"userId"`
	User              User       `gorm:"foreignKey:IDUser"`
	NoState           uint       `gorm:"column:state_no"`
	Total             uint       `gorm:"column:total" json:"total"`
	OrderState        OrderState `gorm:"foreignKey:NoState"`
}

type OrderState struct {
	ID   uint   `gorm:"column:id; primary_key; AUTO_INCREMENT" json:"id"`
	No   uint   `gorm:"column:state_no"`
	Name string `gorm:"column:state_name"`
}

type OrderDetail struct {
	gorm.Model
	IDBorrow uint   `gorm:"column:borrow_id"`
	IDBuku   uint   `gorm:"column:buku_id"`
	Buku     Book   `gorm:"foreignKey:IDBuku"`
	Borrow   Borrow `gorm:"foreignKey:IDBorrow"`
	Qty      uint   `gorm:"column:qty" json:"qty"`
}

type RequestPinjam struct {
	TanggalKembali string       `json:"tanggal_kembali"`
	DetailBuku     []DetailBuku `json:"detail_buku"`
}

type DetailBuku struct {
	BanyakBuku string `json:"banyak_buku"`
	IDBuku     string `json:"id_Buku"`
}

type ReturnBook struct {
	DetailBuku []DetailBuku `json:"detail_buku"`
}

const layoutFormat = "2006-01-02"

func (borrow *RequestPinjam) PinjamBuku(conn *gorm.DB, idMember uint, w http.ResponseWriter) (Borrow, error) {
	var book []Book
	var member User
	var stock []Stock
	var tempIdBuku []string
	var tempJumlahBuku []string

	if len(borrow.DetailBuku) <= 1 {
		tempIdBuku = append(tempIdBuku, borrow.DetailBuku[0].IDBuku)
		tempJumlahBuku = append(tempJumlahBuku, borrow.DetailBuku[0].BanyakBuku)
	} else {
		for index := 0; index < len(borrow.DetailBuku); index++ {
			tempIdBuku = append(tempIdBuku, borrow.DetailBuku[index].IDBuku)
			tempJumlahBuku = append(tempJumlahBuku, borrow.DetailBuku[index].BanyakBuku)
		}
	}

	fmt.Println(tempIdBuku)

	if err := conn.Where("id IN ?", tempIdBuku).Find(&book).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request Book")
	}

	if err := conn.Model(&member).Preload("Role").Find(&member, idMember).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request User")
	}
	if err := conn.Where("book_id IN ?", tempIdBuku).Find(&stock).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request Book")
	}

	t, _ := time.Parse(layoutFormat, borrow.TanggalKembali)
	var pinjam Borrow
	if strings.ToLower(member.Role.Role) == "member" {
		if len(book) == 1 {
			u64, _ := strconv.ParseUint(tempJumlahBuku[0], 10, 32)

			pinjam = Borrow{
				TanggalPeminjaman: time.Now(),
				KodeBorrow:        rand.Int(),
				TanggalKembali:    t,
				IDUser:            member.ID,
				NoState:           1,
				Total:             Calculasi(uint(book[0].Price), "*", uint(u64)),
			}
			err := conn.Debug().Create(&pinjam).Error
			if err != nil {
				helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Save Pinjam")
			}
			orderdetail := OrderDetail{
				IDBorrow: pinjam.ID,
				IDBuku:   book[0].ID,
				Qty:      uint(u64),
			}
			err = conn.Debug().Create(&orderdetail).Error
			if err != nil {
				helpers.ResponseWithError(w, http.StatusBadRequest, "Tidak dapat menyimpan Order Detail")
			}
			stock[0].Qty = Calculasi(stock[0].Qty, "-", uint(u64))
			conn.Save(&stock)
			history := History{
				IDBuku:   book[0].ID,
				IDBorrow: pinjam.ID,
				NoState:  pinjam.NoState,
			}
			err = conn.Debug().Create(&history).Error
			if err != nil {
				helpers.ResponseWithError(w, http.StatusBadRequest, "Tidak bisa Menyimpan History")
			}
		} else {
			kodeBorrow := rand.Int()
			for index := 0; index < len(book); index++ {
				u64, _ := strconv.ParseUint(tempJumlahBuku[index], 10, 32)
				pinjam = Borrow{
					KodeBorrow:        kodeBorrow,
					TanggalPeminjaman: time.Now(),
					TanggalKembali:    t,
					IDUser:            member.ID,
					NoState:           1,
					Total:             Calculasi(uint(book[index].Price), "*", uint(u64)),
				}

				err := conn.Debug().Preload("User").Create(&pinjam).Error
				if err != nil {
					helpers.ResponseWithError(w, http.StatusBadRequest, "Cannot save Borrow")
				}

				orderdetail := OrderDetail{
					IDBorrow: pinjam.ID,
					IDBuku:   book[index].ID,
					Qty:      uint(u64),
				}

				err = conn.Debug().Create(&orderdetail).Error
				if err != nil {
					helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Save orderDetail")
				}

				stock[index].Qty = Calculasi(stock[index].Qty, "-", uint(u64))
				conn.Save(&stock[index])

				history := History{

					IDBuku:   book[index].ID,
					IDBorrow: pinjam.ID,
					NoState:  pinjam.NoState,
				}
				err = conn.Debug().Create(&history).Preload("User").Error
				if err != nil {
					helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Save history")
				}
			}
		}

	}
	return pinjam, nil
}

func (borrow *ReturnBook) ReturnBook(conn *gorm.DB, idMember uint, w http.ResponseWriter) ([]Borrow, error) {
	var tempIdBuku []string
	var tempJumlahBuku []string
	var stock []Stock

	var member User

	var tempBorrows []Borrow
	var orderDetail []OrderDetail

	if len(borrow.DetailBuku) <= 1 {
		tempIdBuku = append(tempIdBuku, borrow.DetailBuku[0].IDBuku)
		tempJumlahBuku = append(tempJumlahBuku, borrow.DetailBuku[0].BanyakBuku)
	} else {
		for index := 0; index < len(borrow.DetailBuku); index++ {
			tempIdBuku = append(tempIdBuku, borrow.DetailBuku[index].IDBuku)
			tempJumlahBuku = append(tempJumlahBuku, borrow.DetailBuku[index].BanyakBuku)
		}
	}

	if err := conn.Where("user_id = ? AND state_no = ?", idMember, 1).Find(&tempBorrows).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid fetch borrows")
	}

	var tempBorrowId []uint
	for index := 0; index < len(tempBorrows); index++ {
		tempBorrowId = append(tempBorrowId, tempBorrows[index].ID)
	}

	if err := conn.Model(&member).Preload("Role").Find(&member, idMember).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request User")
	}
	if err := conn.Model(&orderDetail).Where("buku_id IN ? AND borrow_id IN ?", tempIdBuku, tempBorrowId).Preload("Borrow").Find(&orderDetail).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid orderDetail")
	}
	if err := conn.Model(&stock).Where("book_id IN ?", tempIdBuku).Find(&stock).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Stock	")
	}

	var idBorrow []uint
	for i := 0; i < len(orderDetail); i++ {
		idBorrow = append(idBorrow, orderDetail[i].IDBorrow)
	}

	for index := 0; index < len(stock); index++ {
		if stock[index].Qty > stock[index].MaxStock {
			helpers.ResponseWithError(w, http.StatusBadRequest, "Stock Melebihi Max Stock")
		}
	}
	var tempHistories []History

	for indexCount := 0; indexCount < len(borrow.DetailBuku); indexCount++ {
		u64, _ := strconv.ParseUint(tempJumlahBuku[indexCount], 10, 32)
		idBukus, _ := strconv.ParseUint(tempIdBuku[indexCount], 10, 32)

		tempBorrows[indexCount].NoState = 2
		stock[indexCount].Qty = Calculasi(stock[indexCount].Qty, "+", uint(u64))
		orderDetail[indexCount].Qty = orderDetail[indexCount].Qty - uint(u64)
		conn.Save(&orderDetail)
		conn.Save(&tempBorrows)
		conn.Save(&stock)
		tempHistories = append(tempHistories, History{
			IDBuku:   uint(idBukus),
			IDBorrow: orderDetail[indexCount].Borrow.ID,
			NoState:  tempBorrows[indexCount].NoState,
		})
	}

	for index := 0; index < len(tempHistories); index++ {
		fmt.Println("ID BUKU", tempHistories[index].IDBuku)
	}
	for _, history := range tempHistories {
		err := conn.Debug().Create(&history).Error
		if err != nil {
			log.Fatalf("Failed to create History")
		}
	}

	return tempBorrows, nil
}

func (borrow *OrderDetail) ListBorrow(conn *gorm.DB, w http.ResponseWriter) ([]OrderDetail, error) {
	var tempBorrows []OrderDetail
	if err := conn.Preload("Borrow.User").Preload("Borrow.OrderState").Preload("Buku").Find(&tempBorrows).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "OrderDetail Not Found")
	}

	var borrows []OrderDetail
	for index := 0; index < len(tempBorrows); index++ {
		if tempBorrows[index].Borrow.NoState == 1 {
			borrows = append(borrows, tempBorrows[index])
		}
	}

	return borrows, nil
}

func (borrow *OrderDetail) ListReturnBook(conn *gorm.DB, w http.ResponseWriter) ([]OrderDetail, error) {
	var tempReturnBooks []OrderDetail
	if err := conn.Preload("Borrow.User").Preload("Borrow.OrderState").Preload("Buku").Find(&tempReturnBooks).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "OrderDetail Not Found")
	}
	var borrows []OrderDetail
	for index := 0; index < len(tempReturnBooks); index++ {
		if tempReturnBooks[index].Borrow.NoState == 2 {
			borrows = append(borrows, tempReturnBooks...)
		}
	}
	return borrows, nil
}

func Calculasi(variabel1 uint, operand string, variabel2 uint) uint {
	switch operand {
	case "+":
		return variabel1 + variabel2
	case "-":
		return variabel1 - variabel2
	case "*":
		return variabel1 * variabel2
	default:
		return variabel1 / variabel2
	}
}
