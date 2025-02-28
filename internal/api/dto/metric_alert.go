package dto

import (
	"fmt"

	"github.com/sHelllWalker/heimdallr/internal/enums"
)

type MetricAlert struct {
	Event
}

func (ma *MetricAlert) GetTitle() string {
	descriptionTitle, _ := ma.Data["description_title"].(string)

	return descriptionTitle
}

func (ma *MetricAlert) GetDescription() string {
	descriptionText, _ := ma.Data["description_text"].(string)

	rule, _ := ma.Data["metric_alert"].(map[string]any)["alert_rule"].(map[string]any)

	description := ""
	if ruleName, _ := rule["name"].(string); ruleName != "" {
		description = fmt.Sprintf("rule `%s`\n", ruleName)
	}

	return description + descriptionText
}

func (ma *MetricAlert) GetLevel() string {
	return ""
}

func (ma *MetricAlert) GetLink() string {
	link, _ := ma.Data["web_url"].(string)

	return link
}

func (ma *MetricAlert) GetReadableResource() string {
	return "Metric Alert"
}

func (ma *MetricAlert) GetEventLevel() enums.EventLevel {
	switch ma.Action {
	case "resolved":
		return enums.ProblemResolved
	case "critical", "warning":
		return enums.Problem
	}

	return enums.Problem
}
