package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//var db *sql.DB
type Mydb struct{
	DB *sql.DB
}
//连接数据库
func getConn(dataSourceName string)(*Mydb, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil,err
	}
	db.SetConnMaxLifetime(time.Minute*50)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.Stats()
	mydb:=&Mydb{DB:db}
	return mydb,nil
}
//查询语句，返回查询的结果集
func (db *Mydb) Query(query string, args... interface{})(*sql.Rows,error){
	rows,err:=db.DB.Query(query,args...)
	if err != nil {
		return nil,err
	}
	return rows,nil
}
//插入数据，返回插入数据的id
func (db *Mydb) Insert(strSql string)(int64,error){
	res,err:=db.DB.Exec(strSql)
	if err != nil {
		return 0,err
	}
	insertId, _ := res.LastInsertId()
	return insertId,nil
}
//删除数据，返回被影响的行
func (db *Mydb) Delete(strSql string,args...interface{})(int64,error){
	res,err:=db.DB.Exec(strSql,args...)
	if err != nil {
		return 0,err
	}
	var affNum int64
	affNum, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affNum,nil
}
//更改数据，返回所影响的行
func (db *Mydb) Update(strSql string,args...interface{})(int64,error){
	res,err:=db.DB.Exec(strSql,args...)
	if err != nil {
		return 0,err
	}
	var affNum int64
	affNum, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affNum,nil
}
//封装事务
func ExecTransaction(db *Mydb,handle func(tx *sql.Tx)error)error{
	tx,err:=db.DB.Begin()
	if err != nil {
		return err
	}
	defer func(){
		if err!=nil{
			tx.Rollback()
		}
	}()
	if err=handle(tx);err!=nil{
		return err
	}
	return tx.Commit()
}
func main(){
	//连接数据库
	mydb,_:=getConn("root:Dyf5201314@tcp(127.0.0.1:3306)/students?charset=utf8")
	//例1: 普通查询
	rows,_:=mydb.DB.Query("select * from info_stu")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
	//例2：使用事务进行查询操作
	handle:=func(tx *sql.Tx)error{
		rows, err := tx.Query("select * from info_stu where id=?;",1)
		for rows.Next() {
			var id int
			var name string
			rows.Scan(&id, &name)
			fmt.Printf("id:%d name:%s\n", id, name)
		}
		return err
	}
	ExecTransaction(mydb,handle)  //结果 ： id:1 name:Duyifan
}