package model

import "gorm.io/gorm"

type Microservices struct {
	gorm.Model
	UUID        string `gorm:"column:uuid;type:varchar(60);UNIQUE;NOT NULL;" json:"uuid"`
	Name        string `gorm:"column:name;type:varchar(60);DEFAULT '';" json:"name"`
	Host        string `gorm:"column:host;type:varchar(60);DEFAULT '';" json:"host"`
	Port        string `gorm:"column:port;type:varchar(10);DEFAULT '';" json:"port"`
	Group       string `gorm:"column:group;type:varchar(60);DEFAULT '';" json:"group"`
	IsInternal  bool   `gorm:"column:isinternal;type:bool;DEFAULT false;" json:"isinternal"`
	IsSession   bool   `gorm:"column:issession;type:bool;DEFAULT false;" json:"issession"`
	IsActive    bool   `gorm:"column:isactive;type:bool;DEFAULT true;" json:"isactive"`
	MasterID    string `gorm:"column:masterid;type:varchar(60);DEFAULT '';" json:"masterid"`
	Echo        int64  `gorm:"column:echo;type:uint;DEFAULT '';" json:"echo"`
	Description string `gorm:"column:description;type:longtext;DEFAULT '';" json:"description"`
	Author      string `gorm:"column:author;type:varchar(60);DEFAULT '';" json:"author"`
	Version     string `gorm:"column:version;type:varchar(60);DEFAULT '';" json:"version"`
}
