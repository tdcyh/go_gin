package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func if_login()gin.HandlerFunc {
	return func(context *gin.Context) {
		val,err := context.Cookie("name")
		if err != nil{
			context.JSON(200,gin.H{"notLogin 去登录":"http://localhost:8080/addUser"})
			//context.String(200, "Cookie:%s:", val)

			context.Abort()

			return
		}else{
			context.JSON(200,gin.H{val:"Login"})
			context.String(200, "Cookie:%s:", val)

			context.Next()
		}

	}
}

func main() {

	db, _ := sqlx.Open("mysql", "root:liemaren0.0@tcp(127.0.0.1:3306)/go")

	//初始
	g := gin.Default()
	g.LoadHTMLGlob("templates/*")

	g.GET("/index",if_login(), func(context *gin.Context) {
		context.JSON(200,gin.H{"111":"1232"})
	})

	g.GET("/zhuce", func(c *gin.Context) {

		c.HTML(http.StatusOK, "zhuce.html",nil)

	})
	g.POST("dozhuce", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		question := c.PostForm("question")
		answer := c.PostForm("answer")
		_, err := db.Exec("insert into users (username ,password,question,answer) value (?,?,?,?)", username,password,question,answer)
		if err != nil {
			log.Println(err)
			return
		}


		c.JSON(200,gin.H{"注册成功,去登录:":"http://localhost:8080/addUser"})

	})


	//通过 c.PostForm 接收表单传过来的数据
	g.GET("/addUser", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html",nil)

	})

	g.POST("doAddUser", func(c *gin.Context) {
			flag := false
			username := c.PostForm("username")
			password := c.PostForm("password")
			rows, err := db.Query("select username,password from users")
			//在 Query 语句中 ? 表示占位符  在这里 同样的效果就是
			//   select * from stu where id= 1
			if err != nil {
				log.Println(err)
				return
			}
			//  延迟调用关闭rows释放持有的数据库链接
			defer rows.Close()
			var user struct {
				userName_sql   string
				password_sql    string

			}
			//  迭代查询获取数据  必须调用
			for rows.Next() {
				// row.scan 必须按照先后顺序 &获取数据
				err := rows.Scan(&user.userName_sql, &user.password_sql)
				if err != nil {
					log.Println(err)
					return
				}

				if user.userName_sql == username && user.password_sql == password{
					flag = true
				}
			}
			if flag{
				c.SetCookie("name", password, 60, "/", "localhost", false, true)
				val, _ := c.Cookie("name")
				c.String(200, "Cookie:%s:", val)
				c.JSON(200,gin.H{"登陆成功,重新进入index":"http://localhost:8080/index"})

			} else {
				c.JSON(200,gin.H{"error":"wrong pass or no this user"})
				c.JSON(200,gin.H{"去以下网站注册":"http://localhost:8080/zhuce"})
				c.JSON(200,gin.H{"去以下网站找回密码":"http://localhost:8080/mibao"})

			}


	})

	g.GET("/mibao", func(c *gin.Context) {
		c.HTML(http.StatusOK, "template.html",nil)

	})
	g.POST("domibao", func(c *gin.Context) {
		flag := false
		username := c.PostForm("username")
		question := c.PostForm("question")
		answer := c.PostForm("answer")
		rows, err := db.Query("select * from users")
		//在 Query 语句中 ? 表示占位符  在这里 同样的效果就是
		//   select * from stu where id= 1
		if err != nil {
			log.Println(err)
			return
		}
		//  延迟调用关闭rows释放持有的数据库链接
		defer rows.Close()
		var user struct {
			userName_sql   string
			password_sql   string
			question_sql    string
			answer_sql    string

		}
		//  迭代查询获取数据  必须调用
		for rows.Next() {
			// row.scan 必须按照先后顺序 &获取数据
			err := rows.Scan(&user.userName_sql,&user.password_sql, &user.question_sql, &user.answer_sql)
			if err != nil {
				log.Println(err)
				return
			}

			if user.userName_sql == username && user.question_sql == question &&user.answer_sql == answer{
				flag = true
			}
		}
		if flag{
			db.Exec("update users set password= ? where username=?", "000000", user.userName_sql)

			c.JSON(200,gin.H{"已将密码重置为":"000000"})
			c.JSON(200,gin.H{"去以下网站登录":"http://localhost:8080/addUser"})


		} else {
			c.JSON(200,gin.H{"error":"wrong answer"})
			c.JSON(200,gin.H{"去以下网站注册":"http://localhost:8080/zhuce"})

		}


	})


	g.Run()
}

