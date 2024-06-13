package model

import "gorm.io/gorm"

type Employee struct {
	EmployeeId string  `gorm:"primary key;auto increment" json:"employee_id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Salary     float32 `json:"salary"`
	Department string  `json:"department"`
}

func MigrateEmployee(db *gorm.DB) error {
	err := db.AutoMigrate(&Employee{})
	return err
}
