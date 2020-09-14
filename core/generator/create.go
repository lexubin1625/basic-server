package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"text/template"
	"strings"
	"unicode"

	//"github.com/jinzhu/inflection"
)

var (
	templateDir = "./core/generator/template/"
)

type GenTemplate struct {
	Name    string
	Content string
}

type Config struct {
	JSONNameFormat        string
}

var replaceFuncMap = template.FuncMap{
	//"singular":           inflection.Singular,
	//"pluralize":          inflection.Plural,
	"title":              strings.Title,
	"toLower":            strings.ToLower,
	"toUpper":            strings.ToUpper,
	//"toLowerCamelCase":   camelToLowerCamel,
	"toUpperCamelCase":   camelToUpperCamel,
//	"toSnakeCase":        snaker.CamelToSnake,
	"StringsJoin":        strings.Join,
	//"replace":            replace,
	"stringifyFirstChar": stringifyFirstChar,
	"FmtFieldName":       FmtFieldName,
}

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

// WriteTemplate write a template out
func (c *Config) WriteTemplate(genTemplate *GenTemplate, data map[string]interface{}, outputFile string) error {
	//fmt.Printf("WriteTemplate %s\n", outputFile)

	//if  Exists(outputFile) {
	//	fmt.Printf("not overwriting %s\n", outputFile)
	//	return nil
	//}


	//dir := filepath.Dir(outputFile)

	rt, err := c.GetTemplate(genTemplate)
	if err != nil {
		return fmt.Errorf("error in loading %s template, error: %v", genTemplate.Name, err)
	}
	var buf bytes.Buffer
	err = rt.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("error in rendering %s: %s", genTemplate.Name, err.Error())
	}

	fileContents, err := c.format(genTemplate, buf.Bytes(), outputFile)
	if err != nil {
		return fmt.Errorf("error writing %s - error: %v", outputFile, err)
	}

	err = ioutil.WriteFile(outputFile, fileContents, 0777)
	if err != nil {
		return fmt.Errorf("error writing %s - error: %v", outputFile, err)
	}

	return nil
}

// GetTemplate return a Template based on a name and template contents
func (c *Config) GetTemplate(genTemplate *GenTemplate) (*template.Template, error) {
	//var s State
	var funcMap = template.FuncMap{
		//"ReplaceFileNamingTemplate":  c.ReplaceFileNamingTemplate,
		//"ReplaceModelNamingTemplate": c.ReplaceModelNamingTemplate,
		//"ReplaceFieldNamingTemplate": c.ReplaceFieldNamingTemplate,
		//"stringifyFirstChar":         stringifyFirstChar,
		//"singular":                   inflection.Singular,
		//"pluralize":                  inflection.Plural,
		"title":                      strings.Title,
		"toLower":                    strings.ToLower,
		"toUpper":                    strings.ToUpper,
		//"toLowerCamelCase":           camelToLowerCamel,
		//"toUpperCamelCase":           camelToUpperCamel,
		//"FormatSource":               FormatSource,
		//"toSnakeCase":                snaker.CamelToSnake,
		//"markdownCodeBlock":          markdownCodeBlock,
		//"wrapBash":                   wrapBash,
		//"escape":                     escape,
		//"GenerateTableFile":          c.GenerateTableFile,
		//"GenerateFile":               c.GenerateFile,
		//"ToJSON":                     ToJSON,
		//"spew":                       Spew,
		//"set":                        s.Set,
		//"inc":                        s.Inc,
		"StringsJoin":                strings.Join,
		//"replace":                    replace,
		//"hasField":                   hasField,
		//"FmtFieldName":               FmtFieldName,
		//"copy":                       c.FileSystemCopy,
		//"mkdir":                      c.Mkdir,
		//"touch":                      c.Touch,
		//"pwd":                        Pwd,
		//"config":                     c.DisplayConfig,
	}

	baseName := filepath.Base(genTemplate.Name)

	tmpl, err := template.New(baseName).Option("missingkey=error").Funcs(funcMap).Parse(genTemplate.Content)
	if err != nil {
		return nil, err
	}

	//if baseName == "api.go.tmpl" ||
	//	baseName == "dao_gorm.go.tmpl" ||
	//	baseName == "dao_sqlx.go.tmpl" ||
	//	baseName == "code_dao_sqlx.md.tmpl" ||
	//	baseName == "code_dao_gorm.md.tmpl" ||
	//	baseName == "code_http.md.tmpl" {
	//
	//	operations := []string{"add", "delete", "get", "getall", "update"}
	//	for _, op := range operations {
	//		var filename string
	//		if baseName == "api.go.tmpl" {
	//			filename = fmt.Sprintf("api_%s.go.tmpl", op)
	//		}
	//		if baseName == "dao_gorm.go.tmpl" {
	//			filename = fmt.Sprintf("dao_gorm_%s.go.tmpl", op)
	//		}
	//		if baseName == "dao_sqlx.go.tmpl" {
	//			filename = fmt.Sprintf("dao_sqlx_%s.go.tmpl", op)
	//		}
	//
	//		if baseName == "code_dao_sqlx.md.tmpl" {
	//			filename = fmt.Sprintf("dao_sqlx_%s.go.tmpl", op)
	//		}
	//		if baseName == "code_dao_gorm.md.tmpl" {
	//			filename = fmt.Sprintf("dao_gorm_%s.go.tmpl", op)
	//		}
	//		if baseName == "code_http.md.tmpl" {
	//			filename = fmt.Sprintf("api_%s.go.tmpl", op)
	//		}
	//
	//		var subTemplate *GenTemplate
	//		if subTemplate, err = c.TemplateLoader(filename); err != nil {
	//			fmt.Printf("Error loading template %v\n", err)
	//			return nil, err
	//		}
	//
	//		// fmt.Printf("loading sub template %v\n", filename)
	//		tmpl.Parse(subTemplate.Content)
	//	}
	//}

	return tmpl, nil
}

func (c *Config) format(genTemplate *GenTemplate, content []byte, outputFile string) ([]byte, error) {
	extension := filepath.Ext(outputFile)
	if extension == ".go" {
		formattedSource, err := format.Source([]byte(content))
		if err != nil {
			return nil, fmt.Errorf("error in formatting template: %s outputfile: %s source: %s", genTemplate.Name, outputFile, err.Error())
		}

		fileContents := NormalizeNewlines(formattedSource)
		//if c.LineEndingCRLF {
		//	fileContents = CRLFNewlines(formattedSource)
		//}
		return fileContents, nil
	}

	fileContents := NormalizeNewlines([]byte(content))
	//if c.LineEndingCRLF {
	//	fileContents = CRLFNewlines(fileContents)
	//}
	return fileContents, nil
}

// NormalizeNewlines normalizes \r\n (windows) and \r (mac)
// into \n (unix)
func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}

// LoadTemplate return template from template dir, falling back to the embedded templates
func LoadTemplate(filename string) (tpl *GenTemplate, err error) {
	//baseName := filepath.Base(filename)
	// fmt.Printf("LoadTemplate: %s / %s\n", filename, baseName)

	if templateDir != "" {
		fpath := filepath.Join(templateDir, filename)
		var b []byte
		b, err = ioutil.ReadFile(fpath)
		if err == nil {

			absPath, err := filepath.Abs(fpath)
			if err != nil {
				absPath = fpath
			}
			// fmt.Printf("Loaded template from file: %s\n", fpath)
			tpl = &GenTemplate{Name: "file://" + absPath, Content: string(b)}
			return tpl, nil
		}
	}

	return nil, nil
	//content, err := baseTemplates.FindString(baseName)
	//if err != nil {
	//	return nil, fmt.Errorf("%s not found internally", baseName)
	//}
	//if *verbose {
	//	fmt.Printf("Loaded template from app: %s\n", filename)
	//}

	//tpl = &GenTemplate{Name: "internal://" + filename, Content: content}
	//return tpl, nil
}


// convert first character ints to strings
func stringifyFirstChar(str string) string {
	first := str[:1]

	i, err := strconv.ParseInt(first, 10, 8)

	if err != nil {
		return str
	}

	return intToWordMap[i] + "_" + str[1:]
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}


// FmtFieldName formats a string as a struct key
//
// Example:
// 	fmtFieldName("foo_id")
// Output: FooID
func FmtFieldName(s string) string {
	name := lintFieldName(s)
	runes := []rune(name)
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	fieldName := string(runes)
	fieldName = RenameReservedName(fieldName)
	// fmt.Printf("FmtFieldName:%s=%s\n", s, fieldName)
	return fieldName
}
func isAllLower(name string) (allLower bool) {
	allLower = true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	return
}

func lintAllLowerFieldName(name string) string {
	runes := []rune(name)
	if u := strings.ToUpper(name); commonInitialisms[u] {
		copy(runes[0:], []rune(u))
	} else {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}

func lintFieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}

	for len(name) > 0 && name[0] == '_' {
		name = name[1:]
	}

	allLower := isAllLower(name)

	if allLower {
		return lintAllLowerFieldName(name)
	}

	return lintMixedFieldName(name)
}

func lintMixedFieldName(name string) string {
	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan

	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))

		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

var reservedFieldNames = map[string]bool{
	"TableName":  true,
	"BeforeSave": true,
	"Prepare":    true,
	"Validate":   true,
	"type":       true,
}

// RenameReservedName renames a reserved word
func RenameReservedName(s string) string {
	_, match := reservedFieldNames[s]
	if match {
		return fmt.Sprintf("%s_", s)
	}

	return s
}