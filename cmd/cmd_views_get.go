package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/kamunder/kamunder/process"
	"github.com/spf13/cobra"
)

type RenderMode int

const (
	ModeJSON RenderMode = iota
	ModeOneLine
	ModeKeysOnly
)

// ----- Common single-item and list renderers -----

func itemView[Item any](cmd *cobra.Command, item Item, mode RenderMode, oneLine func(Item) string, keyOf func(Item) string) error {
	switch mode {
	case ModeJSON:
		cmd.Println(ToJSONString(item))
	case ModeKeysOnly:
		cmd.Println(keyOf(item))
	default:
		cmd.Println(strings.TrimSpace(oneLine(item)))
	}
	return nil
}

func listOrJSON[Resp any, Item any](
	cmd *cobra.Command,
	resp Resp,
	items []Item,
	mode RenderMode,
	oneLine func(Item) string,
	keyOf func(Item) string,
) error {
	if len(items) == 0 {
		cmd.Println("found: 0")
		if mode == ModeJSON {
			cmd.Println(ToJSONString(resp))
		}
		return nil
	}
	cmd.Println("found:", len(items))
	switch mode {
	case ModeJSON:
		cmd.Println(ToJSONString(resp))
	case ModeKeysOnly:
		for _, it := range items {
			cmd.Println(keyOf(it))
		}
	default: // ModeOneLine
		for _, it := range items {
			cmd.Println(strings.TrimSpace(oneLine(it)))
		}
	}
	return nil
}

// ----- Process Instances (single and list) -----

func processInstanceView(cmd *cobra.Command, item process.ProcessInstance) error {
	switch {
	case flagAsJson:
		return itemView(cmd, item, ModeJSON, oneLinePI, func(it process.ProcessInstance) string { return it.Key })
	case flagPIKeysOnly:
		return itemView(cmd, item, ModeKeysOnly, oneLinePI, func(it process.ProcessInstance) string { return it.Key })
	default:
		return itemView(cmd, item, ModeOneLine, oneLinePI, func(it process.ProcessInstance) string { return it.Key })
	}
}

func listProcessInstancesView(cmd *cobra.Command, resp process.ProcessInstances) error {
	switch {
	case flagAsJson:
		return listOrJSON(cmd, resp, resp.Items, ModeJSON, oneLinePI, func(it process.ProcessInstance) string { return it.Key })
	case flagPIKeysOnly:
		return listOrJSON(cmd, resp, resp.Items, ModeKeysOnly, oneLinePI, func(it process.ProcessInstance) string { return it.Key })
	default:
		return listOrJSON(cmd, resp, resp.Items, ModeOneLine, oneLinePI, func(it process.ProcessInstance) string { return it.Key })
	}
}

func oneLinePI(it process.ProcessInstance) string {
	pTag := " p:<root>"
	if it.ParentKey != "" {
		pTag = " p:" + it.ParentKey
	}
	eTag := ""
	if it.EndDate != "" {
		eTag = " e:" + it.EndDate
	}
	vTag := ""
	if it.ProcessVersionTag != "" {
		vTag = "/" + it.ProcessVersionTag
	}
	return fmt.Sprintf(
		"%-16s %s %s v%d%s %s s:%s %s%s i:%t",
		it.Key, it.TenantId, it.BpmnProcessId, it.ProcessVersion, vTag,
		it.State, it.StartDate, eTag, pTag, it.Incident,
	)
}

// ----- Process Definitions (single and list) -----

func processDefinitionView(cmd *cobra.Command, item process.ProcessDefinition) error {
	switch {
	case flagAsJson:
		return itemView(cmd, item, ModeJSON, oneLinePD, func(it process.ProcessDefinition) string { return it.Key })
	case flagPDKeysOnly:
		return itemView(cmd, item, ModeKeysOnly, oneLinePD, func(it process.ProcessDefinition) string { return it.Key })
	default:
		return itemView(cmd, item, ModeOneLine, oneLinePD, func(it process.ProcessDefinition) string { return it.Key })
	}
}

func listProcessDefinitionsView(cmd *cobra.Command, resp process.ProcessDefinitions) error {
	switch {
	case flagAsJson:
		return listOrJSON(cmd, resp, resp.Items, ModeJSON, oneLinePD, func(it process.ProcessDefinition) string { return it.Key })
	case flagPDKeysOnly:
		return listOrJSON(cmd, resp, resp.Items, ModeKeysOnly, oneLinePD, func(it process.ProcessDefinition) string { return it.Key })
	default:
		return listOrJSON(cmd, resp, resp.Items, ModeOneLine, oneLinePD, func(it process.ProcessDefinition) string { return it.Key })
	}
}

func oneLinePD(it process.ProcessDefinition) string {
	vTag := ""
	if it.VersionTag != "" {
		vTag = "/" + it.VersionTag
	}
	return fmt.Sprintf("%-16s %s %s v%d%s",
		it.Key, it.TenantId, it.BpmnProcessId, it.Version, vTag,
	)
}

func pickMode() RenderMode {
	switch {
	case flagAsJson:
		return ModeJSON
	case flagPIKeysOnly: // reuse PI keys-only for paths
		return ModeKeysOnly
	default:
		return ModeOneLine
	}
}
