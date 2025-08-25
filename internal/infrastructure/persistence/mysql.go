package persistence

import (
	"fmt"

	"github.com/base-intern-august-b/clipboard-server/internal/pkg/util"
	"github.com/go-sql-driver/mysql"
)

func MySQL() *mysql.Config {
	c := mysql.NewConfig()

	c.User = util.GetEnvOrDefault("DB_USER", "root")
	c.Passwd = util.GetEnvOrDefault("DB_PASS", "password")
	c.Net = util.GetEnvOrDefault("DB_NET", "tcp")
	c.Addr = fmt.Sprintf(
		"%s:%s",
		util.GetEnvOrDefault("DB_HOST", "localhost"),
		util.GetEnvOrDefault("DB_PORT", "3306"),
	)
	c.DBName = util.GetEnvOrDefault("DB_NAME", "app")
	c.Collation = "utf8mb4_general_ci"
	c.AllowNativePasswords = true
	c.ParseTime = true

	return c
}
