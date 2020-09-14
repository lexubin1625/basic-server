package generator

import (
	//"bytes"
	"database/sql"
	"fmt"
	"github.com/jimsmart/schema"
	//"io/ioutil"
	//"path/filepath"
	"regexp"
	"strconv"
	"strings"
)
// ModelInfo info for a sql table
type ModelInfo struct {
	Index           int
	IndexPlus1      int
	PackageName     string
	StructName      string
	ShortStructName string
	TableName       string
	Fields          []string
}



//type dbTableMeta struct {
//	sqlType       string
//	sqlDatabase   string
//	tableName     string
//	columns       []*columnMeta
//	ddl           string
//	primaryKeyPos int
//}

// DbTableMeta table meta data
//type DbTableMeta interface {
//	Columns() []ColumnMeta
//	SQLType() string
//	SQLDatabase() string
//	TableName() string
//	DDL() string
//}

// ModelInfo info for a sql table
//type ModelInfo struct {
//	Index           int
//	IndexPlus1      int
//	PackageName     string
//	StructName      string
//	ShortStructName string
//	TableName       string
//	Fields          []string
//	DBMeta          DbTableMeta
//	Instance        interface{}
//	CodeFields      []*FieldInfo
//}

func LoadMysqlMeta(db *sql.DB, tableName string)(DbTableMeta,error){


	m := &dbTableMeta{
		sqlType:     "mysql",
		sqlDatabase: "sql_type",
		tableName:   tableName,
	}

	ddl,_ := mysqlLoadDDL(db,tableName)

	m.ddl = ddl
	colsDDL, primaryKeys := mysqlParseDDL(ddl)


	cols, err := schema.Table(db, tableName)
	if err != nil {
		return nil,err
	}
	m.columns = make([]*columnMeta, len(cols))

	for i, v := range cols {
		notes := ""
		nullable, ok := v.Nullable()
		if !ok {
			nullable = false
		}

		colDDL := colsDDL[v.Name()]

		isAutoIncrement := strings.Index(colDDL, "AUTO_INCREMENT") > -1
		isUnsigned := strings.Index(colDDL, " unsigned ") > -1 || strings.Index(colDDL, " UNSIGNED ") > -1

		_, isPrimaryKey := find(primaryKeys, v.Name())
		defaultVal := ""
		columnType, columnLen := ParseSQLType(v.DatabaseTypeName())

		if isUnsigned {
			notes = notes + " column is set for unsigned"
			columnType = "u" + columnType
		}

		comment := ""
		commentIdx := strings.Index(colDDL, "COMMENT '")
		if commentIdx > -1 {
			re := regexp.MustCompile("COMMENT '(.*?)'")
			match := re.FindStringSubmatch(colDDL)
			if len(match) > 0 {
				comment = match[1]
			}
		}

		colMeta := &columnMeta{
			index:            i,
			name:             v.Name(),
			databaseTypeName: columnType,
			nullable:         nullable,
			isPrimaryKey:     isPrimaryKey,
			isAutoIncrement:  isAutoIncrement,
			colDDL:           colDDL,
			defaultVal:       defaultVal,
			columnType:       columnType,
			columnLen:        columnLen,
			notes:            strings.Trim(notes, " "),
			comment:          comment,
		}
		m.columns[i] = colMeta
	}

	return m,nil
	//fmt.Println(m.columns[0], err)
}

// 获取mysql表结构
func mysqlLoadDDL(db *sql.DB, tableName string) (ddl string, err error) {
	ddlSQL := fmt.Sprintf("SHOW CREATE TABLE `%s`;", tableName)
	res, err := db.Query(ddlSQL)
	if err != nil {
		return "", fmt.Errorf("unable to load ddl from mysql: %v", err)
	}

	defer res.Close()
	var ddl1 string
	var ddl2 string
	if res.Next() {
		err = res.Scan(&ddl1, &ddl2)
		if err != nil {
			return "", fmt.Errorf("unable to load ddl from mysql Scan: %v", err)
		}
	}
	return ddl2, nil

}

func mysqlParseDDL(ddl string) (colsDDL map[string]string, primaryKeys []string) {
	colsDDL = make(map[string]string)
	lines := strings.Split(ddl, "\n")
	for _, line := range lines {
		line = strings.Trim(line, " \t")
		if strings.HasPrefix(line, "CREATE TABLE") || strings.HasPrefix(line, "(") || strings.HasPrefix(line, ")") {
			continue
		}

		if line[0] == '`' {
			idx := indexAt(line, "`", 1)
			if idx > 0 {
				name := line[1:idx]
				colDDL := line[idx+1 : len(line)-1]
				colsDDL[name] = colDDL
			}
		} else if strings.HasPrefix(line, "PRIMARY KEY") {
			var primaryKeyNums = strings.Count(line, "`") / 2
			var count = 0
			var currentIdx = 0
			var idxL = 0
			var idxR = 0
			for {
				if count >= primaryKeyNums {
					break
				}
				count++
				idxL = indexAt(line, "`", currentIdx)
				currentIdx = idxL + 1
				idxR = indexAt(line, "`", currentIdx)
				currentIdx = idxR + 1
				primaryKeys = append(primaryKeys, line[idxL+1:idxR])
			}
		}
	}
	return
}

func indexAt(s, sep string, n int) int {
	idx := strings.Index(s[n:], sep)
	if idx > -1 {
		idx += n
	}
	return idx
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// ParseSQLType parse sql type and return raw type and length
func ParseSQLType(dbType string) (resultType string, dbTypeLen int64) {

	resultType = strings.ToLower(dbType)
	dbTypeLen = -1
	idx1 := strings.Index(resultType, "(")
	idx2 := strings.Index(resultType, ")")

	if idx1 > -1 && idx2 > -1 {
		sizeStr := resultType[idx1+1 : idx2]
		resultType = resultType[0:idx1]
		i, err := strconv.Atoi(sizeStr)
		if err == nil {
			dbTypeLen = int64(i)
		}
	}

	// fmt.Printf("dbType: %-20s %-20s %d\n", dbType, resultType, dbTypeLen)
	return resultType, dbTypeLen
}


// GenerateModelInfo generates a struct for the given table.
func GenerateModelInfo(dbMeta DbTableMeta,
	tableName string,conf *Config) (*ModelInfo, error) {

	//structName := Replace(conf.ModelNamingTemplate, tableName)
	//structName = CheckForDupeTable(tables, structName)

	fields, err := conf.GenerateFieldsTypes(dbMeta)
	if err != nil {
		return nil, err
	}

	//generator := dynamicstruct.NewStruct()
	//
	//noOfPrimaryKeys := 0
	//for i, c := range fields {
	//	meta := dbMeta.Columns()[i]
	//	jsonName := formatFieldName(conf.JSONNameFormat, meta.Name())
	//	tag := fmt.Sprintf(`json:"%s"`, jsonName)
	//	fakeData := c.FakeData
	//	generator = generator.AddField(c.GoFieldName, fakeData, tag)
	//	if meta.IsPrimaryKey() {
	//		//c.PrimaryKeyArgName = RenameReservedName(strcase.ToLowerCamel(c.GoFieldName))
	//		c.PrimaryKeyArgName = fmt.Sprintf("arg%s", strcase.ToCamel(c.GoFieldName))
	//		noOfPrimaryKeys++
	//	}
	//}
	//
	//instance := generator.Build().New()
	//
	//err = faker.FakeData(&instance)
	//if err != nil {
	//	fmt.Println(err)
	//}
	// fmt.Printf("%+v", instance)

	var code []string
	for _, f := range fields {

		//if f.PrimaryKeyFieldParser == "unsupported" {
		//	return nil, fmt.Errorf("unable to generate code for table: %s, primary key column: [%d] %s has unsupported type: %s / %s",
		//		dbMeta.TableName(), f.ColumnMeta.Index(), f.ColumnMeta.Name(), f.ColumnMeta.DatabaseTypeName(), f.GoFieldType)
		//}
		code = append(code, f.Code)
	}

	var modelInfo = &ModelInfo{
		//PackageName:     conf.ModelPackageName,
		//StructName:      structName,
		//TableName:       tableName,
		//ShortStructName: strings.ToLower(string(structName[0])),
		Fields:          code,
		//CodeFields:      fields,
		//DBMeta:          dbMeta,
		//Instance:        instance,
	}

	return modelInfo, nil
}


