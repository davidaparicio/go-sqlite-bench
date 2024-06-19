package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/cvilsmeier/go-sqlite-bench/app"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	app.Run(func(dbfile string) app.Db {
		db, err := sql.Open("sqlite3", dbfile)
		app.MustBeNil(err)
		return app.NewSqlDb("ncruces", db)
	})
}
