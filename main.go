package main

import (
	"database/sql"
	"embed"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // SQLite驱动
)

//go:embed index.html
var staticFS embed.FS // 内置前端HTML文件

var db *sql.DB // 全局数据库连接

// Group 分组结构体
type Group struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Links []Link `json:"links"`
}

// Link 链接结构体
type Link struct {
	ID   int    `json:"id"`
	GID  int    `json:"gid"`
	Text string `json:"text"`
	URL  string `json:"url"`
}

// 初始化数据库（自动建表+默认用户）
func initDB() {
	var err error
	// 连接SQLite数据库（文件：navigo.db）
	db, err = sql.Open("sqlite3", "./data/navigo.db")
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	// 创建用户表（单用户）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`)
	if err != nil {
		panic("创建用户表失败：" + err.Error())
	}

	// 创建分组表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	)`)
	if err != nil {
		panic("创建分组表失败：" + err.Error())
	}

	// 创建链接表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		gid INTEGER NOT NULL,
		text TEXT NOT NULL,
		url TEXT NOT NULL,
		FOREIGN KEY(gid) REFERENCES groups(id) ON DELETE CASCADE
	)`)
	if err != nil {
		panic("创建链接表失败：" + err.Error())
	}

	// 初始化默认用户（仅第一次运行）
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		panic("查询用户数失败：" + err.Error())
	}
	if count == 0 {
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (1, 'team', '123456')")
		if err != nil {
			panic("初始化默认用户失败：" + err.Error())
		}
	}
}

// 登录验证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "ok" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录，请先登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	// 初始化数据库
	initDB()
	defer db.Close() // 程序退出时关闭数据库

	// 创建Gin引擎
	r := gin.Default()
	// 允许跨域
	r.Use(cors.Default())

	// 提供前端HTML页面
	r.GET("/", func(c *gin.Context) {
		data, err := staticFS.ReadFile("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "页面加载失败")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// ========== API接口 ==========
	api := r.Group("/api")

	// 登录接口
	api.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}

		var password string
		err := db.QueryRow("SELECT password FROM users WHERE username=?", req.Username).Scan(&password)
		if err != nil || password != req.Password {
			c.JSON(http.StatusOK, gin.H{"msg": "账号或密码错误"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": "ok"})
	})

	// 获取所有分组和链接
	api.GET("/groups", func(c *gin.Context) {
		var groups []Group
		// 查询所有分组
		rows, err := db.Query("SELECT id, name FROM groups ORDER BY id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "查询分组失败"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var g Group
			if err := rows.Scan(&g.ID, &g.Name); err != nil {
				continue
			}

			// 查询该分组下的所有链接
			linkRows, err := db.Query("SELECT id, gid, text, url FROM links WHERE gid=? ORDER BY id", g.ID)
			if err != nil {
				continue
			}
			defer linkRows.Close()

			for linkRows.Next() {
				var l Link
				if err := linkRows.Scan(&l.ID, &l.GID, &l.Text, &l.URL); err != nil {
					continue
				}
				g.Links = append(g.Links, l)
			}

			groups = append(groups, g)
		}

		c.JSON(http.StatusOK, groups)
	})

	// 分组CRUD
	api.POST("/group", authMiddleware(), func(c *gin.Context) {
		var g Group
		if err := c.ShouldBindJSON(&g); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}
		_, err := db.Exec("INSERT INTO groups (name) VALUES (?)", g.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建分组失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "创建成功"})
	})

	api.PUT("/group/:id", authMiddleware(), func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}
		var g Group
		if err := c.ShouldBindJSON(&g); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}
		_, err = db.Exec("UPDATE groups SET name=? WHERE id=?", g.Name, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "修改分组失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
	})

	api.DELETE("/group/:id", authMiddleware(), func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}
		_, err = db.Exec("DELETE FROM groups WHERE id=?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除分组失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
	})

	// 链接CRUD
	api.POST("/link", authMiddleware(), func(c *gin.Context) {
		var l Link
		if err := c.ShouldBindJSON(&l); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}
		_, err := db.Exec("INSERT INTO links (gid, text, url) VALUES (?, ?, ?)", l.GID, l.Text, l.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建链接失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "创建成功"})
	})

	api.PUT("/link/:id", authMiddleware(), func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}
		var l Link
		if err := c.ShouldBindJSON(&l); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}
		_, err = db.Exec("UPDATE links SET text=?, url=? WHERE id=?", l.Text, l.URL, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "修改链接失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
	})

	api.DELETE("/link/:id", authMiddleware(), func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			return
		}
		_, err = db.Exec("DELETE FROM links WHERE id=?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除链接失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
	})

	// 启动服务（监听3000端口）
	r.Run(":3000")
}
