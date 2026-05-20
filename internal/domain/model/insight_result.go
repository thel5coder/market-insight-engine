package model

type InsightResult struct {
	Signal     string `json:"signal"`     // Contoh: "BULLISH", "BEARISH", "NEUTRAL"
	Confidence int    `json:"confidence"` // Persentase keyakinan (0-100)
	Reasoning  string `json:"reasoning"`  // Alasan singkat di balik sinyal (maksimal 2 kalimat)
}
