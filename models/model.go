// model/model.go
package models

type Docs struct {
    ID uint `gorm:"primaryKey"`
    Item string `gorm:"unique"`
    Completed int `gorm:"unique"`
}

