package main

import (
	"github.com/cvilsmeier/go-sqlite-bench/app"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	app.Run(func(dbfile string) app.Db {
		return newDb(dbfile)
	})
}

type dbImpl struct {
	conn *sqlite3.Conn
}

var _ app.Db = (*dbImpl)(nil)

func newDb(dbfile string) app.Db {
	conn, err := sqlite3.Open(dbfile)
	app.MustBeNil(err)
	return &dbImpl{conn}
}

func (d *dbImpl) DriverName() string {
	return "ncruces2"
}

func (d *dbImpl) Exec(sqls ...string) {
	for _, s := range sqls {
		d.exec(s)
	}
}

func (d *dbImpl) exec(sql string) {
	err := d.conn.Exec(sql)
	app.MustBeNil(err)
}

func (d *dbImpl) prepare(sql string) *sqlite3.Stmt {
	stmt, _, err := d.conn.Prepare(sql)
	app.MustBeNil(err)
	return stmt
}

func (d *dbImpl) InsertUsers(insertSql string, users []app.User) {
	tx := d.conn.Begin()
	stmt := d.prepare(insertSql)
	for _, u := range users {
		stmt.BindInt64(1, int64(u.Id))
		stmt.BindInt64(2, app.BindTime(u.Created))
		stmt.BindText(3, u.Email)
		stmt.BindBool(4, u.Active)
		app.MustBe(!stmt.Step())
		app.MustBeNil(stmt.Reset())
	}
	app.MustBeNil(stmt.Close())
	app.MustBeNil(tx.Commit())
}

func (d *dbImpl) InsertArticles(insertSql string, articles []app.Article) {
	tx := d.conn.Begin()
	stmt := d.prepare(insertSql)
	for _, u := range articles {
		stmt.BindInt64(1, int64(u.Id))
		stmt.BindInt64(2, app.BindTime(u.Created))
		stmt.BindInt64(3, int64(u.UserId))
		stmt.BindText(4, u.Text)
		app.MustBe(!stmt.Step())
		app.MustBeNil(stmt.Reset())
	}
	app.MustBeNil(stmt.Close())
	app.MustBeNil(tx.Commit())
}

func (d *dbImpl) InsertComments(insertSql string, comments []app.Comment) {
	tx := d.conn.Begin()
	stmt := d.prepare(insertSql)
	for _, u := range comments {
		stmt.BindInt64(1, int64(u.Id))
		stmt.BindInt64(2, app.BindTime(u.Created))
		stmt.BindInt64(3, int64(u.ArticleId))
		stmt.BindText(4, u.Text)
		app.MustBe(!stmt.Step())
		app.MustBeNil(stmt.Reset())
	}
	app.MustBeNil(stmt.Close())
	app.MustBeNil(tx.Commit())
}

func (d *dbImpl) FindUsers(querySql string) []app.User {
	stmt := d.prepare(querySql)
	var users []app.User
	for stmt.Step() {
		user := app.NewUser(
			stmt.ColumnInt(0),                   // id,
			app.UnbindTime(stmt.ColumnInt64(1)), // created,
			stmt.ColumnText(2),                  // email,
			stmt.ColumnInt(3) != 0,              // active,
		)
		users = append(users, user)
	}
	app.MustBeNil(stmt.Close())
	return users
}

func (d *dbImpl) FindArticles(querySql string) []app.Article {
	stmt := d.prepare(querySql)
	var articles []app.Article
	for stmt.Step() {
		article := app.NewArticle(
			stmt.ColumnInt(0),                   // id,
			app.UnbindTime(stmt.ColumnInt64(1)), // created,
			stmt.ColumnInt(2),                   // userId,
			stmt.ColumnText(3),                  // text,
		)
		articles = append(articles, article)
	}
	app.MustBeNil(stmt.Close())
	return articles
}

func (d *dbImpl) FindUsersArticlesComments(querySql string) ([]app.User, []app.Article, []app.Comment) {
	stmt := d.prepare(querySql)
	// collections
	var users []app.User
	userIndexer := make(map[int]int)
	var articles []app.Article
	articleIndexer := make(map[int]int)
	var comments []app.Comment
	commentIndexer := make(map[int]int)
	for stmt.Step() {
		user := app.NewUser(
			stmt.ColumnInt(0),                   // id,
			app.UnbindTime(stmt.ColumnInt64(1)), // created,
			stmt.ColumnText(2),                  // email,
			stmt.ColumnInt(3) != 0,              // active,
		)
		article := app.NewArticle(
			stmt.ColumnInt(4),                   // id,
			app.UnbindTime(stmt.ColumnInt64(5)), // created,
			stmt.ColumnInt(6),                   // userId,
			stmt.ColumnText(7),                  // text,
		)
		comment := app.NewComment(
			stmt.ColumnInt(8),                   // id,
			app.UnbindTime(stmt.ColumnInt64(9)), // created,
			stmt.ColumnInt(10),                  // articleId,
			stmt.ColumnText(11),                 // text,
		)
		_, ok := userIndexer[user.Id]
		if !ok {
			userIndexer[user.Id] = len(users)
			users = append(users, user)
		}
		_, ok = articleIndexer[article.Id]
		if !ok {
			articleIndexer[article.Id] = len(articles)
			articles = append(articles, article)
		}
		_, ok = commentIndexer[comment.Id]
		if !ok {
			commentIndexer[comment.Id] = len(comments)
			comments = append(comments, comment)
		}
	}
	app.MustBeNil(stmt.Close())
	return users, articles, comments
}

func (d *dbImpl) Close() {
	err := d.conn.Close()
	app.MustBeNil(err)
}
