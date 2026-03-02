# NaviGo

轻量、可靠、团队友好的内部网址导航系统（Go + SQLite）。

## ✨ 核心特性
- 🚀 **纯Go实现**：自带HTML页面服务，无需Nginx/Apache
- 📦 **SQLite数据库**：单文件存储，备份/迁移超简单
- 👥 **团队协作**：单用户账号，登录可编辑，游客仅查看
- 🐳 **Docker部署**：一键启动，仅需挂载数据库文件
- 📉 **极低资源占用**：内存仅10-20MB，适合低配服务器
- 🌐 **跨平台**：支持Linux/Windows/macOS

## 🚀 快速部署（Docker Compose）
### 1. 创建docker-compose.yml
```yaml
version: '3'
services:
  navigo:
    image: 你的Docker用户名/navigo:latest
    container_name: navigo
    ports:
      - "80:3000"  # 主机80端口映射到容器3000端口
    volumes:
      - ./navigo.db:/app/navigo.db  # 挂载数据库文件到宿主机
    restart: always  # 开机自启
    # NaviGo
轻量、可靠、团队友好的内部网址导航系统（Go + SQLite）。

## ✨ 核心特性
- 🚀 **纯Go实现**：自带HTML页面服务，无需Nginx/Apache
- 📦 **SQLite数据库**：单文件存储，备份/迁移超简单
- 👥 **团队协作**：单用户账号，登录可编辑，游客仅查看
- 🐳 **Docker部署**：一键启动，仅需挂载数据库文件
- 📉 **极低资源占用**：内存仅10-20MB，适合低配服务器
- 🌐 **跨平台**：支持Linux/Windows/macOS

## 🚀 快速部署（Docker Compose）

### 1. 创建docker-compose.yml
```yaml
version: '3'
services:
  navigo:
    image: 你的Docker用户名/navigo:latest
    container_name: navigo
    ports:
      - "80:3000"  # 主机80端口映射到容器3000端口
    volumes:
      - ./navigo.db:/app/navigo.db  # 挂载数据库文件到宿主机
    restart: always  # 开机自启
```

#### 2. 启动服务

```shell
docker compose up -d
```

#### 3. 访问系统

打开浏览器访问：http://你的服务器IP

> 默认账号：用户名：team 密码：123456

数据库文件：navigo.db

数据库文件：navigo.db
备份：直接复制该文件即可
升级：替换镜像后，重新启动即可（数据不会丢失）

手动运行（二进制）
从 GitHub Release 下载对应系统的二进制文件，直接运行：

```shell
# Linux
./navigo-linux-amd64

# Windows
navigo-windows-amd64.exe

# macOS
./navigo-darwin-amd64
```


功能说明
游客权限
查看所有分组和链接
点击链接跳转
登录后权限
新增 / 修改 / 删除分组
新增 / 修改 / 删除链接
实时编辑分组名和链接名

📝 开发构建


# 下载依赖
go mod tidy

# 本地运行
go run main.go

# 构建二进制
go build -o navigo main.go
