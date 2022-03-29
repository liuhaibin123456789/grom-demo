package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//type User struct {
//	//gorm.Model `json:"gorm_._model"` //gorm会自动将匿名嵌套结构体嵌入父级结构体字段。也可以使用gorm的tag：embedded
//	ID           uint   `json:"id" gorm:"column:id;type:int;primaryKey;autoIncrement;comment:'id主键'"`
//	Name         string `json:"name" gorm:"column:name;type:varchar(20);not null;comment:'用户名'"`
//	Sex          string `json:"sex" gorm:"column:sex;type:varchar(5);not null;comment:'性别'"`
//	Age          uint8  `json:"age" gorm:"column:age;type:int;not null;comment:'年龄'"`
//	Email        string `json:"email" gorm:"column:email;type:varchar(20);not null;unique;comment:'邮箱'"`
//	Phone        string `json:"phone" gorm:"column:phone;type:varchar(11);not null;unique;comment:'手机号'"`
//	QQ           string `json:"qq" gorm:"column:qq;type:varchar(15);not null;unique;comment:'qq号'"`
//	Introduction string `json:"introduction" gorm:"column:introduction;type:varchar(1000);comment:'自我介绍'"`
//}

type User struct {
	//gorm.Model `json:"gorm_._model"` //gorm会自动将匿名嵌套结构体嵌入父级结构体字段。也可以使用gorm的tag：embedded
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:'id主键'"`
	Name         string `json:"name" gorm:"type:varchar(20);comment:'用户名'"`
	Sex          string `json:"sex" gorm:"type:varchar(5);comment:'性别'"`
	Age          uint8  `json:"age" gorm:"comment:'年龄'"`
	Email        string `json:"email" gorm:"type:varchar(20);comment:'邮箱'"`
	Phone        string `json:"phone" gorm:"type:varchar(11);comment:'手机号'"`
	QQ           string `json:"qq" gorm:"type:varchar(15);comment:'qq号'"`
	Introduction string `json:"introduction" gorm:"type:varchar(1000);comment:'自我介绍'"`
}

// TableName 实现接口，自定义表名
func (User) TableName() string {
	return "user"
}

func main() {
	//连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_db?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		//	一般手动禁用事务
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建表
	err = db.AutoMigrate(&User{})
	if err != nil {
		fmt.Println()
		return
	}
	create(db)
	read(db)
	update(db)
	transaction(db)
}

func create(db *gorm.DB) {
	//创建单个学生信息
	user := User{Name: "cold bin", Sex: "man", Age: 20}
	tx := db.Create(&user)
	fmt.Println(user.ID, tx.RowsAffected, tx.Error)
	//创建多个学生信息
	users := []User{{Name: "bwj", Age: 20}, {Name: "wtl", Age: 21}, {Name: "pl", Age: 20}}
	tx = db.Create(&users)
	fmt.Println(tx.Error, users[2].ID)
}

func read(db *gorm.DB) {
	// SELECT * FROM users ORDER BY id LIMIT 1;
	user := User{}
	db.First(&user)
	fmt.Println(user)
	// SELECT * FROM users LIMIT 1;
	user = User{}
	db.Take(&user)
	fmt.Println(user)
	// SELECT * FROM users ORDER BY id DESC LIMIT 1;
	user = User{}
	db.Last(&user)
	fmt.Println(user)
}

func update(db *gorm.DB) {
	//保存更新所有字段 save
	user := User{}
	db.First(&user)
	user.QQ = "3398341353"
	user.Phone = "15736469310"
	db.Save(&user)
	//更新单个列 update
	err := db.Model(&User{}).Where("id = ?", "2").Update("phone", "10102").Error
	fmt.Println(err)
	//更新多列 updates
	err = db.Model(&User{}).Where("id=?", 3).Updates(map[string]interface{}{"qq": "123", "phone": "1111"}).Error
	fmt.Println(err)
}

func transaction(db *gorm.DB) {
	//手动开启事务
	tx := db.Begin()
	//遇到panic回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&User{Name: "Giraffe"}).Error; err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Create(&User{Name: "Lion"}).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
}
