# with
With support for GORM

## Use a With clause
```go
import "github.com/WinterYukky/with"

with := with.New(db).
    Append("`apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = ?)", "apple").
    Append("`orange_buyers` AS (?)", db.Model(&Sale{}).Select("user_id").Where("product = ?", "orange"))
db.Clauses(with).
    Where("users.id IN (?)", db.Table("apple_buyers")).
    Where("users.id IN (SELECT * FROM `orange_buyers`)")).Find(&User{})

// WITH 
//   `apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'apple'),
//   `orange_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'orange')
// SELECT * FROM `users` WHERE users.id IN (SELECT * FROM `apple_buyers`) AND users.id IN (SELECT * FROM `orange_buyers`)
```