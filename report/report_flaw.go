package report

type FlawSummary struct {
	Total                     int `json:"total,omitempty"`
	Fixed                     int `json:"fixed,omitempty"`
	TotalAffectingPolicy      int `json:"total_affecting_policy,omitempty"`
	Mitigated                 int `json:"mitigated,omitempty"`
	OpenAffectingPolicy       int `json:"open_affecting_policy,omitempty"`
	OpenButNotAffectingPolicy int `json:"open_but_not_affecting_policy,omitempty"`
}

type FlawDetails struct {
	ID                      int
	CWE                     int
	AffectsPolicyCompliance bool
	RemediationStatus       string
	MitigationStatus        string
	Mitigation              string
	SourceFile              string
	LineNumber              int
	ModulePath              string
}

func (f *FlawDetails) IsOpen() bool {
	if f.RemediationStatus == "Fixed" {
		return false
	}

	if !(f.MitigationStatus == "none" || f.MitigationStatus == "rejected") {
		return false
	}

	return true
}

func (f *FlawDetails) IsInList(flaws []FlawDetails) bool {
	for _, existingFlaw := range flaws {
		if existingFlaw.ID == f.ID {
			return true
		}
	}

	return false
}
