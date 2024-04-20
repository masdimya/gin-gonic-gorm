init-dependency:
	go get -u github.com/gin-gonic/gin
	go get -u gorm.io/gorm
	go get -u gorm.io/driver/postgres
	go get -u github.com/joho/godotenv
run:
	go run main.go