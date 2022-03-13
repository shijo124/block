package main

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    _ "net/http"
    _ "encoding/json"
    "fmt"
    "time"
    _ "github.com/jinzhu/gorm"    //v1.0
    _ "github.com/jinzhu/gorm/dialects/mysql" // v1.0
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    _ "gorm.io/gorm/logger"
)

type ComModel struct {
    ID int64 `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
    ComModel
    Name string
    Email string
    Pass string
}

type Coin struct {
    ComModel
    User_id int64
    Trn_coin_id int64
    Coin_all int64
}

type CoinTrn struct {
    ComModel
    Send_user_id int64
    Receive_user_id int64
    coin int64
}

// POST 受信用
type Login struct {
    Email string
    Pass string
}

func main(){
    // make router
    router := gin.Default()

    fmt.Println("DB接続")
    //db, err := gorm.Open("mysql", "dix_user:@Shijo1603@/dix_coin?charset=utf8mb4&parseTime=True&loc=Local")
    dsn := "dix_user:@Shijo1603@tcp(127.0.0.1:3306)/dix_coin?charset=utf8mb4&parseTime=True&loc=Local"
    mysql_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("データベースへの接続に失敗しました")
    }
    mysql_db.AutoMigrate(&User{},&Coin{},&CoinTrn{},)

    db, err := mysql_db.DB()    // gorm v2.0 仕様(バグ？) で、DB()を取得して、deferでClose
    defer db.Close()    // Close

    // db.LogMode(true)
    fmt.Println(mysql_db)
    fmt.Println("DB接続完了")

    /*
    router.Use(cors.New(cors.Config{
        // access ok prigin
        AllowOrigins: []string{
            "http://localhost",
        },
        // access ok http method
        AllowMethods: []string{
            "POST",
            "GET",
        },
        // header
        AllowHeaders: []string{
            "Content-Type",
        },
        AllowCredentials: true,
        MaxAge: 24 * time.Hour,
    }))
    */


    // 特定ユーザー作成(UserID 1,)
    var first_user User
    mysql_db.Where("email = ? and pass = ?", "at.shijo@opt-incubate.com", "12345",).First(&first_user)
    fmt.Println(first_user)
    fmt.Println(first_user.Name)
    if first_user.ID != 1 {
        insert_user := User{Name: "shijo", Email: "at.shijo@opt-incubate.com", Pass: "12345"}
        result := mysql_db.Create(&insert_user)
        fmt.Println(insert_user.ID)
        fmt.Println(result.Error)
        fmt.Println(result.RowsAffected)
    }


    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"https://localhost",}    // アクセスを許可したいアクセス元
    config.AllowMethods = []string{"GET","POST",}            // アクセスを許可したいHTTPメソッド
    config.AllowCredentials = true                            // cookie情報を必要(true/false)
    config.AllowHeaders = []string{"Content-Type",}            // アクセスを許可したいHTTPリクエストヘッダ

    router.Use(cors.New(config))                            // gin-routerに設定

    // request is GET return hello world
    router.GET("/", func(c *gin.Context){
        fmt.Println("from localhost!")
        // c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        // c.Writer.Header().Set("Access-Control-Max-Age", "86400")
        // c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
        // c.Writer.Header().Set("Content-Type", "application/json")
        // c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        // c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
        // c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        // c.Writer.Header().Set("Content-Type", "application/json")
        c.JSON(200, gin.H{
            "res_flag":true,
            "message":"hello world",
        })
    })

    router.POST("/login", func(c *gin.Context){
        fmt.Println("from localhost!")

        var user User
        var login Login
        ret := c.Bind(&login)
        fmt.Println(login.Email)
        fmt.Println(login.Pass)
        fmt.Println(ret)

        err := mysql_db.Where("email = ? and pass = ?", login.Email, login.Pass,).First(&user).Error
        fmt.Println(user)
        if err != nil {
            c.JSON(200, gin.H{
                "res_flag":false,
                "message":"i don't know user",
                "user":"",
            })
        } else {
            // c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
            // c.Writer.Header().Set("Access-Control-Max-Age", "86400")
            // c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
            // c.Writer.Header().Set("Content-Type", "application/json")
            // c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
            // c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
            // c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
            // c.Writer.Header().Set("Content-Type", "application/json")
            c.JSON(200, gin.H{
                "res_flag":true,
                "message":"hello world",
                "user":user,
            })
        }
    })

    router.GET("/user_coin", func(c *gin.Context){
        fmt.Println("user_coin!")

        var user User
        var login Login
        ret := c.Bind(&login)
        fmt.Println(login)
        fmt.Println(ret)

        mysql_db.Where("email = ? and pass = ?", "at.shijo@opt-incubate.com", "12345",).First(&user)
        fmt.Println(user)
        // c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        // c.Writer.Header().Set("Access-Control-Max-Age", "86400")
        // c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
        // c.Writer.Header().Set("Content-Type", "application/json")
        // c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        // c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
        // c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        // c.Writer.Header().Set("Content-Type", "application/json")
        c.JSON(200, gin.H{
            "res_flag":true,
            "message":"hello world",
            "user":user,
        })
    })

    // run server
    router.Run(":9000")

    //http.HandleFunc("/", handler)
    //http.ListenAndServe(":9000", nil)
}

/*
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World!!!")
}
*/
