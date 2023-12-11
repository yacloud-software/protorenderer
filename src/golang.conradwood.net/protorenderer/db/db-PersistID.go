package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBPersistID
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence persistid_seq;

Main Table:

 CREATE TABLE persistid (id integer primary key default nextval('persistid_seq'),key text not null  );

Alter statements:
ALTER TABLE persistid ADD COLUMN key text not null default '';


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE persistid_archive (id integer unique not null,key text not null);
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
	default_def_DBPersistID *DBPersistID
)

type DBPersistID struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
}

func DefaultDBPersistID() *DBPersistID {
	if default_def_DBPersistID != nil {
		return default_def_DBPersistID
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBPersistID(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBPersistID = res
	return res
}
func NewDBPersistID(db *sql.DB) *DBPersistID {
	foo := DBPersistID{DB: db}
	foo.SQLTablename = "persistid"
	foo.SQLArchivetablename = "persistid_archive"
	return &foo
}

// archive. It is NOT transactionally save.
func (a *DBPersistID) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBPersistID", "insert into "+a.SQLArchivetablename+" (id,key) values ($1,$2) ", p.ID, p.Key)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBPersistID) Save(ctx context.Context, p *savepb.PersistID) (uint64, error) {
	qn := "DBPersistID_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (key) values ($1) returning id", p.Key)
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
func (a *DBPersistID) SaveWithID(ctx context.Context, p *savepb.PersistID) error {
	qn := "insert_DBPersistID"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,key) values ($1,$2) ", p.ID, p.Key)
	return a.Error(ctx, qn, e)
}

func (a *DBPersistID) Update(ctx context.Context, p *savepb.PersistID) error {
	qn := "DBPersistID_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set key=$1 where id = $2", p.Key, p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBPersistID) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBPersistID_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBPersistID) ByID(ctx context.Context, p uint64) (*savepb.PersistID, error) {
	qn := "DBPersistID_ByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,key from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No PersistID with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) PersistID with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBPersistID) All(ctx context.Context) ([]*savepb.PersistID, error) {
	qn := "DBPersistID_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,key from "+a.SQLTablename+" order by id")
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

// get all "DBPersistID" rows with matching Key
func (a *DBPersistID) ByKey(ctx context.Context, p string) ([]*savepb.PersistID, error) {
	qn := "DBPersistID_ByKey"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,key from "+a.SQLTablename+" where key = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByKey: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByKey: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBPersistID) ByLikeKey(ctx context.Context, p string) ([]*savepb.PersistID, error) {
	qn := "DBPersistID_ByLikeKey"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,key from "+a.SQLTablename+" where key ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByKey: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByKey: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBPersistID) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.PersistID, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBPersistID) Tablename() string {
	return a.SQLTablename
}

func (a *DBPersistID) SelectCols() string {
	return "id,key"
}
func (a *DBPersistID) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".key"
}

func (a *DBPersistID) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.PersistID, error) {
	var res []*savepb.PersistID
	for rows.Next() {
		foo := savepb.PersistID{}
		err := rows.Scan(&foo.ID, &foo.Key)
		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, &foo)
	}
	return res, nil
}

/**********************************************************************
* Helper to create table and columns
**********************************************************************/
func (a *DBPersistID) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),key text not null  );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),key text not null  );`,
	}
	for i, c := range csql {
		_, e := a.DB.ExecContext(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
		if e != nil {
			return e
		}
	}
	return nil
}

/**********************************************************************
* Helper to meaningful errors
**********************************************************************/
func (a *DBPersistID) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}























































































