# sql convertion
##### place
sql should be written in function documents, and start with `//gomodel `

##### syntax
* simple
```
//gomodel sqlname = SELECT Id, Name, Password FROM User
```

* block
```
//gomodel sqlname = [
//  SELECT
//  Id, Name, Password
//  FROM
//  User
//]
```
The form of first and last line can't be changed

# sql syntax
##### ast parsing
Just use structure and field name replace table and column name

##### simple parsing
* `{Structure}` as table 
* `{Structure:Field, Field}` as column, column
* `{Structure.Field, Field}` as table.column, table.column




