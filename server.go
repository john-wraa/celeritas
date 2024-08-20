package celeritas

import (
	"database/sql"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"os"
	"time"
)

// ListenAndServe starts the web server
func (c *Celeritas) ListenAndServe() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("HOST_INTERFACE"), os.Getenv("PORT")),
		ErrorLog:     c.ErrorLog,
		Handler:      c.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	if c.DB.Pool != nil {
		defer func(Pool *sql.DB) {
			_ = Pool.Close()
		}(c.DB.Pool)
	}

	if redisPool != nil {
		defer func(redisPool *redis.Pool) {
			_ = redisPool.Close()
		}(redisPool)
	}

	if badgerConn != nil {
		defer func(badgerConn *badger.DB) {
			_ = badgerConn.Close()
		}(badgerConn)
	}

	go c.listenRPC()

	c.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	return srv.ListenAndServe()
}
