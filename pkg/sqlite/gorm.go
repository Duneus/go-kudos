package sqlite

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func NewGorm(filePath string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", filePath)

	if err != nil {
		return nil, fmt.Errorf("sqlite storage: cannot open database: err=%v", err)
	}

	if err = db.DB().Ping(); err != nil {
		return nil, fmt.Errorf("sqlite  storage: cannot connect to database: err=%v", err)
	}

	fmt.Println("Successfully connected to database")

	return db, nil
}
