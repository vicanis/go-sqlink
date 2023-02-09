package sqlink

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

const tagName = "sql"

// DecodeRows decodes SQL result rows into slice of structs.
// It scans struct field tags to match column name in result row.
// Struct field data type should match column value data type too.
func DecodeRows(rows *sql.Rows, data any) (err error) {
	if data == nil {
		err = errors.New("data is nil")
		return
	}

	val := reflect.ValueOf(data).Elem()
	t := val.Type()

	if t.Kind() != reflect.Slice {
		err = fmt.Errorf("data type is not slice: %s", t.Kind())
		return
	}

	itemType := t.Elem()

	if itemType.Kind() != reflect.Struct {
		err = fmt.Errorf("item type is not struct: %s", itemType.Kind())
		return
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(
				"scan failed: %s",
				e,
			)
		}
	}()

	for rows.Next() {
		dataStruct := reflect.New(itemType)

		row := make([]any, 0)
		names := make([]string, 0)

		for _, col := range columns {
			names = append(names, col.Name())

			switch col.ScanType().Kind() {
			case reflect.Int:
				var i int
				row = append(row, &i)
			case reflect.Int32:
				var i int32
				row = append(row, &i)
			case reflect.Int64:
				var i int64
				row = append(row, &i)

			case reflect.Slice:
				elemType := col.ScanType().Elem()

				if elemType.Kind() != reflect.Uint8 {
					err = fmt.Errorf(
						"column data type is slice with element data type %s, which is not supported yet",
						elemType.Kind(),
					)
					return
				}

				fallthrough

			case reflect.String:
				var i string
				row = append(row, &i)

			default:
				err = fmt.Errorf(
					"column data type %s is not supported yet",
					col.ScanType().Kind(),
				)
				return
			}
		}

		err = rows.Scan(row...)
		if err != nil {
			err = fmt.Errorf("scan failed: %w", err)
			return
		}

		for i, item := range row {
			elem := reflect.ValueOf(item).Elem()

		fieldloop:
			for j := 0; j < itemType.NumField(); j++ {
				field := itemType.Field(j)

				tag := field.Tag.Get(tagName)
				if tag != names[i] {
					continue
				}

				switch field.Type.Kind() {
				case reflect.Int:
					switch elem.Kind() {
					case reflect.Int32:
						fallthrough
					case reflect.Int64:
						i := int(elem.Int())
						dataStruct.Elem().Field(j).Set(reflect.ValueOf(i))
						continue fieldloop
					}

				case reflect.String:
					switch elem.Kind() {
					case reflect.String:
						dataStruct.Elem().Field(j).Set(elem)
						continue fieldloop
					}
				}

				err = fmt.Errorf(
					"data type mismatch: struct field %s and column %s",
					field.Type.Kind(),
					elem.Kind(),
				)
				return
			}
		}

		val.Set(
			reflect.Append(
				val,
				dataStruct.Elem(),
			),
		)
	}

	return
}
