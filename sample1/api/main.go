package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"../api/handler"
	"../api/middleware"
	"../model/admin"
	"../model/user"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
	}
	os.Setenv("secret", viper.GetString("jwt_secret"))
}

func dbConnect(host, port, user, dbname, password, sslmode string) (*gorm.DB, error) {

	// In the case of heroku
	if os.Getenv("DATABASE_URL") != "" {
		return gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	}
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode))

	return db, err
}

func main() {
	mode := viper.GetString("mode")

	// DB binding
	dbprefix := "database_" + mode
	dbhost := viper.GetString(dbprefix + ".host")
	dbport := viper.GetString(dbprefix + ".port")
	dbuser := viper.GetString(dbprefix + ".user")
	dbname := viper.GetString(dbprefix + ".dbname")
	dbpassword := viper.GetString(dbprefix + ".password")
	dbsslmode := viper.GetString(dbprefix + ".sslmode")

	db, err := dbConnect(dbhost, dbport, dbuser, dbname, dbpassword, dbsslmode)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err.Error())
	}
	defer db.Close()
	log.Println("Connected to the database")

	// migrations
	db.AutoMigrate(&user.User{}, &admin.Admin{})

	// todo: 源码的艺术，  建立db， repo对接db， handler对接repo
	// initializing repos and services
	// fixme: 我觉得这里user 和 admin共用一个repo就可以了， 没必要分开两个声明
	// fixme: 但是假如要用一个repo的话， 那么就必须在model文件夹的首层目录就定义repo的声明源码，不能再子文件夹定义
	userRepo := user.NewPostgresRepo(db)
	adminRepo := admin.NewPostgresRepo(db)

	userSvc := user.NewService(userRepo)
	adminSvc := admin.NewService(adminRepo)

	// Initializing handlers
	r := http.NewServeMux()

	// todo: repo数据层分两层， entity+repo+实际数据库操作 为一层， service问一层，
	// todo: 最终的效果都是传内存声明到handler去, userSvc就是这个内存声明
	handler.MakeUserHandler(r, userSvc)
	handler.MakeAdminHandler(r, adminSvc)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})

	// HTTP(s) binding
	serverprefix := "server_" + mode
	host := viper.GetString(serverprefix + ".host")
	port := os.Getenv("PORT")
	timeout := time.Duration(viper.GetInt("timeout"))

	if port == "" {
		port = viper.GetString(serverprefix + ".port")
	}

	conn := host + ":" + port

	// middlewares
	mwCors := middleware.CorsEveryWhere(r)										// 允许跨域
	mwLogs := handlers.LoggingHandler(os.Stdout, mwCors)			// 输出access日志

	srv := &http.Server{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Addr:         conn,
		Handler:      mwLogs,
	}

	log.Printf("Starting in %s mode", mode)
	log.Printf("Server running on %s", conn)
	log.Fatal(srv.ListenAndServe())
}
