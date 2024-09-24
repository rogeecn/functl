package cmd

import (
	"fmt"
	"strings"

	. "functl/config"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	pg "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func genModelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "gen model",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := Global.Gen.Model
			if conf.Schema == "" {
				conf.Schema = "public"
			}

			return postgres.GenerateDSN(
				conf.DSN,
				conf.Schema,
				conf.Path,
				template.Default(pg.Dialect).
					UseSchema(func(schema metadata.Schema) template.Schema {
						return template.
							DefaultSchema(schema).
							UseModel(
								template.
									DefaultModel().
									UseTable(func(table metadata.Table) template.TableModel {
										if lo.Contains(conf.Ignores, table.Name) {
											table := template.DefaultTableModel(table)
											table.Skip = true
											return table
										}

										return template.DefaultTableModel(table).UseField(func(column metadata.Column) template.TableModelField {
											defaultTableModelField := template.DefaultTableModelField(column)
											defaultTableModelField = defaultTableModelField.UseTags(fmt.Sprintf(`json:"%s"`, column.Name))

											if schema.Name != conf.Schema {
												return defaultTableModelField
											}

											fields, ok := conf.Types[table.Name]
											if !ok {
												return defaultTableModelField
											}

											toType, ok := fields[column.Name]
											if !ok {
												return defaultTableModelField
											}

											splits := strings.Split(toType, ".")
											typeName := splits[len(splits)-1]

											pkgSplits := strings.Split(splits[0], "/")
											typePkg := pkgSplits[len(pkgSplits)-1]

											defaultTableModelField = defaultTableModelField.
												UseType(template.Type{
													Name:       fmt.Sprintf("%s.%s", typePkg, typeName),
													ImportPath: splits[0],
												})

											return defaultTableModelField
										})
									}),
							)
					}),
			)
		},
	}

	return cmd
}
