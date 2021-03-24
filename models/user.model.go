package models

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"onboarding/helpers"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type LoginResponse struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type User struct {
	gorm.Model
	Email    string `gorm:"column:email; type:varchar(255); not null; unique" json:"email"`
	Name     string `gorm:"column:name; type:varchar(255); not null" json:"name"`
	Mobile   string `gorm:"column:mobile; type:varchar(255) not null " json:"mobile"`
	Address  string `gorm:"column:alamat; type:varchar(255)" json:"alamat"`
	Password string `gorm:"column:password; type:varchar(255); not null " json:"password"`
	Photo    string `gorm:"column:photo; type: varchar(255)" json:"photo"`
	RoleID   uint   `gorm:"column:roleId" json:"roleId"`
	Role     Role   `gorm:"foreignKey:RoleID"`
}

type Role struct {
	gorm.Model
	Role string `gorm:"column:role; type:varchar(255); not null" json:"role"`
}

func (user *User) Validate() (helpers.MessageResponse, bool) {
	if !strings.Contains(strings.ToLower(user.Email), "@") || len(user.Password) < 6 {
		return *helpers.MessageResponses(false, http.StatusUnprocessableEntity, "Harus dimasukkan"), false
	}
	return *helpers.MessageResponses(true, http.StatusAccepted, "requirement passed"), true

}

func (user *User) Create(conn *gorm.DB) (helpers.MessageResponse, *User) {
	if _, ok := user.Validate(); !ok {
		return *helpers.MessageResponses(false, http.StatusUnprocessableEntity, "Tidak bisa"), nil
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	addPerson := User{
		Email:    user.Email,
		Name:     user.Name,
		Mobile:   user.Mobile,
		Address:  user.Address,
		Password: string(pass),
		RoleID:   user.RoleID,
	}
	err = conn.Debug().Create(&addPerson).Error
	if err != nil {
		return *helpers.MessageResponses(false, http.StatusBadRequest, "Cannot add User"), nil
	}
	return *helpers.MessageResponses(true, 200, "Succesfully"), &addPerson
}

//Login is ...
func (user *User) Login(conn *gorm.DB, mobile string, email string, password string) (*LoginResponse, error) {
	if mobile == "" {
		email = strings.ToLower(email)
		if err := conn.Where("email = ?", email).First(&user).Error; err != nil {
			return nil, err
		}
	} else {
		if err := conn.Where("mobile = ?", mobile).First(&user).Error; err != nil {
			return nil, err
		}
	}
	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil {
		helpers.MessageResponses(false, http.StatusBadRequest, result.Error())
	}
	token, err := GenerateToken(user)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		ID:    uint64(user.ID),
		Email: user.Email,
		Token: token,
	}, nil
}

func (user *User) GetUserById(conn *gorm.DB, id uint) (*User, error) {
	if err := conn.Where("id = ?", id).Preload("Role").First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (user *User) GetUsers(conn *gorm.DB) ([]User, error) {
	var users []User
	if err := conn.Find(&users).Error; err != nil {
		helpers.MessageResponses(false, http.StatusBadRequest, "User Not Foud")
	}
	return users, nil
}

func (user *User) UpdateUsers(conn *gorm.DB, id uint) (*User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	var updateUser User
	if err := conn.Where("id = ?", id).First(&updateUser).Error; err != nil {
		return nil, err
	}
	if user.Email != "" {
		updateUser.Email = user.Email
		updateUser.Password = string(pass)
		updateUser.Mobile = user.Mobile
		updateUser.Address = user.Address
		conn.Save(&updateUser)
	}
	return &updateUser, nil
}

func (user *User) UpdatePhoto(conn *gorm.DB, id uint, r *http.Request) (*User, error) {
	if err := conn.First(&user, id).Error; err != nil {
		return nil, err
	}
	file, header, err := r.FormFile("file")
	filename := header.Filename
	out, err := os.Create("./tmp/" + filename)
	fmt.Print("Outnya", *out)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println(err)
	}

	user.Photo = out.Name()
	GetDB().Save(&user)
	return user, nil
}
func (user *User) DeleteUserByID(conn *gorm.DB, id uint) (*User, error) {
	if err := conn.Delete(&user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func GenerateToken(user *User) (string, error) {
	expTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		ID:    uint64(user.ID),
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("I love bee"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
