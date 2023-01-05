# Struct2SQL
A SQL builder & result scanner for Go. This library is made for study propose. Use it on your own risk.

### What's inside?
* Struct analyzer - mapping your struct to a table in the database
* Easy to use API like `db.Select(&user, "id = ?", 1)`
* Result scanner that scan your result to struct
* Export struct data to CSV
* SQL Build to help you generate Query

## Getting Started
### Start the connection
```
import s2s "github.com/mingchoi/struct2sql"

// Connect to database
db, err := s2s.Open("mysql", "gomysql:123456@/gomysql?charset=utf8&parseTime=True&loc=Local")
if err != nil {
    panic(err)
}
defer db.Close()
```

### Define a struct & create table
```
type Monster struct {
	ID     int `keyword:"NOT NULL AUTO_INCREMENT" primarykey:"true"`
	Name   string
	Level  int
	Health float64
	Boss   bool
	Spawn  time.Time
	MID    int       `foreignkey:"map(id)"`
}

// Create table
db.CreateTable(&Monster{})
```

### Create
```
// Insert
monster := Monster{
    Name:   "Monster 怪物",
    Level:  18,
    Health: 1029.8,
    Boss:   true,
    Spawn:  time.Now(),
    MID:    1,
}
err = db.Insert(&monster)
if err != nil {
    panic(err)
}
```

### Read
```
// Select
sm := Monster{}
err = db.Select(&sm, "id = ?", 1)
if err != nil {
    panic(err)
}
```

### Update
```
// Update
sm.Level = 21
err = db.Update(&sm, "id = ?", sm.ID)
if err != nil {
    panic(err)
}
```

### Delete
```
// Delete
err = db.Delete(&Monster{}, "id = ?", 1)
if err != nil {
    panic(err)
}
```

## For more detail
Please read the test case for the usage
