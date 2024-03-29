package drivers

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreSql string

func (p PostgreSql) Name() string {
	return string(p)
}

func (PostgreSql) Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB {
	var db *gorm.DB
	var err error
	var dsn strings.Builder
	var logConfig logger.Interface

	dsn.WriteString("host=")
	dsn.WriteString(host)
	dsn.WriteString(" user=")
	dsn.WriteString(user)
	dsn.WriteString(" password=")
	dsn.WriteString(password)
	dsn.WriteString(" dbname=")
	dsn.WriteString(dbname)
	dsn.WriteString(" port=")
	dsn.WriteString(strconv.Itoa(port))
	dsn.WriteString(" sslmode=disable TimeZone=UTC")

	if debug {
		logConfig = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      false,
			},
		)
	} else {
		logConfig = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logger.Warn,
				Colorful:      false,
			},
		)
	}

	db, err = gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logConfig,
	})
	if err != nil {
		log.Printf("Gorm PostgreSQL: %+v \n", err)
		panic(err)
	}

	return db
}
