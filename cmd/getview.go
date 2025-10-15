package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/spf13/cobra"
)

func listKeyOnlyProcessInstancesView(cmd *cobra.Command, resp process.ProcessInstances) error {
	return renderListViewV(cmd, resp, func(r process.ProcessInstances) []process.ProcessInstance {
		return r.Items
	}, keyOnlyProcessInstanceView)
}

func listProcessInstancesView(cmd *cobra.Command, resp process.ProcessInstances) error {
	if flagOneLine {
		return renderListViewV(cmd, resp, func(r process.ProcessInstances) []process.ProcessInstance {
			return r.Items
		}, oneLineProcessInstanceView)
	}
	return listJSONViewV(cmd, resp, func(r process.ProcessInstances) []process.ProcessInstance {
		return r.Items
	})
}

func keyOnlyProcessInstanceView(cmd *cobra.Command, item process.ProcessInstance) error {
	cmd.Println(item.Key)
	return nil
}

func processInstanceView(cmd *cobra.Command, item process.ProcessInstance) error {
	if flagOneLine {
		return oneLineProcessInstanceView(cmd, item)
	}
	if flagPDKeysOnly {
		return keyOnlyProcessInstanceView(cmd, item)
	}
	cmd.Println(ToJSONString(item))
	return nil
}

func oneLineProcessInstanceView(cmd *cobra.Command, item process.ProcessInstance) error {
	var pTag, eTag, vTag string
	if item.ParentKey > 0 {
		pTag = fmt.Sprintf(" p:%d", item.ParentKey)
	} else {
		pTag = " p:<root>"
	}
	if item.EndDate != "" {
		eTag = fmt.Sprintf(" e:%s", item.EndDate)
	}
	if item.ProcessVersionTag != "" {
		vTag = "/" + item.ProcessVersionTag
	}

	out := fmt.Sprintf(
		"%-16d %s %s v%d%s %s s:%s%s%s i:%t",
		item.Key, item.TenantId, item.BpmnProcessId, item.ProcessVersion, vTag, item.State, item.StartDate, eTag, pTag, item.Incident,
	)
	cmd.Println(strings.TrimSpace(out))
	return nil
}

func listKeyOnlyProcessDefinitionsView(cmd *cobra.Command, resp process.ProcessDefinitions) error {
	return renderListViewV(cmd, resp, func(r process.ProcessDefinitions) []process.ProcessDefinition {
		return r.Items
	}, keyOnlyProcessDefinitionView)
}

func listProcessDefinitionsView(cmd *cobra.Command, resp process.ProcessDefinitions) error {
	if flagOneLine {
		return renderListViewV(cmd, resp, func(r process.ProcessDefinitions) []process.ProcessDefinition {
			return r.Items
		}, oneLineProcessDefinitionView)
	}
	return listJSONViewV(cmd, resp, func(r process.ProcessDefinitions) []process.ProcessDefinition {
		return r.Items
	})
}

func keyOnlyProcessDefinitionView(cmd *cobra.Command, item process.ProcessDefinition) error {
	cmd.Println(item.Key)
	return nil
}

func processDefinitionView(cmd *cobra.Command, item process.ProcessDefinition) error {
	if flagOneLine {
		return oneLineProcessDefinitionView(cmd, item)
	}
	if flagPDKeysOnly {
		return keyOnlyProcessDefinitionView(cmd, item)
	}
	cmd.Println(ToJSONString(item))
	return nil
}

func oneLineProcessDefinitionView(cmd *cobra.Command, item process.ProcessDefinition) error {
	vTag := ""
	if item.VersionTag != "" {
		vTag = "/" + item.VersionTag
	}
	out := fmt.Sprintf("%-16d %s %s v%s%s",
		item.Key, item.TenantId, item.BpmnProcessId, version, vTag,
	)
	cmd.Println(strings.TrimSpace(out))
	return nil
}

//nolint:unused
func listJSONView[Resp any, Item any](cmd *cobra.Command, resp *Resp, itemsOf func(*Resp) *[]Item) error {
	if resp == nil {
		cmd.Println("{}")
		return nil
	}
	printFound(cmd, itemsOf(resp))
	cmd.Println(ToJSONString(resp))
	return nil
}

func listJSONViewV[Resp any, Item any](cmd *cobra.Command, resp Resp, itemsOf func(Resp) []Item) error {
	items := itemsOf(resp)
	printFoundV(cmd, items)
	cmd.Println(ToJSONString(resp))
	return nil
}

//nolint:unused
func renderListView[Resp any, Item any](cmd *cobra.Command, resp *Resp, itemsOf func(*Resp) *[]Item,
	render func(*cobra.Command, *Item) error) error {
	if resp == nil {
		return nil
	}
	itemsPtr := itemsOf(resp)
	if itemsPtr == nil {
		cmd.Println("found: 0")
		return nil
	}
	items := *itemsPtr
	cmd.Println("found:", len(items))
	for i := range items {
		if err := render(cmd, &items[i]); err != nil {
			return err
		}
	}
	return nil
}

func renderListViewV[Resp any, Item any](cmd *cobra.Command, resp Resp, itemsOf func(Resp) []Item,
	render func(*cobra.Command, Item) error) error {
	items := itemsOf(resp)
	if len(items) == 0 {
		cmd.Println("found: 0")
		return nil
	}
	cmd.Println("found:", len(items))
	for _, it := range items {
		if err := render(cmd, it); err != nil {
			return err
		}
	}
	return nil
}

//nolint:unused
func printFound[T any](cmd *cobra.Command, items *[]T) {
	if items == nil {
		cmd.Println("found: 0")
		return
	}
	cmd.Println("found:", len(*items))
}

func printFoundV[T any](cmd *cobra.Command, items []T) {
	cmd.Println("found:", len(items))
}

func printFilter(cmd *cobra.Command) {
	var filters []string
	if flagPIParentKey != 0 {
		filters = append(filters, fmt.Sprintf("parent-key=%d", flagPIParentKey))
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

//nolint:unused
func valueOr[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
