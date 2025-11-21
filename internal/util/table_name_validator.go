package util

import (
	"fmt"
	"regexp"
)

// starts with letter, then letters/digits/underscore, up to 64 chars total
var tableNameRe = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

// seserved words
var mysqlReserved = map[string]struct{}{
	// core DDL/DML
	"select": {}, "insert": {}, "update": {}, "delete": {}, "from": {}, "where": {}, "into": {}, "values": {},
	"table": {}, "tables": {}, "create": {}, "drop": {}, "alter": {}, "index": {}, "unique": {}, "primary": {},
	"key": {}, "foreign": {}, "references": {}, "constraint": {}, "order": {}, "group": {}, "by": {}, "having": {},
	"join": {}, "left": {}, "right": {}, "inner": {}, "outer": {}, "on": {}, "as": {}, "and": {}, "or": {}, "not": {},
	"union": {}, "all": {}, "distinct": {}, "exists": {}, "like": {}, "in": {}, "is": {}, "null": {}, "case": {},
	// common identifiers that bite
	"user": {}, "role": {}, "grant": {}, "revoke": {}, "schema": {}, "database": {}, "column": {}, "row": {},
}

func ValidateTableName(name string) error {
	if name == "" {
		return fmt.Errorf("table name is empty")
	}
	if len(name) > 64 {
		return fmt.Errorf("table name exceeds 64 characters")
	}
	if !tableNameRe.MatchString(name) {
		return fmt.Errorf("table name must match ^[a-z][a-z0-9_]{0,63}$")
	}
	if _, reserved := mysqlReserved[name]; reserved {
		return fmt.Errorf("table name is a reserved keyword")
	}
	return nil
}
