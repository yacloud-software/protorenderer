package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBSQLMessage
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence sqlmessage_seq;

Main Table:

 CREATE TABLE sqlmessage (id integer primary key default nextval('sqlmessage_seq'),protofile bigint not null  references dbprotofile (id) on delete cascade  ,name text not null  );

Alter statements:
ALTER TABLE sqlmessage ADD COLUMN IF NOT EXISTS protofile bigint not null references dbprotofile (id) on delete cascade  default 0;
ALTER TABLE sqlmessage ADD COLUMN IF NOT EXISTS name text not null default '';


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE sqlmessage_archive (id integer unique not null,protofile bigint not null,name text not null);
*/

import (
	"context"
	gosql "database/sql"
	"fmt"
	savepb "golang.conradwood.net/apis/protorenderer2"
	"golang.conradwood.net/go-easyops/sql"
	"os"
)

var (
	default_def_DBSQLMessage *DBSQLMessage
)

type DBSQLMessage struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
}

func DefaultDBSQLMessage() *DBSQLMessage {
	if default_def_DBSQLMessage != nil {
		return default_def_DBSQLMessage
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBSQLMessage(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBSQLMessage = res
	return res
}
func NewDBSQLMessage(db *sql.DB) *DBSQLMessage {
	foo := DBSQLMessage{DB: db}
	foo.SQLTablename = "sqlmessage"
	foo.SQLArchivetablename = "sqlmessage_archive"
	return &foo
}

// archive. It is NOT transactionally save.
func (a *DBSQLMessage) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBSQLMessage", "insert into "+a.SQLArchivetablename+" (id,protofile, name) values ($1,$2, $3) ", p.ID, p.ProtoFile.ID, p.Name)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBSQLMessage) Save(ctx context.Context, p *savepb.SQLMessage) (uint64, error) {
	qn := "DBSQLMessage_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (protofile, name) values ($1, $2) returning id", a.get_ProtoFile_ID(p), a.get_Name(p))
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
func (a *DBSQLMessage) SaveWithID(ctx context.Context, p *savepb.SQLMessage) error {
	qn := "insert_DBSQLMessage"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,protofile, name) values ($1,$2, $3) ", p.ID, p.ProtoFile.ID, p.Name)
	return a.Error(ctx, qn, e)
}

func (a *DBSQLMessage) Update(ctx context.Context, p *savepb.SQLMessage) error {
	qn := "DBSQLMessage_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set protofile=$1, name=$2 where id = $3", a.get_ProtoFile_ID(p), a.get_Name(p), p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBSQLMessage) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBSQLMessage_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBSQLMessage) ByID(ctx context.Context, p uint64) (*savepb.SQLMessage, error) {
	qn := "DBSQLMessage_ByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,protofile, name from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No SQLMessage with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) SQLMessage with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBSQLMessage) TryByID(ctx context.Context, p uint64) (*savepb.SQLMessage, error) {
	qn := "DBSQLMessage_TryByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,protofile, name from "+a.SQLTablename+" where id = $1", p)
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
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) SQLMessage with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBSQLMessage) All(ctx context.Context) ([]*savepb.SQLMessage, error) {
	qn := "DBSQLMessage_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,protofile, name from "+a.SQLTablename+" order by id")
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

// get all "DBSQLMessage" rows with matching ProtoFile
func (a *DBSQLMessage) ByProtoFile(ctx context.Context, p uint64) ([]*savepb.SQLMessage, error) {
	qn := "DBSQLMessage_ByProtoFile"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,protofile, name from "+a.SQLTablename+" where protofile = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByProtoFile: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByProtoFile: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBSQLMessage) ByLikeProtoFile(ctx context.Context, p uint64) ([]*savepb.SQLMessage, error) {
	qn := "DBSQLMessage_ByLikeProtoFile"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,protofile, name from "+a.SQLTablename+" where protofile ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByProtoFile: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByProtoFile: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBSQLMessage" rows with matching Name
func (a *DBSQLMessage) ByName(ctx context.Context, p string) ([]*savepb.SQLMessage, error) {
	qn := "DBSQLMessage_ByName"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,protofile, name from "+a.SQLTablename+" where name = $1", p)
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
func (a *DBSQLMessage) ByLikeName(ctx context.Context, p string) ([]*savepb.SQLMessage, error) {
	qn := "DBSQLMessage_ByLikeName"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,protofile, name from "+a.SQLTablename+" where name ilike $1", p)
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

/**********************************************************************
* The field getters
**********************************************************************/

// getter for field "ID" (ID) [uint64]
func (a *DBSQLMessage) get_ID(p *savepb.SQLMessage) uint64 {
	return uint64(p.ID)
}

// getter for reference "ProtoFile"
func (a *DBSQLMessage) get_ProtoFile_ID(p *savepb.SQLMessage) uint64 {
	if p.ProtoFile == nil {
		panic("field ProtoFile must not be nil")
	}
	return p.ProtoFile.ID
}

// getter for field "Name" (Name) [string]
func (a *DBSQLMessage) get_Name(p *savepb.SQLMessage) string {
	return string(p.Name)
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBSQLMessage) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.SQLMessage, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBSQLMessage) Tablename() string {
	return a.SQLTablename
}

func (a *DBSQLMessage) SelectCols() string {
	return "id,protofile, name"
}
func (a *DBSQLMessage) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".protofile, " + a.SQLTablename + ".name"
}

func (a *DBSQLMessage) FromRowsOld(ctx context.Context, rows *gosql.Rows) ([]*savepb.SQLMessage, error) {
	var res []*savepb.SQLMessage
	for rows.Next() {
		foo := savepb.SQLMessage{ProtoFile: &savepb.DBProtoFile{}}
		err := rows.Scan(&foo.ID, &foo.ProtoFile.ID, &foo.Name)
		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, &foo)
	}
	return res, nil
}
func (a *DBSQLMessage) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.SQLMessage, error) {
	var res []*savepb.SQLMessage
	for rows.Next() {
		// SCANNER:
		foo := &savepb.SQLMessage{}
		// create the non-nullable pointers
		foo.ProtoFile = &savepb.DBProtoFile{} // non-nullable
		// create variables for scan results
		scanTarget_0 := &foo.ID
		scanTarget_1 := &foo.ProtoFile.ID
		scanTarget_2 := &foo.Name
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
func (a *DBSQLMessage) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),protofile bigint not null ,name text not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),protofile bigint not null ,name text not null );`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS protofile bigint not null default 0;`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS name text not null default '';`,

		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS protofile bigint not null  default 0;`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS name text not null  default '';`,
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
		`ALTER TABLE ` + a.SQLTablename + ` add constraint mkdb_fk_sqlmessage_protofile_dbprotofileid FOREIGN KEY (protofile) references dbprotofile (id) on delete cascade ;`,
	}
	for i, c := range csql {
		a.DB.ExecContextQuiet(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
	}
	return nil
}

/**********************************************************************
* Helper to meaningful errors
**********************************************************************/
func (a *DBSQLMessage) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}

