package data

import (
	"path/filepath"
	"time"

	"answer/pkg/dir"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/segmentfault/pacman/cache"
	"github.com/segmentfault/pacman/contrib/cache/memory"
	"github.com/segmentfault/pacman/log"
	"xorm.io/core"
	"xorm.io/xorm"
	ormlog "xorm.io/xorm/log"
	"xorm.io/xorm/schemas"
)

// Data data
type Data struct {
	DB    *xorm.Engine
	Cache cache.Cache
}

// NewData new data instance
func NewData(db *xorm.Engine, cache cache.Cache) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
		db.Close()
	}
	return &Data{DB: db, Cache: cache}, cleanup, nil
}

// NewDB new database instance
func NewDB(debug bool, dataConf *Database) (*xorm.Engine, error) {
	if dataConf.Driver == "" {
		dataConf.Driver = string(schemas.MYSQL)
	}
	if dataConf.Driver == string(schemas.SQLITE) {
		dbFileDir := filepath.Dir(dataConf.Connection)
		log.Debugf("try to create database directory %s", dbFileDir)
		if err := dir.CreateDirIfNotExist(dbFileDir); err != nil {
			log.Errorf("create database dir failed: %s", err)
		}
	}
	engine, err := xorm.NewEngine(dataConf.Driver, dataConf.Connection)
	if err != nil {
		return nil, err
	}

	if debug {
		engine.ShowSQL(true)
	} else {
		engine.SetLogLevel(ormlog.LOG_ERR)
	}

	if err = engine.Ping(); err != nil {
		return nil, err
	}

	if dataConf.MaxIdleConn > 0 {
		engine.SetMaxIdleConns(dataConf.MaxIdleConn)
	}
	if dataConf.MaxOpenConn > 0 {
		engine.SetMaxOpenConns(dataConf.MaxOpenConn)
	}
	if dataConf.ConnMaxLifeTime > 0 {
		engine.SetConnMaxLifetime(time.Duration(dataConf.ConnMaxLifeTime) * time.Second)
	}
	engine.SetColumnMapper(core.GonicMapper{})
	return engine, nil
}

// NewCache new cache instance
func NewCache(c *CacheConf) (cache.Cache, func(), error) {
	// TODO What cache type should be initialized according to the configuration file
	memCache := memory.NewCache()

	if len(c.FilePath) > 0 {
		cacheFileDir := filepath.Dir(c.FilePath)
		log.Debugf("try to create cache directory %s", cacheFileDir)
		err := dir.CreateDirIfNotExist(cacheFileDir)
		if err != nil {
			log.Errorf("create cache dir failed: %s", err)
		}
		log.Infof("try to load cache file from %s", c.FilePath)
		if err := memory.Load(memCache, c.FilePath); err != nil {
			log.Warn(err)
		}
		go func() {
			ticker := time.Tick(time.Minute)
			for range ticker {
				if err := memory.Save(memCache, c.FilePath); err != nil {
					log.Warn(err)
				}
			}
		}()
	}
	cleanup := func() {
		log.Infof("try to save cache file to %s", c.FilePath)
		if err := memory.Save(memCache, c.FilePath); err != nil {
			log.Warn(err)
		}
	}
	return memCache, cleanup, nil
}
