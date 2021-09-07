package models

import (
	"fmt"
	"net/http"
	"time"
)

type Images struct {
	Id          int    `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	ImageName   string `xorm:"VARCHAR(100) notnull unique(unique_name_version)" json:"imageName"`
	Version     string `xorm:"VARCHAR(100) notnull unique(unique_name_version)" json:"version"`
	Status      string `xorm:"VARCHAR(100) notnull" json:"status"`
	Description string `xorm:"VARCHAR(300)" json:"description"`
	ImageTypeId int    `xorm:"INTEGER" json:"-"`

	// json序列化字段
	Type     string `xorm:"-" json:"type"`
	Category string `xorm:"-" json:"category"`
	// 从系统参数中查询
	Address string `xorm:"-" json:"address"`
}

type Images2 struct {
	Images   `xorm:"extends"`
	Category string `json:"category"`
}

// 设置镜像状态
func (i *Images) SetStatus(address string) {
	checkUrl := fmt.Sprintf("http://%v/api/repositories/%v/tags/%v", address, i.ImageName, i.Version)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(checkUrl)
	if err == nil {
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			i.Status = "Valid"
			return
		}
	}
	i.Status = "Invalid"
}
