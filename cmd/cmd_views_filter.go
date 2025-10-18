package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func printFilter(cmd *cobra.Command) {
	var filters []string
	if flagPIParentKey != "" {
		filters = append(filters, fmt.Sprintf("parent-key=%s", flagPIParentKey))
	}
	if flagPIState != "" && flagPIState != "all" {
		filters = append(filters, fmt.Sprintf("state=%s", flagPIState))
	}
	if flagPIParentsOnly {
		filters = append(filters, "parents-only=true")
	}
	if flagPIChildrenOnly {
		filters = append(filters, "children-only=true")
	}
	if flagPIOrphanParentsOnly {
		filters = append(filters, "orphan-parents-only=true")
	}
	if flagPIIncidentsOnly {
		filters = append(filters, "incidents-only=true")
	}
	if flagPINoIncidentsOnly {
		filters = append(filters, "no-incidents-only=true")
	}
	if len(filters) > 0 {
		cmd.Println("filter: " + strings.Join(filters, ", "))
	}
}
