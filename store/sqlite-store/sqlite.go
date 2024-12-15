package sqlitestore

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type SqliteStore struct {
	db *gorm.DB
}

func NewSqliteStore(ctx context.Context) *SqliteStore {
	db, err := gorm.Open(sqlite.Open(filepath.Join(config.GetString("local.path"), "passwd.db")), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&store.Passwd{})

	return &SqliteStore{
		db: db,
	}
}

func (s *SqliteStore) Init(ctx context.Context) (string, error) {
	return `store:
  type: sqlite`, nil
}

// Save 在使用cryptFunc对密码密文进行存储
func (s *SqliteStore) Save(ctx context.Context, passwd *store.Passwd) error {
	return s.db.Save(passwd).Error
}

// Get 获取密码密文
func (s *SqliteStore) Get(ctx context.Context, name string) (*store.Passwd, error) {
	var passwd *store.Passwd
	if err := s.db.Model(&store.Passwd{}).Where("LOWER(name) = ?", strings.ToLower(name)).First(&passwd).Error; err != nil {
		return nil, err
	}

	return passwd, nil
}

// SearchName 根据名称进行搜索并给出名称列表
func (s *SqliteStore) SearchName(ctx context.Context, name string) ([]*store.Passwd, error) {
	var passwds []*store.Passwd
	if err := s.db.Model(&store.Passwd{}).Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%").Find(&passwds).Error; err != nil {
		return nil, err
	}
	return passwds, nil
}

// Delete 删除一个记录
func (s *SqliteStore) Delete(ctx context.Context, name string) error {
	var passwd store.Passwd
	if err := s.db.Where("name = ?", name).First(&passwd).Error; err != nil {
		return err
	}
	return s.db.Delete(&passwd).Error
}

// Update 更新一个记录
func (s *SqliteStore) Update(ctx context.Context, name string, passwd *store.Passwd) error {
	return s.db.Model(&store.Passwd{}).Where("name = ?", name).Updates(passwd).Error
}
