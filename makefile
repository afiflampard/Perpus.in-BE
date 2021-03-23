install:
	go get github.com/gorilla/mux
	go get gorm.io/gorm
	go get gorm.io/driver/sqlite
	go get github.com/joho/godotenv
	go get github.com/go-pg/pg/v10
	go get gorm.io/driver/postgres
	go get github.com/dgrijalva/jwt-go
	go get github.com/badoux/checkmail
	go get golang.org/x/crypto/bcrypt
	go get github.com/lib/pq
	go get github.com/swaggo/swag/cmd/swag
	go get github.com/swaggo/http-swagger
	go get github.com/alecthomas/template
	

dev:
	nodemon --exec go run main.go --signal SIGTERM

down:
	docker-compose down -v 

up:
	docker-compose up -d
migrate:
	cd models
	go run base.model.go
