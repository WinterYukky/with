package with

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// With is a with clause.
//
// .ex
//	DB.Clauses(New(DB).Append("`buyers` AS (SELECT `user_id` FROM `sales`)")).Where("users.id IN (SELECT * FROM buyers)").Find(&User{})
// Make this query.
//
// WITH `buyers` AS (SELECT user_id FROM `sales`) SELECT * FROM `users` WHERE users.id IN (SELECT * FROM `buyers`)
type With struct {
	tx      *gorm.DB
	queries []withQuery
}
type withQuery struct {
	exprs []clause.Expression
}

// ModifyStatement implements gorm interface
func (with With) ModifyStatement(stmt *gorm.Statement) {
	if len(with.queries) == 0 {
		return
	}
	clause := stmt.Clauses["SELECT"]
	clause.BeforeExpression = with
	stmt.Clauses["SELECT"] = clause
}

// Build implements gorm interface
func (with With) Build(builder clause.Builder) {
	builder.WriteString("WITH ")

	for index, query := range with.queries {
		if index > 0 {
			builder.WriteString(", ")
		}
		for _, expr := range query.exprs {
			expr.Build(builder)
		}
	}
}

// Append a with clause.
func (with With) Append(query string, args ...interface{}) With {
	if conds := with.tx.Statement.BuildCondition(query, args...); len(conds) > 0 {
		with.queries = append(with.queries, withQuery{
			exprs: conds,
		})
	}
	return with
}

// New create a with clause.
func New(tx *gorm.DB) With {
	return With{tx: tx}
}
