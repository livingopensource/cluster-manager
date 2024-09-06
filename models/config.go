package models

import (
	"context"
	"fmt"
	"time"

	"github.com/go-gorm/caches/v4"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type redisCacher struct {
	rdb *redis.Client
}

func (c *redisCacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
	res, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := q.Unmarshal([]byte(res)); err != nil {
		return nil, err
	}

	return q, nil
}

func (c *redisCacher) Store(ctx context.Context, key string, val *caches.Query[any]) error {
	res, err := val.Marshal()
	if err != nil {
		return err
	}

	c.rdb.Set(ctx, key, res, viper.GetDuration("database.cachettl")*time.Second)
	return nil
}

func (c *redisCacher) Invalidate(ctx context.Context) error {
	var (
		cursor uint64
		keys   []string
	)
	for {
		var (
			k   []string
			err error
		)
		k, cursor, err = c.rdb.Scan(ctx, cursor, fmt.Sprintf("%s*", caches.IdentifierPrefix), 0).Result()
		if err != nil {
			return err
		}
		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		if _, err := c.rdb.Del(ctx, keys...).Result(); err != nil {
			return err
		}
	}
	return nil
}

// Database database instance modelling
type database struct {
	DBType     string
	DBName     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBHost     string
}

// MySQLDBConfig initialises a MySQL/MariaDB database instance. Return gorm instance and error
func (db *database) mySQLDBConfig() (*gorm.DB, error) {
	dbConfig := db.DBUser + ":" + db.DBPassword + "@tcp(" + db.DBHost + ":" + db.DBPort + ")/" + db.DBName + "?parseTime=true"
	connection, err := gorm.Open(mysql.Open(dbConfig), &gorm.Config{})
	cachesPlugin := &caches.Caches{Conf: &caches.Config{
		Easer: true,
		Cacher: &redisCacher{
			rdb: redis.NewClient(&redis.Options{
				Addr:     viper.GetString("redis.url"),
				Password: viper.GetString("redis.password"),
				DB:       viper.GetInt("redis.db"),
			}),
		},
	}}
	connection.Use(cachesPlugin)
	return connection, err
}

// PostgreSQLConfig initialises a PostgreSQL db instances, returns gorm, error
func (db *database) postgreSQLConfig() (*gorm.DB, error) {
	dbConfig := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		db.DBUser,
		db.DBPassword,
		db.DBHost,
		db.DBPort,
		db.DBName)

	connection, err := gorm.Open(postgres.Open(dbConfig), &gorm.Config{})
	cachesPlugin := &caches.Caches{Conf: &caches.Config{
		Easer: true,
		Cacher: &redisCacher{
			rdb: redis.NewClient(&redis.Options{
				Addr:     viper.GetString("redis.url"),
				Password: viper.GetString("redis.password"),
				DB:       viper.GetInt("redis.db"),
			}),
		},
	}}
	connection.Use(cachesPlugin)
	return connection, err
}

// dbConfig initialises the required database based on the selected DB
// postgres is selected by default
func (db *database) dbConfig() (*gorm.DB, error) {
	switch db.DBType {
	case "mysql":
		return db.mySQLDBConfig()
	case "postgres":
		return db.postgreSQLConfig()
	default:
		return db.postgreSQLConfig()
	}
}
