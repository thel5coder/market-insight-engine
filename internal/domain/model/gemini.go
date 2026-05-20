package model

type GeminiRequest struct {
	Contents          []GeminiContent  `json:"contents"`
	SystemInstruction GeminiContent    `json:"systemInstruction"`
	GenerationConfig  GenerationConfig `json:"generationConfig"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	ResponseMIMEType string  `json:"responseMimeType"`
	Temperature      float64 `json:"temperature"`
}
