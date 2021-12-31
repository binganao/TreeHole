package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/http"
)

type User struct {
	EncUserName string `gorm:"type:text;not null"`
	EncPassWord string `gorm:"type:text;not null"`
	EncSecret string `gorm:"type:text;not null"`
	gorm.Model
}

type Datas struct {
	Msg string
}

type Resp struct {
	Code string
	User []User
	Data Datas
}

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "treehole"
	databaseUser := "root"
	databasePwd := "root"
	charset := "utf8mb4"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		databaseUser,
		databasePwd,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	return db
}

func main()  {
	db := InitDB()
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	r := gin.Default()

	r.POST("/api/login", func(ctx *gin.Context) {
		encUserName, _ := ctx.GetPostForm("encUserName")
		encPassWord, _ := ctx.GetPostForm("encPassWord")

		if verifyLogin(encUserName, encPassWord, db) {
			ctx.JSON(http.StatusOK, Resp{Code: "200", Data: Datas{Msg: "success"}})
		} else {
			ctx.JSON(http.StatusInternalServerError, Resp{Code: "500", Data: Datas{Msg: "fail"}})
		}
	})

	r.POST("/api/register", func(ctx *gin.Context) {
		encUserName, _ := ctx.GetPostForm("encUserName")
		encPassWord, _ := ctx.GetPostForm("encPassWord")

		if doRegister(encUserName, encPassWord, db) {
			ctx.JSON(http.StatusOK, Resp{Code: "200", Data: Datas{Msg: "success"}})
		} else {
			ctx.JSON(http.StatusInternalServerError, Resp{Code: "500", Data: Datas{Msg: "fail"}})
		}

	})

	r.POST("/api/secretlist", func(ctx *gin.Context) {
		encUserName, _ := ctx.GetPostForm("encUserName")
		encPassWord, _ := ctx.GetPostForm("encPassWord")

		if b, data := getEncSecretList(encUserName, encPassWord, db); b {
			ctx.JSON(http.StatusOK, Resp{Code: "200", Data: Datas{Msg: "success"}, User: data})
		} else {
			ctx.JSON(http.StatusInternalServerError, Resp{Code: "500", Data: Datas{Msg: "fail"}})
		}
	})

	r.POST("/api/addsecret", func(ctx *gin.Context) {
		encUserName, _ := ctx.GetPostForm("encUserName")
		encPassWord, _ := ctx.GetPostForm("encPassWord")
		encSecret, _ := ctx.GetPostForm("encSecret")

		if addSecret(encUserName, encPassWord, encSecret, db) {
			ctx.JSON(http.StatusOK, Resp{Code: "200", Data: Datas{Msg: "success"}})
		} else {
			ctx.JSON(http.StatusInternalServerError, Resp{Code: "500", Data: Datas{Msg: "fail"}})
		}
	})

	panic(r.Run("0.0.0.0:1234"))
}

func verifyLogin(encUserName, encPassWord string, db *gorm.DB) bool {
	var user User

	db.Where("enc_user_name = ? and enc_pass_word = ?", encUserName, encPassWord).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}

func doRegister(encUserName, encPassWord string, db *gorm.DB) bool {
	fmt.Println(encUserName)

	user := User{
		EncUserName: encUserName,
		EncPassWord: encPassWord,
		EncSecret: "",
	}

	fmt.Println(user)

	if err := db.Create(&user).Error; err == nil{
		return true
	}

	return false
}

func getEncSecretList(encUserName, encPassWord string, db *gorm.DB) (bool, []User) {
	var data []User

	db.Where("enc_user_name = ? and enc_pass_word = ?", encUserName, encPassWord).Not("enc_secret = ?", "").Find(&data)
	if len(data) >= 0 {
		return true, data
	}


	return false, []User{}
}

func addSecret(encUserName, encPassWord, encSecret string, db *gorm.DB) bool {
	user := User{
		EncUserName: encUserName,
		EncPassWord: encPassWord,
		EncSecret: encSecret,
	}

	if err := db.Create(&user).Error; err == nil{
		return true
	}

	return false
}