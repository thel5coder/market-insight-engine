package model

type EconomicEvent struct {
	Event    string
	Country  string
	Actual   float64
	Estimate float64
	Previous float64
	Impact   int
	Unit     string
	Time     int64
}

func (e *EconomicEvent) Surprise() float64 {
	if e.Estimate == 0 {
		return 0
	}
	return ((e.Actual - e.Estimate) / e.Estimate) * 100
}

type News struct {
	Title       string
	Content     string
	Source      string
	URL         string
	Timestamp   int64
	ImpactLevel ImpactLevel
}

type ImpactLevel string

const (
	IMPACT_LEVEL_LOW     ImpactLevel = "LOW"
	IMPACT_LEVEL_MEDIUM  ImpactLevel = "MEDIUM"
	IMPACT_LEVEL_HIGHT   ImpactLevel = "HIGH"
	IMPACT_LEVEL_UNKNOWN ImpactLevel = "UNKNOWN"
)
