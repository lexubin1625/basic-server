package main

import (
	"basic-server/core/config"
	"basic-server/core/db"
	"basic-server/core/generator"
	"fmt"
)

func main() {
	// 初始化配置
	var conf config.Config
	conf.InitConf()

	// 数据库初始化
	db.New()

	dbs,_ := db.Get("default",true)
	// 获取表结构信息
	tableName := "logs"
	dbMeta ,_ := generator.LoadMysqlMeta(dbs.DB(),tableName)
	var cc generator.Config
	//fmt.Println(dbMeta)
	ModelTmpl, err := generator.LoadTemplate("model.go.tmpl")
	if err != nil {

	}
	generator.ProcessMappings("./core/generator/template/maping.json")
	model,_ := generator.GenerateModelInfo(dbMeta,tableName,nil)
	fmt.Println(model)
	var modelInfo = map[string]interface{}{
		"TableInfo":model,
		"StructName":"Logs",
	}
	fmt.Println(cc.WriteTemplate(ModelTmpl,modelInfo,"./model/meta/logs.go"))
	fmt.Println(dbMeta)
}
