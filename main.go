package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "go-sql-driver/mysql"
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/register", register)
	http.ListenAndServe(":8080", mux) //监听8080端口
}

type result struct { //定义返回数据格式
	Code int
	Msg  string
	Data []string
}

func login(w http.ResponseWriter, r *http.Request) { //登陆
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许跨域
	r.ParseForm()
	username, found1 := r.Form["username"]
	password, found2 := r.Form["password"]
	if !(found1 && found2) {
		io.WriteString(w, "请勿非法访问")
		return
	}
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test?charset=utf8") //登陆msyql
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var post_data PostData
	post_data.user_name = username[0]
	post_data.pass_word = password[0]

	var filter_data FilterData = post_data
	post_data = filter_data.formatData()

	var row *sql.Row
	row = db.QueryRow("select * from users where username = ? and password = ?", post_data.user_name, post_data.pass_word)
	var user_name, pass_word string
	var id int
	err = row.Scan(&id, &user_name, &pass_word) //遍历结果
	if err != nil {
		arr := &result{
			500,
			"登陆失败",
			[]string{},
		}
		b, json_err := json.Marshal(arr) //json化结果集
		if json_err != nil {
			fmt.Println("encoding faild")
		} else {
			io.WriteString(w, string(b)) //返回结果
		}
	} else {
		arr := &result{
			200,
			"登陆成功",
			[]string{},
		}
		b, json_err := json.Marshal(arr) //json化结果集
		if json_err != nil {
			fmt.Println("encoding faild")
		} else {
			io.WriteString(w, string(b)) //返回结果
		}
	}
}

func register(w http.ResponseWriter, r *http.Request) { //注册
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许跨域
	r.ParseForm()
	username, found1 := r.Form["username"]
	password, found2 := r.Form["password"]
	if !(found1 && found2) {
		io.WriteString(w, "请勿非法访问")
		return
	}
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		io.WriteString(w, "连接数据库失败")
		return
	}
	defer db.Close() //在返回前关闭资源（延迟）

	var post_data PostData
	post_data.user_name = username[0]
	post_data.pass_word = password[0]

	var filter_data FilterData = post_data
	post_data = filter_data.formatData()

	_, err = db.Exec("insert into users(username, password) values(?,?)", post_data.user_name, post_data.pass_word) //插入数据
	if err != nil {
		arr := &result{
			500,
			"注册失败",
			[]string{},
		}
		b, json_err := json.Marshal(arr) //json化结果集
		if json_err != nil {
			fmt.Println("encoding faild")
		} else {
			io.WriteString(w, string(b)) //返回结果
		}
	} else {
		arr := &result{
			200,
			"注册成功",
			[]string{},
		}
		b, json_err := json.Marshal(arr) //json化结果集
		if json_err != nil {
			fmt.Println("encoding faild")
		} else {
			io.WriteString(w, string(b)) //返回结果
		}
	}
}

type FilterData interface { //定义数据接口
	formatData() PostData
}

type PostData struct { //定义传递数据结构体
	user_name string
	pass_word string
}

func (post_data PostData) formatData() PostData { //格式化数据
	post_data.pass_word = mdFormat(post_data.pass_word)
	return post_data
}

func mdFormat(data string) string { //对字符串进行md5加密
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
