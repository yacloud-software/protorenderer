package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBDBProtoFile
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence dbprotofile_seq;

Main Table:

 CREATE TABLE dbprotofile (id integer primary key default nextval('dbprotofile_seq'),name text not null  ,repositoryid bigint not null  );

Alter statements:
ALTER TABLE dbprotofile ADD COLUMN IF NOT EXISTS name text not null default '';
ALTER TABLE dbprotofile ADD COLUMN IF NOT EXISTS repositoryid bigint not null default 0;


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE dbprotofile_archive (id integer unique not null,name text not null,repositoryid bigint not null);
*/

import (
	"context"
	gosql "database/sql"
	"fmt"
	savepb "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/sql"
	"os"
)

var (
	default_def_DBDBProtoFile *DBDBProtoFile
)

type DBDBProtoFile struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
}

func DefaultDBDBProtoFile() *DBDBProtoFile {
	if default_def_DBDBProtoFile != nil {
		return default_def_DBDBProtoFile
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBDBProtoFile(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBDBProtoFile = res
	return res
}
func NewDBDBProtoFile(db *sql.DB) *DBDBProtoFile {
	foo := DBDBProtoFile{DB: db}
	foo.SQLTablename = "dbprotofile"
	foo.SQLArchivetablename = "dbprotofile_archive"
	return &foo
}

// archive. It is NOT transactionally save.
func (a *DBDBProtoFile) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBDBProtoFile", "insert into "+a.SQLArchivetablename+" (id,name, repositoryid) values ($1,$2, $3) ", p.ID, p.Name, p.RepositoryID)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBDBProtoFile) Save(ctx context.Context, p *savepb.DBProtoFile) (uint64, error) {
	qn := "DBDBProtoFile_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (name, repositoryid) values ($1, $2) returning id", a.get_Name(p), a.get_RepositoryID(p))
	if e != nil {
		return 0, a.Error(ctx, qn, e)
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, a.Error(ctx, qn, fmt.Errorf("No rows after insert"))
	}
	var id uint64
	e = rows.Scan(&id)
	if e != nil {
		return 0, a.Error(ctx, qn, fmt.Errorf("failed to scan id after insert: %s", e))
	}
	p.ID = id
	return id, nil
}

// Save using the ID specified
func (a *DBDBProtoFile) SaveWithID(ctx context.Context, p *savepb.DBProtoFile) error {
	qn := "insert_DBDBProtoFile"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,name, repositoryid) values ($1,$2, $3) ", p.ID, p.Name, p.RepositoryID)
	return a.Error(ctx, qn, e)
}

func (a *DBDBProtoFile) Update(ctx context.Context, p *savepb.DBProtoFile) error {
	qn := "DBDBProtoFile_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set name=$1, repositoryid=$2 where id = $3", a.get_Name(p), a.get_RepositoryID(p), p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBDBProtoFile) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBDBProtoFile_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBDBProtoFile) ByID(ctx context.Context, p uint64) (*savepb.DBProtoFile, error) {
	qn := "DBDBProtoFile_ByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,name, repositoryid from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No DBProtoFile with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) DBProtoFile with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBDBProtoFile) TryByID(ctx context.Context, p uint64) (*savepb.DBProtoFile, error) {
	qn := "DBDBProtoFile_TryByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,name, repositoryid from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("TryByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) DBProtoFile with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBDBProtoFile) All(ctx context.Context) ([]*savepb.DBProtoFile, error) {
	qn := "DBDBProtoFile_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,name, repositoryid from "+a.SQLTablename+" order by id")
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("All: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBDBProtoFile" rows with matching Name
func (a *DBDBProtoFile) ByName(ctx context.Context, p string) ([]*savepb.DBProtoFile, error) {
	qn := "DBDBProtoFile_ByName"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,name, repositoryid from "+a.SQLTablename+" where name = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByName: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByName: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBDBProtoFile) ByLikeName(ctx context.Context, p string) ([]*savepb.DBProtoFile, error) {
	qn := "DBDBProtoFile_ByLikeName"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,name, repositoryid from "+a.SQLTablename+" where name ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByName: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByName: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBDBProtoFile" rows with matching RepositoryID
func (a *DBDBProtoFile) ByRepositoryID(ctx context.Context, p uint64) ([]*savepb.DBProtoFile, error) {
	qn := "DBDBProtoFile_ByRepositoryID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,name, repositoryid from "+a.SQLTablename+" where repositoryid = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByRepositoryID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByRepositoryID: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBDBProtoFile) ByLikeRepositoryID(ctx context.Context, p uint64) ([]*savepb.DBProtoFile, error) {
	qn := "DBDBProtoFile_ByLikeRepositoryID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,name, repositoryid from "+a.SQLTablename+" where repositoryid ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByRepositoryID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByRepositoryID: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* The field getters
**********************************************************************/

// getter for field "ID" (ID) [uint64]
func (a *DBDBProtoFile) get_ID(p *savepb.DBProtoFile) uint64 {
	return uint64(p.ID)
}

// getter for field "Name" (Name) [string]
func (a *DBDBProtoFile) get_Name(p *savepb.DBProtoFile) string {
	return string(p.Name)
}

// getter for field "RepositoryID" (RepositoryID) [uint64]
func (a *DBDBProtoFile) get_RepositoryID(p *savepb.DBProtoFile) uint64 {
	return uint64(p.RepositoryID)
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBDBProtoFile) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.DBProtoFile, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBDBProtoFile) Tablename() string {
	return a.SQLTablename
}

func (a *DBDBProtoFile) SelectCols() string {
	return "id,name, repositoryid"
}
func (a *DBDBProtoFile) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".name, " + a.SQLTablename + ".repositoryid"
}

func (a *DBDBProtoFile) FromRowsOld(ctx context.Context, rows *gosql.Rows) ([]*savepb.DBProtoFile, error) {
	var res []*savepb.DBProtoFile
	for rows.Next() {
		foo := savepb.DBProtoFile{}
		err := rows.Scan(&foo.ID, &foo.Name, &foo.RepositoryID)
		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, &foo)
	}
	return res, nil
}
func (a *DBDBProtoFile) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.DBProtoFile, error) {
	var res []*savepb.DBProtoFile
	for rows.Next() {
		// SCANNER:
		foo := &savepb.DBProtoFile{}
		// create the non-nullable pointers
		// create variables for scan results
		scanTarget_0 := &foo.ID
		scanTarget_1 := &foo.Name
		scanTarget_2 := &foo.RepositoryID
		err := rows.Scan(scanTarget_0, scanTarget_1, scanTarget_2)
		// END SCANNER

		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, foo)
	}
	return res, nil
}

/**********************************************************************
* Helper to create table and columns
**********************************************************************/
func (a *DBDBProtoFile) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),name text not null ,repositoryid bigint not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),name text not null ,repositoryid bigint not null );`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS name text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS repositoryid bigint not null default 0;`,

		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS name text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS repositoryid bigint not null  default 0;`,
	}

	for i, c := range csql {
		_, e := a.DB.ExecContext(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
		if e != nil {
			return e
		}
	}

	// these are optional, expected to fail
	csql = []string{
		// Indices:

		// Foreign keys:

	}
	for i, c := range csql {
		a.DB.ExecContextQuiet(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
	}
	return nil
}

/**********************************************************************
* Helper to meaningful errors
**********************************************************************/
func (a *DBDBProtoFile) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}

