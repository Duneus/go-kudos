package sqlite

import "github.com/jinzhu/gorm"

type storage struct {
	gorm *gorm.DB
}

func Migrate(orm *gorm.DB) error {
	if err := orm.Set("gorm:table_options", "CASCADE").DropTableIfExists(
		&kudos{},
		&schedule{},
	).Error
		err != nil {
			panic(err)
	}

	return orm.AutoMigrate(&kudos{}, &schedule{}).Error
}
