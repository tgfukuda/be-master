package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // importing with name _ is special import to tell go not to remove this deps
	"github.com/tgfukuda/be-master/util"
)

var testQueries *Queries
var testDB *sql.DB // to use store_test

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("can't load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
