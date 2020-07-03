package sqlite

type kudos struct {
	Id      string `gorm:"primary_key"`
	Message string
}
