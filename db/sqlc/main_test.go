package db

import (
	"database/sql"
	"github.com/budiharyonoo/simple_bank/utils"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatalln("Error read config file:", err)
		return
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("Error connection to DB:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
