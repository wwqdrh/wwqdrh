package main

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		log.Fatal(err)
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create persons table
		{
			ID: "201608301400",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time
				type Person struct {
					gorm.Model
					Name string
				}
				return tx.AutoMigrate(&Person{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("people")
			},
		},
		// add age column to persons
		{
			ID: "201608301415",
			Migrate: func(tx *gorm.DB) error {
				// when table already exists, it just adds fields as columns
				type Person struct {
					Age int
				}
				return tx.AutoMigrate(&Person{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn("people", "age")
			},
		},
		// add pets table
		{
			ID: "201608301430",
			Migrate: func(tx *gorm.DB) error {
				type Pet struct {
					gorm.Model
					Name     string
					PersonID int
				}
				return tx.AutoMigrate(&Pet{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("pets")
			},
		},
		{
			ID: "202202231130",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`create table table3 (name text, info text)`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("table3")
			},
		},
	})

	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
}
