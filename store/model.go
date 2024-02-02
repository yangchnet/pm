package store

import "time"

type Passwd struct {
	Id            int
	Name          string `gorm:"index:idx_name,unique"` // Name is unique，供用户查找密码
	Url           string // 用户该密码对应的的网址
	UserName      string // 用户在该网站上的用户名
	CryptedPasswd []byte // 用户密码，加密后的
	Note          string // 备注
	CreateTime    time.Time
	UpdateTime    time.Time
}
