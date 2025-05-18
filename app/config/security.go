package config

// AwsConfig - AwsConfig statuses
type AwsConfig struct {
	Activate           string
	AccessKey          string
	SecreteAccessKey   string
	Region             string
	DocumentBucketName string
}

type SchedulerConfig struct {
	MaxRetries     int
	MaxDelay       int
	SchedulerCrons map[string]string
}
