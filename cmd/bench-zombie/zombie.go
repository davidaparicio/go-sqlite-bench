package main

import (
	"github.com/cvilsmeier/go-sqlite-bench/app"
	"zombiezen.com/go/sqlite"
)

func main() {
	app.Run(func(dbfile string) app.Db {
		return newDb(dbfile)
	})
}

type dbImpl struct {
	conn *sqlite.Conn
}

var _ app.Db = (*dbImpl)(nil)

func newDb(dbfile string) app.Db {
	conn, err := sqlite.OpenConn(dbfile, sqlite.OpenReadWrite, sqlite.OpenCreate)
	app.MustBeNil(err)
	return &dbImpl{conn}
}

func (d *dbImpl) Exec(sqls ...string) {
	for _, s := range sqls {
		d.exec(s)
	}
}

func (d *dbImpl) InsertUsers(insertSql string, users []app.User) {
	d.exec("BEGIN")
	stmt := d.conn.Prep(insertSql)
	for _, u := range users {
		//	Id        int
		//	Created   time.Time
		//	Email     string
		//	Active    bool
		stmt.BindInt64(1, int64(u.Id))
		stmt.BindInt64(2, app.BindTime(u.Created))
		stmt.BindText(3, u.Email)
		stmt.BindBool(4, u.Active)
		_, err := stmt.Step()
		app.MustBeNil(err)
		err = stmt.Reset()
		app.MustBeNil(err)
	}
	err := stmt.Finalize()
	app.MustBeNil(err)
	d.exec("COMMIT")
}

func (d *dbImpl) InsertArticles(insertSql string, articles []app.Article) {
	d.exec("BEGIN")
	stmt := d.conn.Prep(insertSql)
	for _, u := range articles {
		stmt.BindInt64(1, int64(u.Id))
		stmt.BindInt64(2, app.BindTime(u.Created))
		stmt.BindInt64(3, int64(u.UserId))
		stmt.BindText(4, u.Text)
		_, err := stmt.Step()
		app.MustBeNil(err)
		err = stmt.Reset()
		app.MustBeNil(err)
	}
	err := stmt.Finalize()
	app.MustBeNil(err)
	d.exec("COMMIT")
}

func (d *dbImpl) InsertComments(insertSql string, comments []app.Comment) {
	d.exec("BEGIN")
	stmt := d.conn.Prep(insertSql)
	for _, u := range comments {
		stmt.BindInt64(1, int64(u.Id))
		stmt.BindInt64(2, app.BindTime(u.Created))
		stmt.BindInt64(3, int64(u.ArticleId))
		stmt.BindText(4, u.Text)
		_, err := stmt.Step()
		app.MustBeNil(err)
		err = stmt.Reset()
		app.MustBeNil(err)
	}
	err := stmt.Finalize()
	app.MustBeNil(err)
	d.exec("COMMIT")
}

func (d *dbImpl) FindUsers(querySql string) []app.User {
	stmt, err := d.conn.Prepare(querySql)
	app.MustBeNil(err)
	more, err := stmt.Step()
	app.MustBeNil(err)
	var users []app.User
	for more {
		user := app.NewUser(
			stmt.ColumnInt(0),                   // id,
			app.UnbindTime(stmt.ColumnInt64(1)), // created,
			stmt.ColumnText(2),                  // email,
			stmt.ColumnInt(3) != 0,              // active,
		)
		users = append(users, user)
		more, err = stmt.Step()
		app.MustBeNil(err)
	}
	return users
}

func (d *dbImpl) FindArticles(querySql string) []app.Article {
	stmt, err := d.conn.Prepare(querySql)
	app.MustBeNil(err)
	more, err := stmt.Step()
	app.MustBeNil(err)
	var articles []app.Article
	for more {
		article := app.NewArticle(
			stmt.ColumnInt(0),                   // id,
			app.UnbindTime(stmt.ColumnInt64(1)), // created,
			stmt.ColumnInt(2),                   // userId,
			stmt.ColumnText(3),                  // text,
		)
		articles = append(articles, article)
		more, err = stmt.Step()
		app.MustBeNil(err)
	}
	return articles
}

func (d *dbImpl) FindUsersArticlesComments(querySql string) ([]app.User, []app.Article, []app.Comment) {
	stmt, err := d.conn.Prepare(querySql)
	app.MustBeNil(err)
	more, err := stmt.Step()
	app.MustBeNil(err)
	// collections
	var users []app.User
	userIndexer := make(map[int]int)
	var articles []app.Article
	articleIndexer := make(map[int]int)
	var comments []app.Comment
	commentIndexer := make(map[int]int)
	for more {
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
		more, err = stmt.Step()
		app.MustBeNil(err)
	}
	return users, articles, comments
}

func (d *dbImpl) Close() {
	err := d.conn.Close()
	app.MustBeNil(err)
}

func (d *dbImpl) exec(sql string) {
	stmt := d.conn.Prep(sql)
	app.MustBeSet(stmt)
	_, err := stmt.Step()
	app.MustBeNil(err)
	err = stmt.Finalize()
	app.MustBeNil(err)
}
