//go:build integration
// +build integration

package tests

import (
	"time"

	"gitlab.ozon.dev/pircuser61/catalog/tests/postgres"
)

var (
	Db      *postgres.TestDB
	Timeout time.Duration
)

func init() {
	Db = postgres.NewFromEnv()
	Timeout = time.Second * 2
}
