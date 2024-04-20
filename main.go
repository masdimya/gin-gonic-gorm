package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseModel struct {
	CreatedAt time.Time `gorm:"->:false;column:created_at" json:"created_at"`
	DeletedAt time.Time `gorm:"->:false;column:deleted_at" json:"-"`
}

type User struct {
	ID       int    `gorm:"column:id; primary_key; not null" json:"id"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"password"`
	Name     string `gorm:"column:name" json:"name"`
	Address  string `gorm:"column:address" json:"address"`
	Phone    string `gorm:"column:phone" json:"phone"`
	Active   bool   `gorm:"column:active" json:"-"`
	BaseModel
}

func (User) TableName() string {
	return "user"
}

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error Load env")
	}

	dsn := os.Getenv("DB_DSN")

	var dbErr error

	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if dbErr != nil {
		fmt.Println("Error connect db")
	}

	if db != nil {
		fmt.Println("Success Connect")
	}

	r := gin.Default()

	r.GET("/users", func(ctx *gin.Context) {
		var users []User

		result := db.Find(&users)

		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "invalid parameter",
			})
			return
		}
		ctx.JSON(http.StatusOK, users)
	})

	r.POST("/users", func(ctx *gin.Context) {
		var user User

		ctx.ShouldBindJSON(&user)
		db.Create(&user)
		ctx.JSON(http.StatusCreated, user)
	})

	r.GET("/users/:id", func(ctx *gin.Context) {
		var user User

		paramId := ctx.Param("id")
		userId, err := strconv.Atoi(paramId)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "invalid parameter",
			})
			return
		}

		result := db.First(&user, userId)

		if result.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "data not found",
			})

			return
		}

		ctx.JSON(http.StatusOK, user)
	})

	r.PATCH("/users/:id", func(ctx *gin.Context) {
		var user User

		ctx.ShouldBindJSON(&user)

		paramId := ctx.Param("id")
		userId, err := strconv.Atoi(paramId)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "invalid parameter",
			})
			return
		}

		user.ID = userId

		result := db.Save(&user)

		if result.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "data not found",
			})

			return
		}

		ctx.JSON(http.StatusCreated, user)
	})

	r.DELETE("/users/:id", func(ctx *gin.Context) {

		var user User

		paramId := ctx.Param("id")
		userId, err := strconv.Atoi(paramId)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "invalid parameter",
			})
			return
		}

		result := db.Clauses(clause.Returning{}).Where("id = ?", userId).Delete(&user)

		if result.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "data not found",
			})

			return

		}

		ctx.JSON(http.StatusOK, user)
	})

	r.Run()

}
