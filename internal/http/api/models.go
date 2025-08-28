package api

// Models

type UpsertStorage struct {
	Name        string `json:"name" binding:"required"`
	IsMain      bool   `json:"is_main"`
	Address     string `json:"address" binding:"required"`
	Provider    string `json:"provider" binding:"required"`
	Credentials map[string]struct {
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
	} `json:"credentials" binding:"required"`
	IsSecure            bool   `json:"is_secure"`
	DefaultRegion       string `json:"default_region"`
	HealthCheckInterval string `json:"health_check_interval"`
	HttpTimeout         string `json:"http_timeout"`
	RateLimit           struct {
		Enable bool   `json:"enable"`
		Rpm    uint32 `json:"rpm"`
	} `json:"rate_limit"`
}

type CreateReplicationRequest struct {
	User     string   `json:"user" binding:"required"`
	From     string   `json:"from" binding:"required"`
	To       string   `json:"to" binding:"required"`
	Buckets  []string `json:"buckets"`
	ToBucket string   `json:"to_bucket"`
	AgentURL string   `json:"agent_url"`
}

type replicationIdent struct {
	User     string `json:"user" binding:"required"`
	Bucket   string `json:"bucket" binding:"required"`
	From     string `json:"from" binding:"required"`
	To       string `json:"to" binding:"required"`
	ToBucket string `json:"to_bucket"`
}
