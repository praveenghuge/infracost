package output

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type ReportInput struct {
	Metadata map[string]string
	Root     Root
}

func Load(data []byte) (Root, error) {
	var out Root
	err := json.Unmarshal(data, &out)
	return out, err
}

func Combine(inputs []ReportInput, opts Options) Root {
	var combined Root

	var totalHourlyCost *decimal.Decimal
	var totalMonthlyCost *decimal.Decimal

	projects := make([]ProjectResult, 0)

	for _, input := range inputs {

		projects = append(projects, input.Root.ProjectResults...)

		for _, r := range input.Root.Resources {
			for k, v := range input.Metadata {
				if r.Metadata == nil {
					r.Metadata = make(map[string]string)
				}
				r.Metadata[k] = v
			}

			combined.Resources = append(combined.Resources, r)
		}

		if input.Root.TotalHourlyCost != nil {
			if totalHourlyCost == nil {
				totalHourlyCost = decimalPtr(decimal.Zero)
			}

			totalHourlyCost = decimalPtr(totalHourlyCost.Add(*input.Root.TotalHourlyCost))
		}
		if input.Root.TotalMonthlyCost != nil {
			if totalMonthlyCost == nil {
				totalMonthlyCost = decimalPtr(decimal.Zero)
			}

			totalMonthlyCost = decimalPtr(totalMonthlyCost.Add(*input.Root.TotalMonthlyCost))
		}
	}

	sortResources(combined.Resources, opts.GroupKey)

	combined.Version = outputVersion
	combined.ProjectResults = projects
	combined.TotalHourlyCost = totalHourlyCost
	combined.TotalMonthlyCost = totalMonthlyCost
	combined.TimeGenerated = time.Now()

	return combined
}
