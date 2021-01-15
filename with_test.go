package with

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
	DryRun: true,
})

type User struct {
	ID   uint
	Name string
}

type Sale struct {
	ID      uint
	UserID  uint
	Product string
}

func AssertSQL(t *testing.T, result *gorm.DB, sql string) {
	if result.Statement.SQL.String() != sql {
		t.Errorf("SQL expects: %v, got %v", sql, result.Statement.SQL.String())
	}
}

func TestWith(t *testing.T) {
	result := DB.Clauses(New()).Find(&User{})

	AssertSQL(t, result, "SELECT * FROM `users`")

	with := New().Append("apple_buyers", DB.Model(&Sale{}).Select("user_id").Where("product = ?", "apple"))
	result = DB.Clauses(with).Where("users.id IN (SELECT * FROM apple_buyers)").Find(&User{})
	AssertSQL(t, result, "WITH `apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = ?) SELECT * FROM `users` WHERE users.id IN (SELECT * FROM apple_buyers)")

	with = New().Append("apple_buyers", DB.Model(&Sale{}).Select("user_id").Where("product = ?", "apple")).
		Append("orange_buyers", DB.Model(&Sale{}).Select("user_id").Where("product = ?", "orange"))
	result = DB.Clauses(with).
		Where("users.id IN (SELECT * FROM apple_buyers)").
		Where("users.id IN (SELECT * FROM orange_buyers)").Find(&User{})
	AssertSQL(t, result, "WITH `apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = ?), `orange_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = ?) SELECT * FROM `users` WHERE users.id IN (SELECT * FROM apple_buyers) AND (users.id IN (SELECT * FROM orange_buyers))")
}
