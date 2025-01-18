package models

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImageModel struct {
	Model
	Filename string `gorm:"size:64" json:"filename"`
	Path     string `gorm:"size:256" json:"path"`
	Size     int64  `json:"size"`
	Hash     string `gorm:"size:32" json:"hash"`
}

func (i ImageModel) WebPath() string {
	return fmt.Sprintf("/" + i.Path)
}

func (l ImageModel) BeforeDelete(tx *gorm.DB) error {
	err := os.Remove(l.Path)
	if err != nil {
		logrus.Warnf("删除文件失败 %s", err)
	}
	return nil
}
