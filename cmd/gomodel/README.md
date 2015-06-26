# gomodel
```sh
$ gomodel [OPTIONS] DIR|FILES...
```

# Structure tags
* `table`: table name
* `column`: column name

Both using "`-`" to prevent from parsing.

# SQL convertion
### Why
* I don't like to build sql in program manually
* I don't mention what the real table and column name is
* Just writing sql

### Synax
```Go
type User struct {
    Id       string `table:"user" column:"user_id"`
    Name     string
    Password string
}
```

* **AST parsing**: parsing sql AST  
    Just use structure and field name replace table and column name.
    Example: replace 
    ```sql
    SELECT user_id, name, password FROM user
    ```
    with 
    ```sql
    SELECT Id, Name, Password FROM User
    ```

* **Simple parsing**: simple lexing 
    + `{Structure}` as table 
    + `{Structure:Field, Field}` as column, column
    + `{Structure.Field, Field}` as table.column, table.column
    
    Example: replace 
    ```sql
    SELECT user_id, name, password FROM user
    ```
    with 
    ```sql 
    SELECT {User:Id, Name, Password} FROM {User}
    ```

The only reason to use simple parsing is that AST parsing is not enough tested, 
AST parsing is enabled by default, use `-AST=false` to enable simple parsing.

### Position
sql should be put in function documents, and start with `//gomodel `
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
The form of first and last line can't be changed, others will be joined with an " ".

# Output
see model.tmpl.
