package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"errors"
)

func BatchInsert(objArr []interface{}) (int64, error) {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return 0, errors.New("slice must not be empty")
	}

	mainObj := objArr[0]
	mainScope := db.NewScope(mainObj)
	mainFields := mainScope.Fields()
	quoted := make([]string, 0, len(mainFields))
	for i := range mainFields {
		// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
		// If field is ignore field, skip it.
		if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) || mainFields[i].Relationship != nil {
			continue
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}

	placeholdersArr := make([]string, 0, len(objArr))

	for _, obj := range objArr {
		scope := db.NewScope(obj)
		fields := scope.Fields()

		placeholders := make([]string, 0, len(fields))
		for i := range fields {
			if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) || mainFields[i].Relationship != nil {
				continue
			}
			var vars interface{}
			if (fields[i].Name == "CreatedAt" || fields[i].Name == "UpdatedAt") && fields[i].IsBlank {
				vars = gorm.NowFunc()
			} else {
				vars = fields[i].Field.Interface()
			}
			placeholders = append(placeholders, mainScope.AddToVars(vars))
		}

		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
	}

	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
		mainScope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholdersArr, ", "),
	))

	// Execute and Log
	if err := mainScope.Exec().DB().Error; err != nil {
		return 0, err
	}
	return mainScope.DB().RowsAffected, nil
}

func BatchInsertV2(objArr []interface{},db *gorm.DB) (int64, error) {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return 0, errors.New("slice must not be empty")
	}

	mainObj := objArr[0]
	mainScope := db.NewScope(mainObj)
	mainFields := mainScope.Fields()
	quoted := make([]string, 0, len(mainFields))
	for i := range mainFields {
		// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
		// If field is ignore field, skip it.
		if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || mainFields[i].IsIgnored || mainFields[i].Relationship != nil {
			continue
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}

	placeholdersArr := make([]string, 0, len(objArr))

	for _, obj := range objArr {
		scope := db.NewScope(obj)
		fields := scope.Fields()

		placeholders := make([]string, 0, len(fields))
		for i := range fields {
			if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) || mainFields[i].Relationship != nil {
				continue
			}
			var vars interface{}
			if (fields[i].Name == "CreatedAt" || fields[i].Name == "UpdatedAt") && fields[i].IsBlank {
				vars = gorm.NowFunc()
			} else {
				vars = fields[i].Field.Interface()
			}
			placeholders = append(placeholders, mainScope.AddToVars(vars))
		}

		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
	}

	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
		mainScope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholdersArr, ", "),
	))

	// Execute and Log
	if err := mainScope.Exec().DB().Error; err != nil {
		return 0, err
	}
	return mainScope.DB().RowsAffected, nil
}
