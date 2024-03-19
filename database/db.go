package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type CodeBaseFolder struct {
	Id        int
	Lang      string
	Domain    string
	Subdomain string
	Path      string
}

func GetConnection() (database *sql.DB) {
	databaseDriver := "mysql"
	databaseUser := "root"
	databaseName := "codebase"
	databasePassword := ""

	database, err := sql.Open(databaseDriver, databaseUser+":"+databasePassword+"@/"+databaseName)
	if err != nil {
		panic(err.Error())
	}

	return database
}

func GetCodeBaseFolders() []CodeBaseFolder {
	var database *sql.DB
	database = GetConnection()
	defer database.Close()

	var error error
	var rows *sql.Rows
	rows, error = database.Query("SELECT * FROM folder ORDER BY id DESC")

	if error != nil {
		panic(error.Error())
	}

	codeBase := CodeBaseFolder{}
	codeBases := []CodeBaseFolder{}

	for rows.Next() {
		error = rows.Scan(&codeBase.Id, &codeBase.Lang, &codeBase.Domain, &codeBase.Subdomain, &codeBase.Path)

		if error != nil {
			panic(error.Error())
		}

		codeBases = append(codeBases, codeBase)
	}

	return codeBases
}

func GetCodeBaseFoldersByLang(lang string) []CodeBaseFolder {
	var database *sql.DB
	database = GetConnection()
	defer database.Close()

	var error error
	var rows *sql.Rows
	rows, error = database.Query("SELECT * FROM folder WHERE lang = ? ORDER BY id DESC")

	if error != nil {
		panic(error.Error())
	}

	codeBase := CodeBaseFolder{}
	codeBases := []CodeBaseFolder{}

	for rows.Next() {
		error = rows.Scan(&codeBase.Id, &codeBase.Lang, &codeBase.Domain, &codeBase.Subdomain, &codeBase.Path)

		if error != nil {
			panic(error.Error())
		}

		codeBases = append(codeBases, codeBase)
	}

	return codeBases
}

func GetCodeBaseFoldersByLangDistinctDomain(lang string) []CodeBaseFolder {
	var database *sql.DB
	database = GetConnection()
	defer database.Close()

	var error error
	var rows *sql.Rows
	rows, error = database.Query("SELECT DISTINCT domain, path FROM folder WHERE lang = ?", lang)

	if error != nil {
		panic(error.Error())
	}

	codeBase := CodeBaseFolder{}
	codeBases := []CodeBaseFolder{}

	for rows.Next() {
		error = rows.Scan(&codeBase.Domain, &codeBase.Path)

		if error != nil {
			panic(error.Error())
		}

		codeBases = append(codeBases, codeBase)
	}

	return codeBases
}

func InsertCodeBaseFolder(codeBase CodeBaseFolder) {
	database := GetConnection()
	defer database.Close()

	var error error
	var insert *sql.Stmt

	insert, error = database.Prepare("INSERT IGNORE INTO folder(lang, domain, path, subdomain) VALUES(?, ?, ?, ?)")

	if error != nil {
		panic(error.Error())
	}
	_, err := insert.Exec(codeBase.Lang, codeBase.Domain, codeBase.Path, codeBase.Subdomain)
	if err != nil {
		panic(err.Error())
	}
}

func UpdateCodeBaseFolder(codeBase CodeBaseFolder) {
	database := GetConnection()
	defer database.Close()

	update, error := database.Prepare("UPDATE folder SET lang=?, domain=?, path=?, subdomain=? WHERE id=?")
	if error != nil {
		panic(error.Error())
	}

	update.Exec(codeBase.Lang, codeBase.Domain, codeBase.Path, codeBase.Subdomain, codeBase.Id)
}

func deleteCodeBaseFolder(codebase CodeBaseFolder) {
	database := GetConnection()
	defer database.Close()

	del, err := database.Prepare("DELETE FROM folder WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	del.Exec(codebase.Id)
}
