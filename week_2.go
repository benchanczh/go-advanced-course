package main

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

/**
应该Wrap这个error抛给上层的原因：
1. 防止获取数据信息的指针为空被业务层使用，引起空指针异常或其他错误；
2. 业务层能够获取数据库查询错误的根本原因，可以将错误原因映射到用户层
3. 记录sql查询语句的日志
4. 保存错误根因
5. 业务在调用第三库或标准库时，应该使用errors.Wrap 保存堆栈信息及错误详细详细，方便业务层问题定位
6. 根据查询的结果进行业务处理，如果查询结果为零不进行相关业务处理，忽略或进行其他操作
**/

var db *sql.DB

type User struct {
	Id   string
	Name string
}

func QueryUserName(name string) (*User, error) {
	var user User
	sql := "SELECT name FROM users WHERE id = ?"
	err := db.QueryRow(sql, name).Scan(&user.Id, &user.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "main:QueryUserName Scan failed:%s,%s", sql, name)
	}
	return &user, nil
}

func handleSQLErrNoRows() {
	// do something
}

func main() {
	name := "foo"
	user, err := QueryUserName(name)
	if err != nil {
		fmt.Printf("Original error, %v\n", errors.Cause(err))
		fmt.Printf("Stack trace:\n%+v\n", err)
		if errors.Cause(err) == sql.ErrNoRows {
			handleSQLErrNoRows()
		}
		return
	}

	fmt.Print(user)
}
