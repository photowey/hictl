package schema

import (
	`bytes`
	`database/sql`
	`errors`
	`fmt`
	`io/ioutil`
	`log`
	`os`
	`path/filepath`
	`strings`
	`text/template`
	`unicode`

	_ "github.com/go-sql-driver/mysql"
	`github.com/spf13/cobra`

	`github.com/photowey/hictl`
	`github.com/photowey/hictl/cmd/internal/common/cmdz`
	`github.com/photowey/hictl/cmd/internal/common/constant`
	`github.com/photowey/hictl/cmd/internal/common/filez`
	`github.com/photowey/hictl/cmd/internal/common/helper`
	`github.com/photowey/hictl/cmd/internal/common/jsonz`
	`github.com/photowey/hictl/cmd/internal/common/stringz`
	db `github.com/photowey/hictl/cmd/internal/database`
	`github.com/photowey/hictl/hictlc/gen`
)

const (
	defaultSchema = "ent/schema"
	genFile       = "package ent\n\n//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema\n"

	schema = "schema"
	mixin  = "mixin"

	hictlConfig = "hictl.json"
)

const (
	int64DataTypeTemplate  = "field.Int64(\"%s\"), // %s"
	stringDataTypeTemplate = "field.String(\"%s\"), // %s"
	timeDataTypeTemplate   = "field.Time(\"%s\").Immutable().Default(time.Now().Local), // %s"
)

const (
	indexTemplate      = "index.Fields(\"%s\"),"
	indexMultiTemplate = "index.Fields(%s),"
)

var (
	SkipFields = []string{"created_at", "updated_at"}
)

func Cmd() *cobra.Command {
	var target string
	var database string
	cmd := &cobra.Command{
		Use:   "init [flags] [schemas]",
		Short: "initialize an environment with zero or more schemas",
		Example: helper.Examples(
			"hictl init Example",
			"hictl init --target entv1/schema --database users User Group",
		),
		Args: func(cmd *cobra.Command, names []string) error {
			for _, name := range names {
				if !unicode.IsUpper(rune(name[0])) {
					return errors.New("schema names must begin with uppercase")
				}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, names []string) {
			if err := writeSchema(target, database, names); err != nil {
				log.Fatalln(fmt.Errorf("hictl/schema: %w", err))
			}
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", defaultSchema, "target directory for schemas")
	cmd.Flags().StringVarP(&database, "database", "d", constant.EmptyString, "reverse Engineering database for schemas")

	return cmd
}

type Template struct {
	Name       string   // name
	Fields     []string // field list
	Indexs     []string //  index list
	Imports    []string // import packages
	EntImports []string // ent import packages
}

// writeSchema initialize an environment for ent codegen.
//
//
func writeSchema(target string, databaseName string, names []string) error {
	var database *db.Database
	if constant.EmptyString != strings.TrimSpace(databaseName) {
		hictlHome := hictl.HomeDir
		hiConf := &db.Hictl{}
		hictlConfigFile := filepath.Join(hictlHome, strings.ToLower(hictlConfig))
		if filez.FileExists(hictlHome, hictlConfig) {
			conf, _ := ioutil.ReadFile(hictlConfigFile)
			err := jsonz.UnmarshalStruct(conf, hiConf)
			if err != nil {
				return err
			}
		}

		conf, ok := hiConf.AcquireDatabase(databaseName)
		if !ok {
			return fmt.Errorf("the database config not present at hiConf home-config file:%s", hictlConfigFile)
		}

		databasez, err := populateDatabaseInfo(conf, databaseName)
		if err != nil {
			return fmt.Errorf("retrieve database table info error:%v", err)
		}
		database = databasez
	}

	if err := createDir(target); err != nil {
		return fmt.Errorf("create dir %s: %w", target, err)
	}
	if len(names) > 0 {
		var table *db.Table
		for _, name := range names {
			if nil != database {
				for _, tableInfo := range database.Tables {
					entity := stringz.Pascal(tableInfo.Name)
					if name == entity {
						table = tableInfo
					}
				}
			}
			if err := gen.ValidSchemaName(name); err != nil {
				return fmt.Errorf("init schema %s: %w", name, err)
			}
			if filez.FileExists(target, strings.ToLower(name+".go")) {
				return fmt.Errorf("init schema %s: already exists", name)
			}
			err := writeSchemaTmpl(target, name, table)
			if err != nil {
				return err
			}
		}
	} else {
		if nil != database {
			for _, table := range database.Tables {
				entity := stringz.Pascal(table.Name)
				if err := gen.ValidSchemaName(entity); err != nil {
					return fmt.Errorf("init schema %s: %w", entity, err)
				}
				if filez.FileExists(target, strings.ToLower(entity+".go")) {
					return fmt.Errorf("init schema %s: already exists", entity)
				}
				err := writeSchemaTmpl(target, entity, table)
				if err != nil {
					return err
				}
			}
		}
	}

	err := writeMixinTmpl(target)
	if err != nil {
		return err
	}

	return nil
}

func populateDatabaseInfo(conf db.Config, databaseName string) (*db.Database, error) {
	dsn := fmt.Sprintf(db.DsnTemplate, conf.UserName, conf.Password, conf.Host, conf.Port, conf.Database)
	driver, err := sql.Open(db.DriverMysql, dsn)
	if err != nil {
		return nil, err
	}
	err = driver.Ping()
	if err != nil {
		return nil, err
	}

	rows, err := driver.Query(db.TableInfoSql, databaseName)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if err != nil {
		return nil, err
	}

	database := &db.Database{
		Name:   databaseName,
		Tables: make([]*db.Table, 0),
	}

	for rows.Next() {
		table := &db.Table{
			Columns: make([]*db.Column, 0),
			Indexs:  make([]*db.Index, 0),
		}
		err := rows.Scan(&table.Name, &table.Comment)
		if err != nil {
			return nil, fmt.Errorf("scan table failed, err:%v", err)
		}
		err = handleTableColumn(driver, database, table)
		if err != nil {
			return nil, err
		}

		err = handleTableIndex(driver, database, table)
		if err != nil {
			return nil, err
		}

		database.Tables = append(database.Tables, table)
	}

	return database, nil
}

func handleTableColumn(driver *sql.DB, database *db.Database, table *db.Table) error {
	// handle column
	rowColumns, err := driver.Query(db.ColumnInfoSql, database.Name, table.Name)
	if err != nil {
		return fmt.Errorf("query table:[%s] column failed, err:%v", table.Name, err)
	}
	for rowColumns.Next() {
		column := &db.Column{}
		err := rowColumns.Scan(
			&column.TableName, &column.ColumnName, &column.ColumnComment,
			&column.NotNull, &column.DataType, &column.DataLength,
			&column.PrimaryKey, &column.MaxBit, &column.MinBit,
		)
		if err != nil {
			return fmt.Errorf("scan table:[%s] column failed, err:%v", table.Name, err)
		}
		table.Columns = append(table.Columns, column)
	}
	_ = rowColumns.Close()

	return nil
}

func handleTableIndex(driver *sql.DB, database *db.Database, table *db.Table) error {
	// handle index
	indexColumns, err := driver.Query(db.IndexInfoSql, database.Name, table.Name)
	if err != nil {
		return fmt.Errorf("query table:[%s] index failed, err:%v", table.Name, err)
	}
	for indexColumns.Next() {
		indexz := &db.Index{}
		err := indexColumns.Scan(
			&indexz.TableName, &indexz.IndexName, &indexz.IndexColumn,
		)
		if err != nil {
			return fmt.Errorf("scan table:[%s] index failed, err:%v", table.Name, err)
		}
		table.Indexs = append(table.Indexs, indexz)
	}

	_ = indexColumns.Close()

	return nil
}

func createDir(target string) error {
	_, err := os.Stat(target)
	if err == nil || !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(target, os.ModePerm); err != nil {
		return fmt.Errorf("creating schema directory: %w", err)
	}
	if target != defaultSchema {
		return nil
	}
	if err := os.WriteFile("ent/generate.go", []byte(genFile), 0644); err != nil {
		return fmt.Errorf("creating generate.go file: %w", err)
	}
	return nil
}

// writeSchemaTmpl writeSchema the schema.tmpl
func writeSchemaTmpl(target string, name string, table *db.Table) error {
	st, _ := populateSchemaTemplate(name, table)

	buffer := bytes.NewBuffer(nil)
	tmplSchema := template.Must(template.New(schema).Parse(db.SchemaTemplate))
	if err := tmplSchema.Execute(buffer, st); err != nil {
		return fmt.Errorf("executing template %s: %w", name, err)
	}
	schemaTarget := filepath.Join(target, strings.ToLower(name+".go"))

	content := string(buffer.Bytes())
	content = strings.ReplaceAll(content, "    \n", "")
	content = strings.ReplaceAll(content, "\t\n", "")
	buf := bytes.NewBufferString(content)
	if err := os.WriteFile(schemaTarget, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", schemaTarget, err)
	}

	cmdz.FormatSourceCode(schemaTarget)
	// 解析-依赖
	cmdz.ResolveDependencies()

	return nil
}

func populateSchemaTemplate(name string, table *db.Table) (*Template, error) {
	st := &Template{
		Name:       name,
		Fields:     make([]string, constant.EmptyCollection),
		Indexs:     make([]string, constant.EmptyCollection),
		Imports:    make([]string, constant.EmptyCollection),
		EntImports: make([]string, constant.EmptyCollection),
	}

	if nil != table && len(table.Columns) > constant.EmptyCollection {
		for _, column := range table.Columns {
			if stringz.ArrayContains(SkipFields, column.ColumnName) {
				continue
			}
			fieldTemplate, pkg := determineFieldTemplate(column.DataType)
			if constant.EmptyString != pkg {
				if stringz.ArrayNotContains(st.Imports, pkg) {
					st.Imports = append(st.Imports, pkg)
				}
			}
			field := fmt.Sprintf(fieldTemplate, column.ColumnName, column.ColumnComment.String)
			st.Fields = append(st.Fields, field)
		}
	}

	if nil != table && len(table.Indexs) > constant.EmptyCollection {
		for _, indexz := range table.Indexs {
			idx, fill := determineIndexSlice(indexz)
			if fill {
				pkg := "\"entgo.io/ent/schema/index\""
				indexTemplate, multi := determineIndexTemplate(indexz)
				if multi {
					wrap := make([]string, 0)
					idxSlice := strings.Split(idx, ",")
					for _, index := range idxSlice {
						wrap = append(wrap, fmt.Sprintf("\"%s\"", index))
					}
					idx = strings.Join(wrap, ",")
				}

				idxFill := fmt.Sprintf(indexTemplate, idx)
				st.Indexs = append(st.Indexs, idxFill)

				if stringz.ArrayNotContains(st.EntImports, pkg) {
					st.EntImports = append(st.EntImports, pkg)
				}
			}
		}
	}

	return st, nil
}

func determineFieldTemplate(databaseDataType string) (string, string) {
	switch databaseDataType {
	case "char", "varchar":
		return stringDataTypeTemplate, constant.EmptyString
	case "int", "tinyint", "middleint", "bigint":
		return int64DataTypeTemplate, constant.EmptyString
	case "timestamp", "datetime":
		return timeDataTypeTemplate, "time"
	}

	return stringDataTypeTemplate, constant.EmptyString
}

func determineIndexSlice(idx *db.Index) (string, bool) {
	switch idx.IndexColumn {
	case "id": // the primary key is id
		return constant.EmptyString, false
	default:
		if idx.IndexName == "PRIMARY" {
			return constant.EmptyString, false // the primary key isn't id, such as: pid
		}

		return idx.IndexColumn, true
	}
}

func determineIndexTemplate(idx *db.Index) (string, bool) {
	if strings.Contains(idx.IndexColumn, ",") {
		return indexMultiTemplate, true
	}

	return indexTemplate, false
}

// writeSchemaTmpl writeSchema the mixin.tmpl
func writeMixinTmpl(target string) error {
	buffer := bytes.NewBuffer(nil)
	tmplMixin := template.Must(template.New(mixin).Parse(db.MixinTemplate))
	if err := tmplMixin.Execute(buffer, constant.EmptyString); err != nil {
		return fmt.Errorf("executing template %s: %w", "mixin", err)
	}
	mixinTarget := filepath.Join(target, strings.ToLower(mixin+".go"))
	if err := os.WriteFile(mixinTarget, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", mixinTarget, err)
	}
	return nil
}
