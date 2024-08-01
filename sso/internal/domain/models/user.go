package models

type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement:true;uniqueIndex"`
	Username string `gorm:"unique;not null;type:varchar(64);uniqueIndex"`
	PassHash []byte `gorm:"not null;type:bytea"`
}

type UserAPI struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Admin struct {
	ID      int64 `gorm:"primaryKey;autoIncrement:true;uniqueIndex"`
	IsAdmin bool  `gorm:"default:false"`
}
