package component

import (
	"fmt"
	"github.com/allegro/bigcache"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"time"
)

var (
	DB          *gorm.DB
	GlobalCache *bigcache.BigCache
)

func init() {
	// Connect to DB
	//var err error
	//DB, err = gorm.Open("mysql", "your_db_url")
	//if err != nil {
	//	panic(fmt.Errorf("failed to connect to DB: %w", err))
	//}

	database, err := gorm.Open(sqlite.Open("test2.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database

	// Initialize cache
	var err2 error
	GlobalCache, err2 = bigcache.NewBigCache(bigcache.DefaultConfig(30 * time.Minute))
	if err2 != nil {
		panic(fmt.Errorf("failed to initialize cahce: %w", err))
	}
}
