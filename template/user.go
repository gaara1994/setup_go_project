package models

import "github.com/jinzhu/gorm"

// User represents the User model with comments for each field.
type User struct {
	gorm.Model
	Username     string `gorm:"column:user_name;not null;uniqueIndex;comment:用户名"` // 用户名，设置为唯一索引
	PasswordHash string `gorm:"column:password_hash;not null;comment:密码散列值"`       // 密码散列值
	Email        string `gorm:"column:email;comment:邮箱"`                           // 邮箱，设置为唯一索引
}

// BeforeSave is a GORM callback, used for operations like hashing the password before saving.
func (u *User) BeforeSave(_ *gorm.DB) error {
	// 实际应用中在此处使用安全的密码哈希方法
	// u.Password = hashedPassword(u.Password)
	return nil
}
