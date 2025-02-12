package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

func TestMySQLConn(t *testing.T) {
	type dbconf struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Db       string `yaml:"db"`
		Table    string `yaml:"table"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	// 读取 YAML 配置文件
	file, err := os.Open("../conf/testdb.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var cfg dbconf
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		t.Fatal(err)
	}

	// 数据库连接信息
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Db,
	)
	// 打开数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	// 执行查询
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", cfg.Table))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var id int
		var account string
		var username string
		var password string
		// 根据你的表的结构调整以下变量
		if err := rows.Scan(&id, &account, &username, &password); err != nil {
			t.Fatal(err)
		}
		t.Log("query result:", id, account, username, password)
	}

	// 检查遍历期间是否出现错误
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}
