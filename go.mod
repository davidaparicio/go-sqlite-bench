module github.com/cvilsmeier/go-sqlite-bench

go 1.22

require (
	crawshaw.io/sqlite v0.3.2
	github.com/cvilsmeier/sqinn-go v1.2.0
	github.com/eatonphil/gosqlite v0.9.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/ncruces/go-sqlite3 v0.16.3
	modernc.org/sqlite v1.30.1
	zombiezen.com/go/sqlite v1.3.0
)

replace github.com/ncruces/go-sqlite3 => ../go-sqlite3

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/ncruces/julianday v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/tetratelabs/wazero v1.7.3 // indirect
	golang.org/x/sys v0.21.0 // indirect
	modernc.org/gc/v3 v3.0.0-20240304020402-f0dba7c97c2b // indirect
	modernc.org/libc v1.53.3 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/strutil v1.2.0 // indirect
	modernc.org/token v1.1.0 // indirect
)
