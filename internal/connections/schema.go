// internal/connections/schema.go
// Schema discovery — the IT-side analog of internal/discovery.BrowseNodeTree.
// See docs/analysis_log.md Entry 115 for the full Track B plan this is
// Phase 1 of.
package connections

import (
	"context"
	"database/sql"
)

// ColumnSchema is one discovered column.
type ColumnSchema struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	IsKey    bool   `json:"is_key"`
}

// TableSchema is one discovered table and its columns, in ordinal order.
type TableSchema struct {
	Name    string         `json:"name"`
	Columns []ColumnSchema `json:"columns"`
}

// DiscoverSchema browses the connection's own database via information_schema
// — MySQL-only, matching the sql_query connector's current V1a scope. Returns
// tables in the order MySQL's information_schema naturally groups them
// (alphabetical by table_name), columns in ordinal (declaration) order.
func DiscoverSchema(db *sql.DB) ([]TableSchema, error) {
	rows, err := db.QueryContext(context.Background(), `
		SELECT table_name, column_name, data_type, column_key
		FROM information_schema.columns
		WHERE table_schema = DATABASE()
		ORDER BY table_name, ordinal_position`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order []string
	byTable := map[string]*TableSchema{}
	for rows.Next() {
		var table, column, dataType, key string
		if err := rows.Scan(&table, &column, &dataType, &key); err != nil {
			return nil, err
		}
		t, ok := byTable[table]
		if !ok {
			t = &TableSchema{Name: table}
			byTable[table] = t
			order = append(order, table)
		}
		t.Columns = append(t.Columns, ColumnSchema{Name: column, DataType: dataType, IsKey: key == "PRI"})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]TableSchema, 0, len(order))
	for _, name := range order {
		out = append(out, *byTable[name])
	}
	return out, nil
}

// DatabaseSchema is one database and every table visible in it.
type DatabaseSchema struct {
	Name   string        `json:"name"`
	Tables []TableSchema `json:"tables"`
}

// ListDatabasesAndTables browses every database and table visible to this
// connection's user in one pass — the "connect and see everything" browse
// experience. Deliberately separate from DiscoverSchema: that function stays
// scoped to the connection's own configured database, since that's what the
// canonical-mapping heuristic and SchemaMapping KG nodes (Entry 115/116)
// assume — one connection maps to one database's worth of structure. This
// function is read-only visibility, nothing it returns feeds the KG. System
// schemas (mysql/information_schema/performance_schema/sys) are excluded —
// never useful to a customer browsing their own data.
func ListDatabasesAndTables(db *sql.DB) ([]DatabaseSchema, error) {
	rows, err := db.QueryContext(context.Background(), `
		SELECT table_schema, table_name, column_name, data_type, column_key
		FROM information_schema.columns
		WHERE table_schema NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
		ORDER BY table_schema, table_name, ordinal_position`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbOrder []string
	byDB := map[string]*DatabaseSchema{}
	tableOrder := map[string][]string{}
	byTable := map[string]map[string]*TableSchema{}

	for rows.Next() {
		var schema, table, column, dataType, key string
		if err := rows.Scan(&schema, &table, &column, &dataType, &key); err != nil {
			return nil, err
		}
		d, ok := byDB[schema]
		if !ok {
			d = &DatabaseSchema{Name: schema}
			byDB[schema] = d
			byTable[schema] = map[string]*TableSchema{}
			dbOrder = append(dbOrder, schema)
		}
		t, ok := byTable[schema][table]
		if !ok {
			t = &TableSchema{Name: table}
			byTable[schema][table] = t
			tableOrder[schema] = append(tableOrder[schema], table)
		}
		t.Columns = append(t.Columns, ColumnSchema{Name: column, DataType: dataType, IsKey: key == "PRI"})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]DatabaseSchema, 0, len(dbOrder))
	for _, dbName := range dbOrder {
		d := byDB[dbName]
		for _, tableName := range tableOrder[dbName] {
			d.Tables = append(d.Tables, *byTable[dbName][tableName])
		}
		out = append(out, *d)
	}
	return out, nil
}
