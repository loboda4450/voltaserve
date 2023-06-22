package config

type Config struct {
	Port           int
	APIURL         string
	LanguageURL    string
	ExiftoolURL    string
	FFMPEGURL      string
	ImageMagickURL string
	LibreOfficeURL string
	OCRMyPDFURL    string
	PopplerURL     string
	TesseractURL   string
	Security       SecurityConfig
	Limits         LimitsConfig
	S3             S3Config
}

type SecurityConfig struct {
	APIKey string `json:"api_key"`
}

type LimitsConfig struct {
	ExternalCommandTimeoutSeconds int
	FileProcessingMaxSizeMB       int
	ImagePreviewMaxWidth          int
	ImagePreviewMaxHeight         int
}

type S3Config struct {
	URL       string
	AccessKey string
	SecretKey string
	Region    string
	Secure    bool
}
