package gormex

import (
	_ "embed"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
	"log"
	"net/http"
	"testing"
	"time"
)

//go:embed mysql.yaml
var mysqlDSN string

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100"`
	Age  int
}

func TestExample(t *testing.T) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		fmt.Println("running: 9191")
		if err := http.ListenAndServe(":9191", nil); err != nil {
			log.Fatalf("failed to start metrics server: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)

	db := InitDB(mysqlDSN)

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	for {
		go func() {
			createUser(db, "a", 30)
			createUser(db, "b", 25)

			queryUser(db, 1)
			queryUser(db, 2)

			updateUser(db, 1, "aa", 31)

			deleteUser(db, 2)

			listAllUsers(db)
		}()
		time.Sleep(100 * time.Millisecond)
	}
}

func createUser(db *gorm.DB, name string, age int) {
	newUser := User{Name: name, Age: age}
	if err := db.Create(&newUser).Error; err != nil {
		log.Fatalf("failed to create sso: %v", err)
	}
	fmt.Printf("Created sso: %+v\n", newUser)
}

func queryUser(db *gorm.DB, id uint) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		log.Printf("failed to find sso with UID %d: %v\n", id, err)
		return
	}
	fmt.Printf("Found sso: %+v\n", user)
}

func updateUser(db *gorm.DB, id uint, newName string, newAge int) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		log.Printf("failed to find sso with UID %d for update: %v\n", id, err)
		return
	}

	user.Name = newName
	user.Age = newAge
	if err := db.Save(&user).Error; err != nil {
		log.Printf("failed to update sso: %v\n", err)
		return
	}
	fmt.Printf("Updated sso: %+v\n", user)
}

func deleteUser(db *gorm.DB, id uint) {
	if err := db.Delete(&User{}, id).Error; err != nil {
		log.Printf("failed to delete sso with UID %d: %v\n", id, err)
		return
	}
	fmt.Printf("Deleted sso with UID %d\n", id)
}

func listAllUsers(db *gorm.DB) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("failed to list users: %v\n", err)
		return
	}
	fmt.Println("All users:")
	for _, user := range users {
		fmt.Printf("- %+v\n", user)
	}
}
