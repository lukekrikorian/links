package database

import (
	"links/config"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	db  *sqlx.DB
	err error
)

const linkSchema = `
CREATE TABLE IF NOT EXISTS links (
	id SERIAL,
	url text NOT NULL,
	title text NOT NULL,
	author text,
	comment text,
	type text NOT NULL,
	time timestamp DEFAULT current_timestamp
)
`

const tagSchema = `
CREATE TABLE IF NOT EXISTS tags (
	id integer NOT NULL,
	tag text NOT NULL
)
`

const tagInsert = `
INSERT INTO tags
	(id, tag)
VALUES
	($1, $2)
`

type Link struct {
	URL     string `db:"url"`
	Title   string
	Author  string
	Time    time.Time
	Tags    []string
	Comment string
	Type    string
	ID      int `db:"id"`
}

func (l *Link) GetTags() {
	err = db.Select(&l.Tags, "SELECT UPPER(tag) FROM tags WHERE id = $1", l.ID)
	if err != nil {
		log.Println("Couldn't get tags:", err)
	}
}

func (l *Link) PrettyTags() string {
	return strings.Join(l.Tags, ", ")
}

func (l *Link) PrettyTime() string {
	return l.Time.Format("2 January 2006")
}

func tagsFromString(s string) []string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ToLower(s)
	return strings.Split(s, ",")
}

func insertTags(tags []string, id int) {
	for _, tag := range tags {
		_, err := db.Exec(`
			INSERT INTO tags
				(id, tag)
			VALUES
				($1, $2)
		`, id, tag)
		if err != nil {
			log.Println(err)
		}
	}
}

func GetLinks() []Link {
	var links []Link

	err = db.Select(&links, "SELECT * FROM links")
	if err != nil {
		log.Println("Couldn't get links:", err)
	}

	for i := range links {
		links[i].GetTags()
	}

	return links
}

func init() {
	db, err = sqlx.Connect("pgx", config.Config.DatabaseURL)
	if err != nil {
		log.Fatal("Couldn't connect to database.", err)
	}

	db.MustExec(linkSchema)
	db.MustExec(tagSchema)
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to decode form", http.StatusBadRequest)
		return
	}

	f := flatForm(r.Form)

	var id int
	err := db.QueryRowx(`
		INSERT INTO links 
			(url, title, author, comment, type)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id
	`, f["url"], f["title"], f["author"], f["comment"], f["type"]).Scan(&id)

	if err != nil {
		log.Println(err)
	}

	tags := tagsFromString(f["tags"])
	log.Println("Got tags", tags)
	insertTags(tags, id)

	http.Redirect(w, r, "/upload", http.StatusSeeOther)
}

func flatForm(f url.Values) map[string]string {
	flat := make(map[string]string)
	for key, val := range f {
		flat[key] = val[0]
	}
	return flat
}
