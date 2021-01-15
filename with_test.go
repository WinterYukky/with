package with

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID   uint
	Name string
}

type Sale struct {
	ID      uint
	UserID  uint
	Product string
}

func TestWith(t *testing.T) {
	assertSQL := func(t *testing.T, result *gorm.DB, sql string) {
		if result.Statement.SQL.String() != sql {
			t.Errorf("SQL expects: %v, got %v", sql, result.Statement.SQL.String())
		}
	}
	var db, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DryRun: true,
	})
	result := db.Clauses(New(db)).Find(&User{})
	assertSQL(t, result, "SELECT * FROM `users`")

	with := New(db).Append("`apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'apple')")
	result = db.Clauses(with).Where("users.id IN (?)", db.Table("apple_buyers")).Find(&User{})
	assertSQL(t, result, "WITH `apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'apple') SELECT * FROM `users` WHERE users.id IN (SELECT * FROM `apple_buyers`)")

	with = New(db).Append("`apple_buyers` AS (?)", db.Model(&Sale{}).Select("user_id").Where("product = ?", "apple"))
	result = db.Clauses(with).Where("users.id IN (?)", db.Table("apple_buyers")).Find(&User{})
	assertSQL(t, result, "WITH `apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = ?) SELECT * FROM `users` WHERE users.id IN (SELECT * FROM `apple_buyers`)")

	with = New(db).Append("`apple_buyers` AS (?)", db.Model(&Sale{}).Select("user_id").Where("product = ?", "apple")).
		Append("`orange_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'orange')")
	result = db.Clauses(with).
		Where("users.id IN (?)", db.Table("apple_buyers")).
		Where("users.id IN (SELECT * FROM `orange_buyers`)").Find(&User{})
	assertSQL(t, result, "WITH `apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = ?), `orange_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'orange') SELECT * FROM `users` WHERE users.id IN (SELECT * FROM `apple_buyers`) AND (users.id IN (SELECT * FROM `orange_buyers`))")
}
