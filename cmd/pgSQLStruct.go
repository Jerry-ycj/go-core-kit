package cmd

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(pgssCmd)
}

var pgssCmd = &cobra.Command{
	Use:   "pgss",
	Short: "postgres sql to struct",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		sQL2Struct("/Users/ycj/Downloads/demo.sql", "/Users/ycj/Downloads/dest.go")
		println(time.Since(t).Milliseconds())
	},
}

type field struct {
	Name string
	Type string
	Tags []string
}

func sQL2Struct(sqlFile, destFile string) {
	sqls, err := filekit.ReadString(sqlFile)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(sqls, "\n")
	var dest, temp string
	var fields []field
	var table string
	for i := 0; i < len(parts); i++ {
		val := strings.TrimSpace(parts[i])
		if strings.Index(val, "-") == 0 || strings.Index(val, "create index") == 0 || strings.Index(val, "insert") == 0 || strings.Index(val, "update") == 0 {
			continue
		}
		if strings.Index(val, "create table") == 0 {
			table = val[strings.Index(val, "create")+13 : len(val)-1]
			temp = "type " + stringkit.CamelCase(table) + " struct{\n"
			fields = []field{}
		} else if strings.Index(val, ")") == 0 {
			// end
			for _, f := range fields {
				temp += "\t" + stringkit.CamelCase(f.Name) + " " + f.Type + " `" + strings.Join(f.Tags, " ") + "`\n"
			}
			temp += "}\n\n"
			dest += temp
			temp = ""
		} else if temp != "" {
			es := stringkit.Split(val, "[ ,\t]+")
			if es[0] == "primary" {
				// todo
				continue
			}
			f := field{Name: es[0]}
			f.Tags = append(f.Tags, fmt.Sprintf("json:\"%s\" db:\"%s\"", es[0], strings.ToLower(es[0])))
			if arraykit.StringContains(es, "primary") {
				f.Tags = append(f.Tags, fmt.Sprintf("pk:\"true\" tablename:\"%s\"", table))
			}
			switch es[1] {
			case "varchar", "text":
				f.Type = "class.String"
			case "serial":
				f.Type = "int32"
				f.Tags = append(f.Tags, "autoincrement:\"true\"")
			case "bigserial":
				f.Type = "int64"
				f.Tags = append(f.Tags, "autoincrement:\"true\"")
			case "int", "smallint":
				f.Type = "class.Int32"
				if arraykit.StringContains(es, "primary") {
					f.Type = "int32"
				}
			case "bigint":
				f.Type = "class.Int64"
				if arraykit.StringContains(es, "primary") {
					f.Type = "int64"
				}
			case "timestamp", "date":
				f.Type = "class.Time"
			case "jsonb":
				f.Type = "class.MapString"
			case "varchar[]", "text[]":
				f.Type = "class.ArrString"
			case "int[]":
				f.Type = "class.ArrInt"
			case "boolean":
				f.Type = "bool"
			default:
				if strings.Index(es[1], "decimal") == 0 {
					f.Type = "class.Decimal"
				}
			}
			fields = append(fields, f)
		}
	}
	_ = filekit.WriteFile(destFile, []byte(dest))
}
