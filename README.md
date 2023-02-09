# go-sqlink
Decodes sql.Rows into slice of structs

DecodeRows decodes SQL result rows into slice of structs.
It scans struct field tags to match column name in result row.
Struct field data type should match column value data type too.

Usage:
