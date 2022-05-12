package config

type Config struct {
	// HideBanner (optional) if true, we don't display the go-feature-flag relay proxy banner
	HideBanner bool `json:"hideBanner"`

	// EnableSwagger (optional) to have access to the swagger
	EnableSwagger bool `json:"enableSwagger"`

	// Debug (optional) if true, go-feature-flag relay proxy will run on debug mode, with more logs and custom responses
	Debug bool `json:"debug"`

	// StartWithRetrieverError (optional) If true, the SDK will start even if we did not get any flags from the retriever.
	// It will serve only default values until the retriever returns the flags.
	// The init method will not return any error if the flag file is unreachable.
	// Default: false
	StartWithRetrieverError bool
}
