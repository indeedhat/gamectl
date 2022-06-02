package apps

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/indeedhat/gamectl/internal/config"
	"github.com/indeedhat/gamectl/internal/juniper"
	"gorm.io/gorm"
)

const (
	ListKey   = "apps"
	ListUsage = "list apps"
)

func List(*gorm.DB) juniper.CliCommandFunc {
	return func([]string) error {
		var (
			green  = color.New(color.FgGreen).SprintFunc()
			yellow = color.New(color.FgYellow).SprintFunc()
			apps   = *config.Apps()
			rows   = make([]string, 0, len(apps))
			keys   = make([]string, 0, len(apps))
		)

		for key := range apps {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, key := range keys {
			app := apps[key]
			for i, line := range strings.Split(app.Description, "\n") {
				if i != 0 {
					key = "..."
				}

				rows = append(rows, fmt.Sprintf("%s\t%s\n", green(key), yellow(line)))
			}
		}

		tbl := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', tabwriter.StripEscape)
		for _, row := range rows {
			fmt.Fprint(tbl, row)
		}

		tbl.Flush()
		return nil
	}
}
