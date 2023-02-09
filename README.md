# go-sqlink

DecodeRows decodes SQL result rows into slice of structs.
It scans struct field tags to match column name in result row.
Struct field data type should match column value data type too.

Usage:

Consider database table:

```
CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL
);
```

We want to scan rows into slices of struct:

```
// tag name should match column name of result row
// data type should match too
type User struct {
	Id   int    `sql:"id"`
	Name string `sql:"name"`
}

func Process() {
    // database connection, etc

	rows, err := db.Query("SELECT * FROM `users`")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	users := make([]User, 0)

	err = sqlink.DecodeRows(rows, &users)
	if err != nil {
		log.Fatal(err)
	}

    // other data processing code
}
```
