package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
	"github.com/iancoleman/strcase"
)
var sqlMappings = make(map[string]*SQLMapping)
//func LoadTableInfo(db *sql.DB, tableName string)DbTableMeta{
//	dbMeta ,err := LoadMysqlMeta(db,tableName)
//	return dbMeta,err
//}
// DbTableMeta table meta data
type DbTableMeta interface {
	Columns() []ColumnMeta
	SQLType() string
	SQLDatabase() string
	TableName() string
	DDL() string
}

type ColumnMeta interface {
	Name() string
	String() string
	Nullable() bool
	DatabaseTypeName() string
	DatabaseTypePretty() string
	Index() int
	IsPrimaryKey() bool
	IsAutoIncrement() bool
	IsArray() bool
	ColumnType() string
	Notes() string
	Comment() string
	ColumnLength() int64
	DefaultValue() string
}

type columnMeta struct {
	index int
	// ct              *sql.ColumnType
	nullable         bool
	isPrimaryKey     bool
	isAutoIncrement  bool
	isArray          bool
	colDDL           string
	columnType       string
	columnLen        int64
	defaultVal       string
	notes            string
	comment          string
	databaseTypeName string
	name             string
}


type dbTableMeta struct {
	sqlType       string
	sqlDatabase   string
	tableName     string
	columns       []*columnMeta
	ddl           string
	primaryKeyPos int
}

// ColumnType column type
func (ci *columnMeta) ColumnType() string {
	return ci.columnType
}

// Notes notes on column generation
func (ci *columnMeta) Notes() string {
	return ci.notes
}

// Comment column comment
func (ci *columnMeta) Comment() string {
	return ci.comment
}

// ColumnLength column length for text or varhar
func (ci *columnMeta) ColumnLength() int64 {
	return ci.columnLen
}

// DefaultValue default value of column
func (ci *columnMeta) DefaultValue() string {
	return ci.defaultVal
}

// Name name of column
func (ci *columnMeta) Name() string {
	return ci.name
}

// Index index of column in db
func (ci *columnMeta) Index() int {
	return ci.index
}

// IsAutoIncrement return is column is a primary key column
func (ci *columnMeta) IsPrimaryKey() bool {
	return ci.isPrimaryKey
}

// IsArray return is column is an array type
func (ci *columnMeta) IsArray() bool {
	return ci.isArray
}

// IsAutoIncrement return is column is an auto increment column
func (ci *columnMeta) IsAutoIncrement() bool {
	return ci.isAutoIncrement
}

// String friendly string for columnMeta
func (ci *columnMeta) String() string {
	return fmt.Sprintf("[%2d] %-45s  %-20s null: %-6t primary: %-6t isArray: %-6t auto: %-6t col: %-15s len: %-7d default: [%s]",
		ci.index, ci.name, ci.DatabaseTypePretty(),
		ci.nullable, ci.isPrimaryKey, ci.isArray,
		ci.isAutoIncrement, ci.columnType, ci.columnLen, ci.defaultVal)
}

// Nullable reports whether the column may be null.
// If a driver does not support this property ok will be false.
func (ci *columnMeta) Nullable() bool {
	return ci.nullable
}

// ColDDL string of the ddl for the column
func (ci *columnMeta) ColDDL() string {
	return ci.colDDL
}

// DatabaseTypeName returns the database system name of the column type. If an empty
// string is returned the driver type name is not supported.
// Consult your driver documentation for a list of driver data types. Length specifiers
// are not included.
// Common type include "VARCHAR", "TEXT", "NVARCHAR", "DECIMAL", "BOOL", "INT", "BIGINT".
func (ci *columnMeta) DatabaseTypeName() string {
	return ci.databaseTypeName
}

// DatabaseTypePretty string of the db type
func (ci *columnMeta) DatabaseTypePretty() string {
	if ci.columnLen > 0 {
		return fmt.Sprintf("%s(%d)", ci.columnType, ci.columnLen)
	}

	return ci.columnType
}

// PrimaryKeyPos ordinal pos of primary key
func (m *dbTableMeta) PrimaryKeyPos() int {
	return m.primaryKeyPos
}

// SQLType sql db type
func (m *dbTableMeta) SQLType() string {
	return m.sqlType
}

// SQLDatabase sql database name
func (m *dbTableMeta) SQLDatabase() string {
	return m.sqlDatabase
}

// TableName sql table name
func (m *dbTableMeta) TableName() string {
	return m.tableName
}

// Columns ColumnMeta for columns in a sql table
func (m *dbTableMeta) Columns() []ColumnMeta {

	cols := make([]ColumnMeta, len(m.columns))
	for i, v := range m.columns {
		cols[i] = ColumnMeta(v)
	}
	return cols
}

// DDL string for a sql table
func (m *dbTableMeta) DDL() string {
	return m.ddl
}


// FieldInfo codegen info for each column in sql table
type FieldInfo struct {
	Index                 int
	GoFieldName           string
	GoFieldType           string
	GoAnnotations         []string
	JSONFieldName         string
	ProtobufFieldName     string
	ProtobufType          string
	ProtobufPos           int
	Comment               string
	Notes                 string
	Code                  string
	FakeData              interface{}
	ColumnMeta            ColumnMeta
	PrimaryKeyFieldParser string
	PrimaryKeyArgName     string
	GormAnnotation        string
	JSONAnnotation        string
	XMLAnnotation         string
	DBAnnotation          string
	GoGoMoreTags          string
}

// GenerateFieldsTypes FieldInfo slice from DbTableMeta
func (c *Config) GenerateFieldsTypes(dbMeta DbTableMeta) ([]*FieldInfo, error) {

	var fields []*FieldInfo
	field := ""
	for i, col := range dbMeta.Columns() {
		fieldName := col.Name()

		fi := &FieldInfo{
			Index: i,
		}

		//valueType := "int64"
		valueType, err := SQLTypeToGoType(strings.ToLower(col.DatabaseTypeName()), col.Nullable(), false)
		if err != nil { // unknown type
			fmt.Printf("table: %s unable to generate struct field: %s type: %s error: %v\n", dbMeta.TableName(), fieldName, col.DatabaseTypeName(), err)
			continue
		}

		//fieldName = camelToUpperCamel(fieldName)
		//fmt.Println(fieldName)
		fieldName = Replace("{{FmtFieldName (stringifyFirstChar .) }}", fieldName)
		//fieldName = checkDupeFieldName(fields, fieldName)

		fi.GormAnnotation = createGormAnnotation(col)
		fi.JSONAnnotation = createJSONAnnotation("snake", col)
		//fi.XMLAnnotation = createXMLAnnotation(c.XMLNameFormat, col)
		//fi.DBAnnotation = createDBAnnotation(col)

		var annotations []string
		//if c.AddGormAnnotation {
			annotations = append(annotations, fi.GormAnnotation)
		//}

		//if c.AddJSONAnnotation {
			annotations = append(annotations, fi.JSONAnnotation)
		//}
		//
		//if c.AddXMLAnnotation {
		//	annotations = append(annotations, fi.XMLAnnotation)
		//}
		//
		//if c.AddDBAnnotation {
		//	annotations = append(annotations, fi.DBAnnotation)
		//}

		gogoTags := []string{fi.GormAnnotation, fi.JSONAnnotation, fi.XMLAnnotation, fi.DBAnnotation}
		GoGoMoreTags := strings.Join(gogoTags, " ")

		//if c.AddProtobufAnnotation {
		//	annotation, err := createProtobufAnnotation(c.ProtobufNameFormat, col)
		//	if err == nil {
		//		annotations = append(annotations, annotation)
		//	}
		//}

		if len(annotations) > 0 {
			field = fmt.Sprintf("%s %s `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))
		} else {
			field = fmt.Sprintf("%s %s", fieldName, valueType)
		}

		field = fmt.Sprintf("//%s\n    %s", col.String(), field)
		if col.Comment() != "" {
			field = fmt.Sprintf("%s // %s", field, col.Comment())
		}

	//	sqlMapping, _ := SQLTypeToMapping(strings.ToLower(col.DatabaseTypeName()))
		goType, _ := SQLTypeToGoType(strings.ToLower(col.DatabaseTypeName()), false, false)
		//protobufType, _ := SQLTypeToProtobufType(col.DatabaseTypeName())

		// fmt.Printf("protobufType: %v  DatabaseTypeName: %v\n", protobufType, col.DatabaseTypeName())

		//fakeData := createFakeData(goType, fieldName)

		//if c.Verbose {
		//	fmt.Printf("table: %-10s type: %-10s fieldname: %-20s val: %v\n", c.DatabaseTypeName(), goType, fieldName, fakeData)
		//	spew.Dump(fakeData)
		//}

		//fmt.Printf("%+v", fakeData)
		primaryKeyFieldParser := ""
		if col.IsPrimaryKey() {
			var ok bool
			primaryKeyFieldParser, ok = parsePrimaryKeys[goType]
			if !ok {
				primaryKeyFieldParser = "unsupported"
			}
		}

		fi.Code = field
		fi.GoFieldName = fieldName
		fi.GoFieldType = valueType
		fi.GoAnnotations = annotations
		//fi.FakeData = fakeData
		fi.Comment = col.String()
		//fi.JSONFieldName = formatFieldName(c.JSONNameFormat, col.Name())
	///	fi.ProtobufFieldName = formatFieldName(c.ProtobufNameFormat, col.Name())
	//	fi.ProtobufType = protobufType
		fi.ProtobufPos = i + 1
		fi.ColumnMeta = col
		fi.PrimaryKeyFieldParser = primaryKeyFieldParser
		//fi.SQLMapping = sqlMapping
		fi.GoGoMoreTags = GoGoMoreTags

		//fi.JSONFieldName = checkDupeJSONFieldName(fields, fi.JSONFieldName)
		//fi.ProtobufFieldName = checkDupeProtoBufFieldName(fields, fi.ProtobufFieldName)

		fields = append(fields, fi)
	}
	return fields, nil
}

func createGormAnnotation(c ColumnMeta) string {
	buf := bytes.Buffer{}

	key := c.Name()
	buf.WriteString("gorm:\"")

	if c.IsPrimaryKey() {
		buf.WriteString("primary_key;")
	}
	if c.IsAutoIncrement() {
		buf.WriteString("AUTO_INCREMENT;")
	}

	buf.WriteString("column:")
	buf.WriteString(key)
	buf.WriteString(";")

	if c.DatabaseTypeName() != "" {
		buf.WriteString("type:")
		buf.WriteString(c.DatabaseTypeName())
		buf.WriteString(";")

		if c.ColumnLength() > 0 {
			buf.WriteString(fmt.Sprintf("size:%d;", c.ColumnLength()))
		}

		if c.DefaultValue() != "" {
			value := c.DefaultValue()
			value = strings.Replace(value, "\"", "'", -1)

			if value == "NULL" || value == "null" {
				value = ""
			}

			if value != "" && !strings.Contains(value, "()") {
				buf.WriteString(fmt.Sprintf("default:%s;", value))
			}
		}

	}

	buf.WriteString("\"")
	return buf.String()
}

func createJSONAnnotation(nameFormat string, c ColumnMeta) string {
	name := formatFieldName(nameFormat, c.Name())
	return fmt.Sprintf("json:\"%s\"", name)
}

func formatFieldName(nameFormat string, name string) string {

	var jsonName string
	switch nameFormat {
	case "snake":
		jsonName = strcase.ToSnake(name)
	case "camel":
		jsonName = strcase.ToCamel(name)
	case "lower_camel":
		jsonName = strcase.ToLowerCamel(name)
	case "none":
		jsonName = name
	default:
		jsonName = name
	}
	return jsonName
}

// SQLTypeToGoType map a sql type to a go type
func SQLTypeToGoType(sqlType string, nullable bool, gureguTypes bool) (string, error) {


	mapping, err := SQLTypeToMapping(sqlType)
	if err != nil {
		return "", err
	}

	if nullable && gureguTypes {
		return mapping.GureguType, nil
	} else if nullable {
		return mapping.GoNullableType, nil
	} else {
		return mapping.GoType, nil
	}
}


var parsePrimaryKeys = map[string]string{
	"uint8":     "parseUint8",
	"uint16":    "parseUint16",
	"uint32":    "parseUint32",
	"uint64":    "parseUint64",
	"int":       "parseInt",
	"int8":      "parseInt8",
	"int16":     "parseInt16",
	"int32":     "parseInt32",
	"int64":     "parseInt64",
	"string":    "parseString",
	"uuid.UUID": "parseUUID",
}

// SQLMappings mappings for sql types to json, go etc
type SQLMappings struct {
	SQLMappings []*SQLMapping `json:"mappings"`
}

// SQLMapping mapping
type SQLMapping struct {
	// SQLType sql type reported from db
	SQLType string `json:"sql_type"`

	// GoType mapped go type
	GoType string `json:"go_type"`

	// JSONType mapped json type
	JSONType string `json:"json_type"`

	// ProtobufType mapped protobuf type
	ProtobufType string `json:"protobuf_type"`

	// GureguType mapped go type using Guregu
	GureguType string `json:"guregu_type"`

	// GoNullableType mapped go type using nullable
	GoNullableType string `json:"go_nullable_type"`

	// SwaggerType mapped type
	SwaggerType string `json:"swagger_type"`
}

// ProcessMappings process the json for mappings to load sql mappings
func ProcessMappings(fpath string) error {
	absPath, err := filepath.Abs(fpath)
	if err != nil {
		absPath = fpath
	}
	// fmt.Printf("Loaded template from file: %s\n", fpath)
	//tpl = &GenTemplate{Name: "file://" + absPath, Content: string(b)}

	mappingJsonstring, err := ioutil.ReadFile(absPath)
	fmt.Println(absPath)
	if err != nil {
		return err
	}
	var mappings = &SQLMappings{}
	err = json.Unmarshal(mappingJsonstring, mappings)
	if err != nil {
		fmt.Printf("Error unmarshalling json error: %v\n", err)
		return err
	}

	for _, value := range mappings.SQLMappings {

		sqlMappings[value.SQLType] = value
	}
	fmt.Println(mappings)

	return nil
}

// SQLTypeToMapping retrieve a SQLMapping based on a sql type
func SQLTypeToMapping(sqlType string) (*SQLMapping, error) {
	sqlType = cleanupSQLType(sqlType)

	mapping, ok := sqlMappings[sqlType]
	if !ok {
		return nil, fmt.Errorf("unknown sql type: %s", sqlType)
	}
	return mapping, nil
}

func cleanupSQLType(sqlType string) string {
	sqlType = strings.ToLower(sqlType)
	sqlType = strings.Trim(sqlType, " \t")
	sqlType = strings.ToLower(sqlType)
	idx := strings.Index(sqlType, "(")
	if idx > -1 {
		sqlType = sqlType[0:idx]
	}
	return sqlType
}

func camelToUpperCamel(s string) string {
	ss := strings.Split(s, "")
	ss[0] = strings.ToUpper(ss[0])
	return strings.Join(ss, "")
}

// Replace takes a template based name format and will render a name using it
func Replace(nameFormat, name string) string {
	var tpl bytes.Buffer
	//fmt.Printf("Replace: %s\n",nameFormat)
	t := template.Must(template.New("t1").Funcs(replaceFuncMap).Parse(nameFormat))

	if err := t.Execute(&tpl, name); err != nil {
		//fmt.Printf("Error creating name format: %s error: %v\n", nameFormat, err)
		return name
	}
	result := tpl.String()

	result = strings.Trim(result, " \t")
	result = strings.Replace(result, " ", "_", -1)
	result = strings.Replace(result, "\t", "_", -1)

	//fmt.Printf("Replace( '%s' '%s')= %s\n",nameFormat, name, result)
	return result
}

