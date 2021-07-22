package models

import (
	"errors"
	"github.com/chujieyang/commonops/ops/infrastructure/database"
)

type Nacos struct {
	Id    int64   `json:"Id" gorm:"column:id;type:int;PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	DataStatus    int8   `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	EndPoint  string `json:"EndPoint" gorm:"column:end_point;type:varchar(255)"`
	Alias  string `json:"Alias" gorm:"column:alias;type:varchar(255)"`
	Username string `json:"Username" gorm:"column:username;type:varchar(255)"`
	Password  string `json:"Password" gorm:"column:password;type:varchar(255)"`
}

func (Nacos) TableName() string {
	return "nacos"
}

func AddNewNacosServer(alias, endpoint, username, password string) (id int, err error) {
	count := 0
	if err = database.Mysql().Raw("select count(*) from nacos where end_point = ? and data_status = 1", endpoint).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		err = errors.New("已存在相同EndPoint的集群")
		return
	}
	nacos := &Nacos{
		Alias: alias,
		EndPoint: endpoint,
		Username: username,
		Password: password,
	}
	err = database.Mysql().Create(nacos).Error
	id = int(nacos.Id)
	return
}

func GetNacosInfoById(id string) (info Nacos, err error) {
	err = database.Mysql().Raw("select * from nacos where id = ? limit 1", id).Scan(&info).Error
	return
}

func GetNacosList() (data []Nacos, err error) {
	querySql := "select id, end_point, alias from nacos where data_status = 1"
	err = database.Mysql().Raw(querySql).Scan(&data).Error
	return
}
