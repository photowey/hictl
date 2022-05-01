package executor

import (
	`bytes`
	`errors`
	`fmt`
	`io/ioutil`
	`log`
	`os`
	`path/filepath`
	`strings`
	`text/template`
	`unicode`

	`github.com/hictl/cmd/internal/common/filez`
	`github.com/hictl/cmd/internal/common/jsonz`
	`github.com/hictl/hictlc/gen`
	`github.com/spf13/cobra`
)

const (
	defaultSchema = "ent/schema"
	genFile       = "package ent\n\n//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema\n"

	schema = "schema"
	mixin  = "mixin"

	hictlHomeDir = ".hictl"
	hictlConfig  = "hictl.json"
)

func InitCmd() *cobra.Command {
	var target string
	var database string
	cmd := &cobra.Command{
		Use:   "init [flags] [schemas]",
		Short: "initialize an environment with zero or more schemas",
		Example: examples(
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
				log.Fatalln(fmt.Errorf("hictl/init: %w", err))
			}
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", defaultSchema, "target directory for schemas")
	cmd.Flags().StringVarP(&database, "database", "d", "", "reverse Engineering database for schemas")

	return cmd
}

type SchemaTemplate struct {
	Name    string
	Fields  []string
	Indexes []string
}

// writeSchema initialize an environment for ent codegen.
func writeSchema(target string, database string, names []string) error {
	if err := createDir(target); err != nil {
		return fmt.Errorf("create dir %s: %w", target, err)
	}
	for _, name := range names {
		if err := gen.ValidSchemaName(name); err != nil {
			return fmt.Errorf("init schema %s: %w", name, err)
		}
		if filez.FileExists(target, strings.ToLower(name+".go")) {
			return fmt.Errorf("init schema %s: already exists", name)
		}
		err := writeSchemaTmpl(target, name)
		if err != nil {
			return err
		}
	}

	err := writeMixinTmpl(target)
	if err != nil {
		return err
	}

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
func writeSchemaTmpl(target string, name string) error {

	st, _ := populateSchemaTemplate(name)

	buffer := bytes.NewBuffer(nil)
	tmplSchema := template.Must(template.New(schema).Parse(schemaTemplate))
	if err := tmplSchema.Execute(buffer, st); err != nil {
		return fmt.Errorf("executing template %s: %w", name, err)
	}
	schemaTarget := filepath.Join(target, strings.ToLower(name+".go"))
	if err := os.WriteFile(schemaTarget, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", schemaTarget, err)
	}
	return nil
}

func populateSchemaTemplate(name string) (*SchemaTemplate, error) {
	home, _ := os.UserHomeDir()
	hictlHome := filepath.Join(home, hictlHomeDir)
	if filez.FileExists(hictlHome, hictlConfig) {
		hictlConfigFile := filepath.Join(hictlHome, strings.ToLower(hictlConfig))
		conf, _ := ioutil.ReadFile(hictlConfigFile)
		hictl := &Hictl{}
		err := jsonz.UnmarshalStruct(conf, hictl)
		if err != nil {
			return nil, err
		}
		pretty, _ := jsonz.Pretty(hictl)
		fmt.Printf("the hictl database config pretty is:%s", pretty)
	}
	st := &SchemaTemplate{
		Name:    name,
		Fields:  []string{"field.Int64(\"id\"),", "field.String(\"order_number\"),"},
		Indexes: []string{},
	}

	return st, nil
}

// writeSchemaTmpl writeSchema the mixin.tmpl
func writeMixinTmpl(target string) error {
	buffer := bytes.NewBuffer(nil)
	tmplMixin := template.Must(template.New(mixin).Parse(mixinTemplate))
	if err := tmplMixin.Execute(buffer, ""); err != nil {
		return fmt.Errorf("executing template %s: %w", "mixin", err)
	}
	mixinTarget := filepath.Join(target, strings.ToLower(mixin+".go"))
	if err := os.WriteFile(mixinTarget, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", mixinTarget, err)
	}
	return nil
}

func examples(ex ...string) string {
	for i := range ex {
		ex[i] = "  " + ex[i]
	}
	return strings.Join(ex, "\n")
}
