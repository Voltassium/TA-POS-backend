package database

import (
	"backend-ta/pkg/logger"
	"time"

	"github.com/uptrace/bun"
)

type (
	QueryHook struct {
		bun.QueryHook

		logger       *logger.ZapLogger
		slowDuration time.Duration
	}

	Database struct {
		*bun.DB
		Query *bun.SelectQuery
	}
)
