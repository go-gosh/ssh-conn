package db

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type HostTag struct {
	ID   uint   `gorm:"primaryKey"`
	Host string `gorm:"index;not null"`
	Tag  string `gorm:"index;not null"`
}

func getDBPath() string {
	home, _ := os.UserHomeDir()
	sshDir := filepath.Join(home, ".ssh")
	os.MkdirAll(sshDir, 0700)
	return filepath.Join(sshDir, "sshconfig.db")
}

func OpenDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(getDBPath()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&HostTag{})
	return db, nil
}

// 添加多个标签
func AddTags(host string, tags []string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		db.Create(&HostTag{Host: host, Tag: tag})
	}
	return nil
}

// 移除多个标签
func RemoveTags(host string, tags []string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		db.Where("host = ? AND tag = ?", host, tag).Delete(&HostTag{})
	}
	return nil
}

// 设置主机标签（覆盖）
func SetTags(host string, tags []string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	db.Where("host = ?", host).Delete(&HostTag{})
	return AddTags(host, tags)
}

// 查询主机所有标签
func GetTags(host string) ([]string, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	var tags []HostTag
	err = db.Where("host = ?", host).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	var result []string
	for _, t := range tags {
		result = append(result, t.Tag)
	}
	return result, nil
}
