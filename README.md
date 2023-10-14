# trier
[![Go Reference](https://pkg.go.dev/badge/github.com/syke99/trier.svg)](https://pkg.go.dev/github.com/syke99/trier)
[![go reportcard](https://goreportcard.com/badge/github.com/syke99/trier)](https://goreportcard.com/report/github.com/syke99/trier)
[![License](https://img.shields.io/github/license/syke99/trier)](https://github.com/syke99/trier/blob/master/LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/syke99/trier)</br>
a heavily Zig-inspired approach to error handling in Go

Why use `trier`?
===
One of Go's biggest complaints is verbose error handling. You might have a function with 10 or more `if err != nil` checks, causing functions to bloat. With `trier`, you can create a trier and then chain function calls and return the first error encountered without having to clutter your code with `if err != nil` checks everywhere. It's heavily inspired by Zig's `try` keyword that doesn't use `catch` blocks. Instead, it internally stores and checks errors before trying the next chained function. You can chain as many functions together as you want, and then do your own final check on nil errors by calling `tr.Err()` (example below)

How do I use trier?
====

### Installation

```
go get github.com/syke99/trier
```

### Basic Usage

```go
package main

import (
    "errors"

    "github.com/syke99/trier"
)

// create the functions you want to be
// tried (they can also just be passed 
// as anonymous functions)
func passOrFail(args ...any) error {
    if len(args) != 0 {
        return errors.New("failed passOrFail")
    }
    return nil
}

func failIfString(args ...any) error {
    var err error
    
    switch args[0].(type) {
    case string:
        err = errors.New("failedIfString")
    }
    return err
}

func main() {
    // create a new trier by calling trier.NewTrier()
    tr := trier.NewTrier()
    
    // try your functions by 
    // chaining them in tr.Try() calls
    // and passing the appropriate args
    // (in this example, even though the
    // second time trying failIfString 
    // should return an error, the trier's
    // error will be "failed passOrFail"
    // because the second time trying
    // passOrFail returned an error before
    // the second time trying failIfString)
    tr.Try(passOrFail).
        Try(failIfString, 0).
        Try(passOrFail, true).
        Try(failIfString, "hi")
    
    // prints "failed passOrFail
    println(tr.Err().Error())
}
```

### More Advanced Usage

```go
package main

import (
	"database/sql"
        "errors"
	"fmt"


	"github.com/syke99/trier"
	_ "github.com/go-sql-driver/mysql"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func main() {
        // create a new trier by calling trier.NewTrier()
        tr := trier.NewTrier()

	var err error
	
	var db *sql.DB
	
	var tx *sql.Tx
	
	tr.Try(func(args ...any) error {
		db, err = sql.Open("mysql", dsn(""))
		return err
	}).Try(func(args ...any) error {
                tx, err = db.Begin()
		return err
	}).Try(func(args ...any) error {
                _, err = tx.Exec("UPDATE customers SET name = \"Jane Doe\" WHERE ID = 1")
		return err
	})
	
	if err != nil {
		if tx != nil {
			tx.Rollback()
                }
		panic(err)
        }
	
	tx.Commit()
}
```

More examples (such as anonymous functions, `*.TryJoin()`, etc) can be found [here](https://github.com/syke99/trier/blob/main/trier_test.go)

Who?
====

This library was developed by Quinn Millican ([@syke99](https://github.com/syke99))


## License

This repo is under the MIT license, see [LICENSE](LICENSE) for details.
