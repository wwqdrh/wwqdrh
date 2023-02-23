# 简介

- migrate: `https://github.com/golang-migrate/migrate`, 提供命令行模式进行管理
- gormmigrate: `https://github.com/go-gormigrate/gormigrate`，与gorm进行适配

## migrate

安装migrate，下载链接`https://github.com/golang-migrate/migrate/releases/tag/v4.15.2`

```bash
cd handle/dbmigrate/

mkdir migrations #首次执行，用于存放变更的 sql 文件

migrate create -ext sql -dir ./migrations -seq create_users_table

# 分别在生成的sql里面写装载以及回退的逻辑

# 首次执行启动逻辑
migrate -database sqlite3://test.db -path ./migrations up

migrate -database 'postgres://postgres:mysecretpassword@192.168.10.212:5432/example?sslmode=disable' -path ./migrations up 1
```

## gormigrate

具体代码查看同目录下的main.go文件


在gomigrate中，创建了migrate之后，可以自己通过调用方法去控制数据库版本

```bash
func (g *Gormigrate) InitSchema(initSchema InitSchemaFunc)
func (g *Gormigrate) Migrate() error
func (g *Gormigrate) MigrateTo(migrationID string) error
func (g *Gormigrate) RollbackLast() error
func (g *Gormigrate) RollbackMigration(m *Migration) error
func (g *Gormigrate) RollbackTo(migrationID string) error
```
