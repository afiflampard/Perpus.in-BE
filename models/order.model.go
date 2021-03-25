package models

import (
	"fmt"
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
}

type RequestPinjam struct {
	TanggalKembali string     `json:"tanggal_kembali"`
	DetailBuku     DetailBuku `json:"detail_buku"`
}

type DetailBuku struct {
	BanyakBuku string `json:"banyak_buku"`
	IDBuku     string `json:"id_Buku"`
}

type ReturnBook struct {
	DetailBuku DetailBuku `json:"detail_buku"`
}

const layoutFormat = "2006-01-02"

func (borrow *RequestPinjam) PinjamBuku(conn *gorm.DB, idMember uint, w http.ResponseWriter) (Borrow, error) {
	var book []Book
	var member User
	var stock []Stock
	var tempIdBuku []string
	var tempJumlahBuku []string

	if strings.Contains(borrow.DetailBuku.IDBuku, ",") {
		tempIdBuku = append(tempIdBuku, strings.Split(borrow.DetailBuku.IDBuku, ",")...)
		tempJumlahBuku = append(tempJumlahBuku, strings.Split(borrow.DetailBuku.BanyakBuku, ",")...)
	} else {
		tempIdBuku = append(tempIdBuku, borrow.DetailBuku.IDBuku)
		tempJumlahBuku = append(tempJumlahBuku, borrow.DetailBuku.BanyakBuku)
	}

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
			for index := 0; index < len(book); index++ {
				u64, _ := strconv.ParseUint(tempJumlahBuku[index], 10, 32)
				pinjam = Borrow{
					TanggalPeminjaman: time.Now(),
					TanggalKembali:    t,
					IDUser:            member.ID,
					NoState:           1,
					Total:             Calculasi(uint(book[index].Price), "*", uint(u64)),
				}
				fmt.Println("Indexnya adalah", index)
				err := conn.Debug().Preload("User").Create(&pinjam).Error
				if err != nil {
					helpers.ResponseWithError(w, http.StatusBadRequest, "Cannot save Borrow")
				}

				orderdetail := OrderDetail{
					IDBorrow: pinjam.ID,
					IDBuku:   book[index].ID,
				}
				fmt.Println("Harusnya kesini")
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
	var borrows []Borrow
	var member User

	var orderDetail []OrderDetail
	if strings.Contains(borrow.DetailBuku.BanyakBuku, ",") {
		tempIdBuku = append(tempIdBuku, strings.Split(borrow.DetailBuku.IDBuku, ",")...)
		tempJumlahBuku = append(tempJumlahBuku, strings.Split(borrow.DetailBuku.BanyakBuku, ",")...)
	} else {
		tempIdBuku = append(tempIdBuku, borrow.DetailBuku.IDBuku)
		tempJumlahBuku = append(tempJumlahBuku, borrow.DetailBuku.BanyakBuku)
	}
	if err := conn.Model(&member).Preload("Role").Find(&member, idMember).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request User")
	}
	if err := conn.Model(&orderDetail).Where("buku_id IN ?", tempIdBuku).Preload("Borrow").Find(&orderDetail).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid orderDetail")
	}
	if err := conn.Model(&stock).Where("book_id IN ?", tempIdBuku).Find(&stock).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Stock	")
	}

	var idBorrow []uint
	for i := 0; i < len(orderDetail); i++ {
		idBorrow = append(idBorrow, orderDetail[i].IDBorrow)
	}
	if err := conn.Where("id IN ?", idBorrow).Find(&borrows).Error; err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Borrow	")
	}
	for index := 0; index < len(stock); index++ {
		if stock[index].Qty > stock[index].MaxStock {
			helpers.ResponseWithError(w, http.StatusBadRequest, "Stock Melebihi Max Stock")
		}
	}
	for indexCount := 0; indexCount < len(orderDetail); indexCount++ {
		u64, _ := strconv.ParseUint(tempJumlahBuku[indexCount], 10, 32)
		idBukus, _ := strconv.ParseUint(tempIdBuku[indexCount], 10, 32)
		borrows[indexCount].NoState = 2
		stock[indexCount].Qty = Calculasi(stock[indexCount].Qty, "+", uint(u64))
		conn.Save(&borrows[indexCount])
		conn.Save(&stock[indexCount])
		history := History{
			IDBuku:   uint(idBukus),
			IDBorrow: orderDetail[indexCount].Borrow.ID,
			NoState:  2,
		}
		err := conn.Debug().Create(&history).Error
		if err != nil {
			helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Save")
		}
	}
	return borrows, nil
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
