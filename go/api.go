package main

import (
    "github.com/gin-contrib/cors"
    _ "github.com/gin-contrib/sessions"
    _ "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "net/http"
    _ "encoding/json"
    "fmt"
    "time"
    "math/rand" // 乱数
    "strconv"  // キャスト
    "reflect"  // 型確認
    _ "github.com/jinzhu/gorm"    //v1.0
    _ "github.com/jinzhu/gorm/dialects/mysql" // v1.0
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    _ "gorm.io/gorm/logger"
    _ "errors"
    "strings"
)

type ComModel struct {
    ID uint64 `gorm:"primaryKey"`  // MAX:18446744073709551615
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
    User_id uint64
    // Trn_coin_id int64
    Coin_all uint64
}

type CoinTrn struct {
    ComModel
    Send_user_id uint64
    Receive_user_id uint64
    Coin_id uint64
    Is_coin uint64
}

type DailyReportTrn struct {
    ComModel
    User_id uint64
    Date time.Time `sql:"not null;type:date"`
    Report string
}

// POST 受信用
type Login struct {
    Email string
    Pass string
}

// POST アカウント作成
type Account struct {
    Name string
    Email string
    Pass string
}

// POST 日報
type DailyReport struct {
    Date string
    Report string
}

// POST 日報＋mining
type ReportMining struct {
    Date string
    Mining_coin uint64
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
    mysql_db.AutoMigrate(&User{},&Coin{},&CoinTrn{},&DailyReportTrn{},)

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


    // // 特定ユーザー作成(UserID 1,)
    // var first_user User
    // mysql_db.Where("email = ? and pass = ?", "at.shijo@opt-incubate.com", "12345",).First(&first_user)
    // fmt.Println(first_user)
    // fmt.Println(first_user.Name)
    // if first_user.ID != 1 {
    //     insert_user := User{Name: "shijo", Email: "at.shijo@opt-incubate.com", Pass: "12345"}
    //     result := mysql_db.Create(&insert_user)
    //     fmt.Println(insert_user.ID)
    //     fmt.Println(result.Error)
    //     fmt.Println(result.RowsAffected)
    // }

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"https://localhost","https://dix.front.hello-oi.com",}    // アクセスを許可したいアクセス元
    config.AllowMethods = []string{"GET","POST",}           // アクセスを許可したいHTTPメソッド
    config.AllowCredentials = true                          // cookie情報を必要(true/false)
    config.AllowHeaders = []string{"Content-Type",}         // アクセスを許可したいHTTPリクエストヘッダ
    router.Use(cors.New(config))                            // gin-routerに設定

    // store := cookie.NewStore([]byte("secret"))
    // router.Use(sessions.Sessions("mysession", store))

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

    router.POST("/create_account", func(c *gin.Context){
        fmt.Println("create_account!")

        var user User
        var account Account
        ret := c.Bind(&account)
        fmt.Println(account.Name)
        fmt.Println(account.Email)
        fmt.Println(account.Pass)
        fmt.Println(ret)

        err := mysql_db.Where("email = ?", account.Email,).First(&user).Error
        fmt.Println(user)
        fmt.Println(user.ID)
        fmt.Println(reflect.TypeOf(user.ID))

        if err != nil {
            // アカウント未登録は作成
            insert_user := User{Name: account.Name, Email:  account.Email, Pass:  account.Pass}
            result := mysql_db.Create(&insert_user)
            fmt.Println(insert_user.ID)
            fmt.Println(result.Error)
            fmt.Println(result.RowsAffected)
        
            // Cookieをセット
            // cookie := new(http.Cookie)
            // cookie.Value = strconv.FormatUint(insert_user.ID, 10) //Cookieに入れる値

            // http.SameSiteNoneModeをNoneにしないと、アクセス元ドメインとアクセス先ドメインが違う場合にcookieがはれない
            // c.SetSameSite(http.SameSiteNoneMode)

            // SetCookie(key, value, 保存期間(秒), パス範囲, 利用許可ドメイン, httpsでcookie利用, httpで利用不可)
            // c.SetCookie("user_login", cookie.Value, 604800, "/", "hello-oi.com", true, true)

            c.JSON(200, gin.H{
                "res_flag":true,
                "message":"OK! create user",
            })
        } else {
            // // ローカルの場合
            // if os.Getenv("ENV") == "local" {
            //     log.Println("cookieをセットする")
            //     c.SetCookie("jwt", cookie.Value, 604800, "/", "localhost", true, true)
            // }

            // // 本番環境の場合
            // if os.Getenv("ENV") == "production" {
            //     log.Println("productionでcookieをセットする")
            //     c.SetCookie("jwt", cookie.Value, 604800, "/", "your_domain", true, true)
            // }

            c.JSON(200, gin.H{
                "res_flag":false,
                "message":"allready user",
            })
        }
    })

    router.POST("/login", func(c *gin.Context){
        fmt.Println("from localhost!")

        var user User
        var login Login
        ret := c.Bind(&login)
        fmt.Println(login.Email)
        fmt.Println(login.Pass)
        fmt.Println(ret)

        // session := sessions.Default(c)

        err := mysql_db.Where("email = ? and pass = ?", login.Email, login.Pass,).First(&user).Error
        fmt.Println(user)
        fmt.Println(user.ID)
        fmt.Println(reflect.TypeOf(user.ID))

        if err != nil {
            // session.Set("hello", "world")
            // session.Save()
            c.JSON(200, gin.H{
                "res_flag":false,
                "message":"i don't know user",
                "user":"",
                // "session":session.Get("hello"),
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
            //v := session.Get("count")
            // session.Set("user_login", strconv.FormatUint(user.ID, 10))
            // session.Save()
            // // ローカルの場合
            // if os.Getenv("ENV") == "local" {
            //     log.Println("cookieをセットする")
            //     c.SetCookie("jwt", cookie.Value, 604800, "/", "localhost", true, true)
            // }

            // // 本番環境の場合
            // if os.Getenv("ENV") == "production" {
            //     log.Println("productionでcookieをセットする")
            //     c.SetCookie("jwt", cookie.Value, 604800, "/", "your_domain", true, true)
            // }

            // Cookieをセット
            cookie := new(http.Cookie)
            cookie.Value = strconv.FormatUint(user.ID, 10) //Cookieに入れる値

            // http.SameSiteNoneModeをNoneにしないと、アクセス元ドメインとアクセス先ドメインが違う場合にcookieがはれない
            c.SetSameSite(http.SameSiteNoneMode)

            // SetCookie(key, value, 保存期間(秒), パス範囲, 利用許可ドメイン, httpsでcookie利用, httpで利用不可)
            c.SetCookie("user_login", cookie.Value, 604800, "/", "hello-oi.com", true, true)

            c.JSON(200, gin.H{
                "res_flag":true,
                "message":"hello world",
                "user":user,
            })
        }
    })

    router.POST("/user_wallet", func(c *gin.Context){
        fmt.Println("user_wallet!")

        var coin Coin
        var user_id string
        //var uint64_user_id uint64
        user_id, err = c.Cookie("user_login")
        if err != nil {
            fmt.Println(err)
        }
        // fmt.Printf("%T\n", user_id)
        // fmt.Println(user_id)
        // fmt.Println(c.Cookie("user_login"))

        // strconv.ParseUint(文字列, 基数（10進数）,ビット長)
        uint64_user_id, _ := strconv.ParseUint(user_id, 10, 64)
        fmt.Println("----------------------")
        fmt.Printf("%T\n", uint64_user_id)
        fmt.Println(uint64_user_id)

        if len(user_id) == 0 {
            c.JSON(200, gin.H{
                "res_flag":false,
                "message":"user not found",
            })    
        } else {
            // User情報取得
            var first_user User
            ret_user := mysql_db.Where("id = ?", uint64_user_id,).First(&first_user)
            if ret_user.Error != nil {
                fmt.Println("ユーザー、見つからず")
            }
            fmt.Println("---------------first_user-----------------")
            fmt.Println(first_user)

            // RecordNotFound エラーが返却されたかチェックする,これでもデータなしが判別できる
            // coin_err := mysql_db.First(&coin, 1).Error
            // fmt.Println(errors.Is(coin_err, gorm.ErrRecordNotFound))

            var take_coin uint64
            // coin レコード　なし　作成
            ret_query := mysql_db.Where("user_id = ?", uint64_user_id,).First(&coin)
            fmt.Println("---------------ret_query-----------------")
            fmt.Println(ret_query)
            // sql := mysql_db.ToSQL(func(tx *gorm.DB) *gorm.DB {
            //     return tx.Model(&Coin{}).Where("user_id = ?", uint64_user_id).First(&[]Coin{})
            // })
            if ret_query.Error != nil {
                print("でーたなし、だからレコード作成、初期レコードは１０枚プレゼント")
                insert_coin := Coin{User_id: uint64_user_id, Coin_all: 10}
                result := mysql_db.Create(&insert_coin)
                fmt.Println(insert_coin.ID)
                fmt.Println(result.Error)
                fmt.Println(result.RowsAffected)
                take_coin = insert_coin.Coin_all
            } else {
                take_coin = coin.Coin_all
            }

            fmt.Println("---------------coin-----------------")
            fmt.Println(coin)
            c.JSON(http.StatusOK, gin.H{
                "res_flag": true,
                "message": "wallet",
                "user_name": first_user.Name,
                "have_coin": take_coin,
            })
        }
    })

    router.POST("/wallet_mining", func(c *gin.Context){
        fmt.Println("wallet_mining!")

        // 乱数初期化
        rand.Seed(time.Now().UnixNano())

        // 乱数取得(1-100までの範囲)
        rand_max := 100
        mining_coin := rand.Intn(rand_max)

        var coin Coin
        var user_id string
        //var uint64_user_id uint64
        user_id, err = c.Cookie("user_login")
        if err != nil {
            fmt.Println(err)
        }
        // fmt.Printf("%T\n", user_id)
        // fmt.Println(user_id)
        // fmt.Println(c.Cookie("user_login"))

        // strconv.ParseUint(文字列, 基数（10進数）,ビット長)
        uint64_user_id, _ := strconv.ParseUint(user_id, 10, 64)
        fmt.Println("----------------------")
        fmt.Printf("%T\n", uint64_user_id)
        fmt.Println(uint64_user_id)

        if len(user_id) == 0 {
            c.JSON(200, gin.H{
                "res_flag":false,
                "message":"user not found",
            })    
        } else {
            // coin レコード　確認
            ret_query := mysql_db.Where("user_id = ?", uint64_user_id,).First(&coin)
            if ret_query.Error != nil {
                print("でーたなし、だから異常")
            } else {
                // coinレコードをUPDATE、ただし、mining_coinが30以下の場合は更新しない
                fmt.Println("coin.Coin_all")
                fmt.Println(coin.Coin_all)
                var ret_message string = "採掘できませんでした。"
                if mining_coin >= 30 {
                    mysql_db.Model(&coin).Select("Coin_all").Updates(Coin{Coin_all: (coin.Coin_all + uint64(mining_coin))})
                    ret_message = strconv.FormatUint(uint64(mining_coin), 10) + "採掘できました。"
                    fmt.Println("coin.Coin_all")
                    fmt.Println(coin.Coin_all)
                }
                c.JSON(http.StatusOK, gin.H{
                    "res_flag": true,
                    "message": ret_message,
                    "have_coin": coin.Coin_all,
                })
            }
        }
    })

    router.POST("/create_daily_report", func(c *gin.Context){
        fmt.Println("create_daily_report!")

        // 乱数初期化
        rand.Seed(time.Now().UnixNano())

        // 乱数取得(1-100までの範囲)
        rand_max := 100
        mining_coin := rand.Intn(rand_max)

        var date_report DailyReport
        c.Bind(&date_report)
        // input_date := c.PostForm("date")
        // input_report := c.PostForm("report")
        fmt.Println(date_report)

        var user_id string
        //var uint64_user_id uint64
        user_id, err = c.Cookie("user_login")
        if err != nil {
            fmt.Println(err)
            c.JSON(200, gin.H{
                "res_flag":false,
                "message":"user not find",
            })
        } else {
            // strconv.ParseUint(文字列, 基数（10進数）,ビット長)
            uint64_user_id, _ := strconv.ParseUint(user_id, 10, 64)
            fmt.Println("----------------------")
            fmt.Printf("%T\n", uint64_user_id)
            fmt.Println(uint64_user_id)

            strTime := strings.Split(date_report.Date, "-")
            yyyy, _ := strconv.Atoi(strTime[0])
            mm, _ := strconv.Atoi(strTime[1])
            dd, _ := strconv.Atoi(strTime[2])
            fmt.Println(yyyy)
            fmt.Println(mm)
            fmt.Println(dd)
            t := time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.Local)

            insert_daily_report := DailyReportTrn{User_id: uint64_user_id, Date: t, Report: date_report.Report}
            result := mysql_db.Create(&insert_daily_report)
            fmt.Println(result.Error)
            fmt.Println(result.RowsAffected)

            c.JSON(200, gin.H{
                "res_flag":true,
                "message":"OK! create daily_report",
                "mining_coin":mining_coin,
            })
        }
    })

    router.POST("/get_dix_coin_report", func(c *gin.Context){
        fmt.Println("get_dix_coin_report!")

        // // 乱数初期化
        // rand.Seed(time.Now().UnixNano())

        // // 乱数取得(1-100までの範囲)
        // rand_max := 100
        // mining_coin := rand.Intn(rand_max)

        var report_mining ReportMining
        ret := c.Bind(&report_mining)
        fmt.Println("===============================")
        fmt.Println(report_mining)
        fmt.Println("===============================")
        fmt.Println(ret)

        var coin Coin
        var daily_report DailyReportTrn
        var user_id string
        //var uint64_user_id uint64
        user_id, err = c.Cookie("user_login")
        if err != nil {
            fmt.Println(err)
        }

        // strconv.ParseUint(文字列, 基数（10進数）,ビット長)
        uint64_user_id, _ := strconv.ParseUint(user_id, 10, 64)
        fmt.Println("----------------------")
        fmt.Printf("%T\n", uint64_user_id)
        fmt.Println(uint64_user_id)

        if len(user_id) == 0 {
            c.JSON(200, gin.H{
                "res_flag":false,
                "message":"user not found",
            })    
        } else {
            // coin レコード　確認
            ret_query := mysql_db.Where("user_id = ?", uint64_user_id,).First(&coin)
            if ret_query.Error != nil {
                c.JSON(http.StatusOK, gin.H{
                    "res_flag": true,
                    "message": "coins レコード なし 異常",
                    "have_coin": coin.Coin_all,
                })
                return
            }
            
            // daily_report レコード　確認　対象日のレコードが存在していたらdix_coinを付与しない
            strTime := strings.Split(report_mining.Date, "-")
            yyyy, _ := strconv.Atoi(strTime[0])
            mm, _ := strconv.Atoi(strTime[1])
            dd, _ := strconv.Atoi(strTime[2])
            fmt.Println(yyyy)
            fmt.Println(mm)
            fmt.Println(dd)
            t := time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.Local)
            ret_query2 := mysql_db.Where("user_id = ? and date = ?", uint64_user_id, t,).First(&daily_report)
            if ret_query2.Error != nil {
                c.JSON(http.StatusOK, gin.H{
                    "res_flag": true,
                    "message": "既に登録済",
                    "have_coin": coin.Coin_all,
                })
                return
            }
            
            // coinレコードをUPDATE
            fmt.Println("report_mining.Mining_coin")
            fmt.Println(report_mining.Mining_coin)
            var update_coin uint64 = coin.Coin_all + uint64(report_mining.Mining_coin)
            fmt.Println("coin.Coin_all")
            fmt.Println(coin.Coin_all)
            
            mysql_db.Model(&coin).Select("Coin_all").Updates(Coin{Coin_all: update_coin})
            ret_message := strconv.FormatUint(uint64(update_coin), 10) + "所有しています"
            fmt.Println("coin.Coin_all")
            fmt.Println(coin.Coin_all)

            c.JSON(http.StatusOK, gin.H{
                "res_flag": true,
                "message": ret_message,
                "have_coin": coin.Coin_all,
            })
            return
        }
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