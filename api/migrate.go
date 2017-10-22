package api

import (
	"database/sql"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/database"
	"os"
)

func Migration(mConn string) {

	log.Info("Run migration")

	db, err := sql.Open("mysql", mConn)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// OR: Read migrations from a folder:
	//migrations := &migrate.FileMigrationSource{
	//	Dir: "database/migrations",
	//}

	// OR: Use migrations from bindata:
	migrations := &migrate.AssetMigrationSource{
		Asset:    database.Asset,
		AssetDir: database.AssetDir,
		Dir:      "database/migrations",
	}

	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Infof("Applied %d migrations!", n)
}