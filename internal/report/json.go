package report

import (
	"encoding/json"
	"io"

	"github.com/user/envdiff/internal/diff"
)

type jsonReport struct {
	BaseFile    string          `json:"base_file"`
	CompareFile string          `json:"compare_file"`
	Missing     []string        `json:"missing"`
	Extra       []string        `json:"extra"`
	Mismatched  []jsonMismatch  `json:"mismatched"`
	Clean       bool            `json:"clean"`
}

type jsonMismatch struct {
	Key          string `json:"key"`
	BaseValue    string `json:"base_value"`
	CompareValue string `json:"compare_value"`
}

func renderJSON(w io.Writer, result diff.Result, baseFile, compareFile string) error {
	mismatched := make([]jsonMismatch, 0, len(result.Mismatched))
	for _, m := range result.Mismatched {
		mismatched = append(mismatched, jsonMismatch{
			Key:          m.Key,
			BaseValue:    m.BaseValue,
			CompareValue: m.CompareValue,
		})
	}

	missing := result.Missing
	if missing == nil {
		missing = []string{}
	}
	extra := result.Extra
	if extra == nil {
		extra = []string{}
	}

	report := jsonReport{
		BaseFile:    baseFile,
		CompareFile: compareFile,
		Missing:     missing,
		Extra:       extra,
		Mismatched:  mismatched,
		Clean:       len(missing) == 0 && len(extra) == 0 && len(mismatched) == 0,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
