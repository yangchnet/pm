package filestore

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/store"
)

type FileStore struct {
	localPath string
}

func NewFileStore(ctx context.Context) *FileStore {
	return &FileStore{
		localPath: config.GetString("local.path"),
	}
}

var _ store.Store = &FileStore{}

func (s *FileStore) Init(ctx context.Context) (string, error) {
	return `store:
  type: file`, nil
}

// Save 在使用cryptFunc对密码密文进行存储
func (s *FileStore) Save(ctx context.Context, passwd *store.Passwd) error {
	files, err := readAllPasswd(s.localPath)
	if err != nil {
		return err
	}

	_, ok := files[passwd.Name+".passwd"]
	if ok {
		return store.ErrAlreadyExists
	}

	f, err := os.Create(filepath.Join(s.localPath, passwd.Name+".passwd"))
	if err != nil {
		return err
	}
	defer f.Close()

	passwdByte, err := json.Marshal(passwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	encoded := base64.StdEncoding.EncodeToString(passwdByte)
	_, err = f.WriteString(encoded)
	if err != nil {
		return err
	}

	return nil
}

// Get 获取密码
func (s *FileStore) Get(ctx context.Context, name string) (*store.Passwd, error) {
	files, err := readAllPasswd(s.localPath)
	if err != nil {
		return nil, err
	}

	path, ok := files[name+".passwd"]
	if !ok {
		return nil, store.ErrNotFound
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(f))
	if err != nil {
		return nil, err
	}

	var passwd store.Passwd
	if err := json.Unmarshal(decoded, &passwd); err != nil {
		return nil, err
	}

	return &passwd, nil
}

// SearchName 根据名称进行搜索并给出名称列表
func (s *FileStore) SearchName(ctx context.Context, name string) ([]string, error) {
	files, err := readAllPasswd(s.localPath)
	if err != nil {
		return nil, err
	}

	var names []string
	for k, _ := range files {
		list := strings.Split(k, ".")
		if len(list) < 1 {
			continue
		}

		if strings.Contains(strings.ToLower(list[0]), strings.ToLower(name)) {
			names = append(names, list[0])
		}
	}

	return names, nil
}

// Delete 删除一个记录
func (s *FileStore) Delete(ctx context.Context, name string) error {
	return os.Remove(filepath.Join(s.localPath, name+".passwd"))
}

// Update 更新一个记录
func (s *FileStore) Update(ctx context.Context, name string, passwd *store.Passwd) error {
	files, err := readAllPasswd(s.localPath)
	if err != nil {
		return err
	}

	path, ok := files[name+".passwd"]
	if !ok {
		return store.ErrNotFound
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	passwdByte, err := json.Marshal(passwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	encoded := base64.StdEncoding.EncodeToString(passwdByte)
	_, err = f.WriteString(encoded)
	if err != nil {
		return err
	}

	return nil
}

func readAllPasswd(dir string) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".passwd" {
			files[info.Name()] = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
