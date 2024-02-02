package store

import (
	"context"
	"path/filepath"

	"github.com/yangchnet/pm/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteStore struct {
	db *gorm.DB
}

func NewSqliteStore(ctx context.Context) *SqliteStore {
	db, err := gorm.Open(sqlite.Open(filepath.Join(config.GetString("local.path"), "passwd.db")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&Passwd{})

	return &SqliteStore{
		db: db,
	}
}

// Save 在使用cryptFunc对密码密文进行存储
func (s *SqliteStore) Save(ctx context.Context, name, url, username string, cryptedPasswd []byte, note string) error {
	return s.db.Save(&Passwd{
		Name:          name,
		Url:           url,
		UserName:      username,
		Note:          note,
		CryptedPasswd: cryptedPasswd,
	}).Error
}

// Get 获取密码密文
func (s *SqliteStore) Get(ctx context.Context, name string) (*Passwd, error) {
	var passwd *Passwd
	if err := s.db.Model(&Passwd{}).Where("name = ?", name).First(&passwd).Error; err != nil {
		return passwd, nil
	}

	return passwd, nil
}

// SearchName 根据名称进行搜索并给出名称列表
func (s *SqliteStore) SearchName(ctx context.Context, name string) ([]string, error) {
	var names []string
	if err := s.db.Model(&Passwd{}).Where("name LIKE ?", "%"+name+"%").Select("name").Scan(&names).Error; err != nil {
		return nil, err
	}
	return names, nil
}

// Delete 删除一个记录
func (s *SqliteStore) Delete(ctx context.Context, name string) error {
	var passwd Passwd
	if err := s.db.Where("name = ?", name).First(&passwd).Error; err != nil {
		return err
	}
	return s.db.Delete(&passwd).Error
}
