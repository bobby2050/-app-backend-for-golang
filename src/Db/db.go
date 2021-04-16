package Db
//package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)
var DBHelper *gorm.DB
var err error


func init() {

	// 打印日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,         // Disable color
		},
	)

	dsn := "meizi:123456@tcp(127.0.0.1:8889)/meizi?charset=utf8mb4&parseTime=True&loc=Local"
	DBHelper, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "tbl_",   // table name prefix, table for `User` would be `t_users`
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})

	if err != nil {
		log.Fatal("连接数据库异常-1")
	}
	sqlDB, err := DBHelper.DB()
	if err != nil {
		log.Fatal("连接数据库异常-2")
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(50)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	//sqlDB.SetConnMaxLifetime(time.Second * 60)


	//var tc []Model.Answer
	////db.Where("typeId=?", "abc").Find(&tc)
	//db.Where(&Model.Answer{TypeId:"abc"}).Find(&tc)
	//fmt.Println(tc)



	//tc := &Model.Answer{}
	//db.First(&tc, 10)
	//fmt.Println(tc)




	//rows,_ := db.Raw("select id, option_name from answer limit 10").Rows()
	//
	//for rows.Next() {
	//	var id int
	//	var option_name string
	//	rows.Scan(&id, &option_name)
	//
	//	fmt.Println(id, option_name)
	//}



}