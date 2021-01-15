# with
With support for GORM

## Single CTE
```go
import "github.com/WinterYukky/with"

with := New("apple_buyers", DB.Model(&Sale{}).Select("user_id").Where("product = ?", "apple"))
DB.Clauses(with).Where("users.id IN (?)", DB.Table("apple_buyers")).Find(&User{})
// WITH `apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'apple')
// SELECT * FROM `users` WHERE users.id IN (SELECT * FROM apple_buyers)
```

## Multiple CTE
```go
import "github.com/WinterYukky/with"

with := New("apple_buyers", DB.Model(&Sale{}).Select("user_id").Where("product = ?", "apple")).
    Append("orange_buyers", DB.Model(&Sale{}).Select("user_id").Where("product = ?", "orange"))
DB.Clauses(with).
    Where("users.id IN (?)", DB.Table("apple_buyers")).
    Where("users.id IN (?)", DB.Table("orange_buyers")).Find(&User{})
// WITH 
// `apple_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'apple'),
// `orange_buyers` AS (SELECT `user_id` FROM `sales` WHERE product = 'orange'),
// SELECT * FROM `users` WHERE users.id IN (SELECT * FROM apple_buyers)
```