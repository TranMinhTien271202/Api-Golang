package main

import (
	"database/sql"
	"log"
	"net/http"

	"Test/controllers" // Import controllers
	"Test/routes"      // Sử dụng module routes

	_ "github.com/go-sql-driver/mysql" // Driver MySQL
)

var db *sql.DB

func initDB() {
	var err error
	// Chuỗi kết nối MySQL
	dsn := "tien:@tcp(localhost:3306)/golang"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Không thể kết nối đến cơ sở dữ liệu: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Cơ sở dữ liệu không khả dụng: ", err)
	}
	log.Println("Kết nối cơ sở dữ liệu thành công")
}

func main() {
	initDB()
	defer db.Close() // Đóng kết nối khi không còn sử dụng

	// Khởi tạo controller với db
	controllers.InitController(db)

	routes.SetupRoutes() // Định nghĩa routes

	// Bắt đầu server trên cổng 9000
	log.Println("Server đang chạy tại http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
