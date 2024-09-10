module github.com/xhermitx/gitpulse-01/backend

go 1.22.4

require (
	github.com/go-sql-driver/mysql v1.7.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.11
)

require github.com/xhermitx/my-utils v0.0.0-20240910030802-724cdb129e72 // indirect

require (
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/google/uuid v1.6.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.26.0
	golang.org/x/text v0.17.0 // indirect
)
