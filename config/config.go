package config

// Config represents s3url configurations
type Config struct {
	Bucket   string
	Duration int64
	Key      string
	Profile  string
	Upload   string
	Version  bool
}
