package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/kamunder/kamunder/process"
)

type Chain map[string]process.ProcessInstance
type KeysPath []string

type Label func(process.ProcessInstance) string

func (p KeysPath) KeysOnly(c Chain) string {
	return p.join(c, func(item process.ProcessInstance) string {
		return fmt.Sprint(item.Key)
	}, "\n")
}

func (p KeysPath) StandardLine(c Chain) string {
	return p.join(c, func(item process.ProcessInstance) string {
		var pTag, eTag, vTag string
		if item.ProcessVersion > 0 {
			pTag = fmt.Sprintf(" p:%s", item.ParentKey)
		} else {
			pTag = " p:<root>"
		}
		if item.EndDate != "" {
			eTag = fmt.Sprintf(" e:%s", item.EndDate)
		}
		if item.ProcessVersionTag != "" {
			vTag = "/" + item.ProcessVersionTag
		}

		return fmt.Sprintf(
			"%-16s %s %s v%s%s %s s:%s%s%s i:%t",
			item.Key, item.TenantId, item.BpmnProcessId, version, vTag, item.State, item.StartDate, eTag, pTag, item.Incident,
		)
	}, "\n")
}

func (p KeysPath) PrettyLine(c Chain) string {
	return p.join(c, func(it process.ProcessInstance) string {
		return fmt.Sprintf("%s (%s)", it.Key, it.BpmnProcessId)
	}, " → ")
}

func (p KeysPath) join(c Chain, label Label, sep string) string {
	if len(p) == 0 {
		return ""
	}
	if label == nil {
		label = func(it process.ProcessInstance) string {
			return fmt.Sprintf("%s (%s)", it.Key, it.BpmnProcessId)
		}
	}
	out := make([]string, 0, len(p))
	for _, k := range p {
		out = append(out, label(c[k]))
	}
	return strings.Join(out, sep)
}
