package models

type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement:true;uniqueIndex"`
	Username string `gorm:"unique;not null;type:varchar(64);uniqueIndex"`
	PassHash []byte `gorm:"not null;type:bytea"`
}

type Admin struct {
	ID      int64 `gorm:"primaryKey;autoIncrement:true;uniqueIndex"`
	IsAdmin bool  `gorm:"default:false"`
}

//func UserToApiUser(user *User) *api.User{
//	return &api.User{
//		Username: user.Username,
//		Password: user.,
//	}
//}
