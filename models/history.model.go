package models

import "gorm.io/gorm"

type History struct {
	gorm.Model
	IDBorrow   uint       `gorm:"column:borrowId" json:"order_id"`
	Order      Borrow     `gorm:"foreignKey:IDBorrow"`
	IDBuku     uint       `gorm:"column:buku_id" json:"buku_id"`
	Buku       Book       `gorm:"foreignKey:IDBuku"`
	NoState    uint       `gorm:"column:state_no" json:"no_state"`
	OrderState OrderState `gorm:"foreignKey:NoState"`
}
