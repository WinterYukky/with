package with

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// With is with clause.
//
// .ex
//	db.Table("users").Clause(with.New("buyers", db.Table("sales").Select("user_id"))).Where("users.id IN (db.Table("buyers"))"))
// Make this query.
//
// WITH `buyers` AS (SELECT * FROM `sales`) SELECT * FROM `users` WHERE users.id IN (SELECT user_id FROM buyers)
type With []query
type query struct {
	name     string
	subquery *gorm.DB
}

// ModifyStatement implements gorm interface
func (with With) ModifyStatement(stmt *gorm.Statement) {
	if len(with) == 0 {
		return
	}
	clause := stmt.Clauses["SELECT"]
	clause.BeforeExpression = with
	stmt.Clauses["SELECT"] = clause
}

// Build implements gorm interface
func (with With) Build(builder clause.Builder) {
	builder.WriteString("WITH ")
	for index, query := range with {
		if index > 0 {
			builder.WriteString(", ")
		}
		builder.WriteQuoted(query.name)
		builder.WriteString(" AS (" + query.subquery.Session(&gorm.Session{DryRun: true}).Find(nil).Statement.SQL.String() + ")")

	}
}

// Append a With clause.
func (with With) Append(name string, subquery *gorm.DB) With {
	return append(with, query{
		name:     name,
		subquery: subquery,
	})
}

// New create a With clause.
func New() With {
	return With{}
}
