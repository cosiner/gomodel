# gomodel
```sh
$ gomodel [OPTIONS] DIR|FILES...
```

# Structure tags
* `table`: table name
* `column`: column name

Both using "`-`" to prevent from parsing.

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
    Example: 
    replace 
    ```sql
    SELECT user_id, name, password FROM user
    ```
    with 
    ```sql
    SELECT Id, Name, Password FROM User
    ```

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
The form of first and last line can't be changed.
