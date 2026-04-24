package entity

type Page struct {
	Model
	Name  string `gorm:"unique;type:varchar(20)" json:"name"`
	Label string `gorm:"unique;type:varchar(30)" json:"label"`
	Cover string `gorm:"type:varchar(255)" json:"cover"`
}

type Config struct {
	Model
	Key   string `gorm:"unique;type:varchar(256)" json:"key"`
	Value string `gorm:"type:varchar(256)" json:"value"`
	Desc  string `gorm:"type:varchar(256)" json:"desc"`
}
