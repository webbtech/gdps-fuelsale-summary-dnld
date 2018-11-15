package config

// defaults struct
type defaults struct {
	AWSRegion       string `yaml:"AWSRegion"`
	S3Bucket        string `yaml:"S3Bucket"`
	S3FilePrefix    string `yaml:"S3FilePrefix"`
	CognitoClientID string `yaml:"CognitoClientID"`
	CognitoPoolID   string `yaml:"CognitoPoolID"`
	CognitoRegion   string `yaml:"CognitoRegion"`
	GraphqlURI      string `yaml:"GraphqlURI"`
	SsmPath         string `yaml:"SsmPath"`
	Stage           string `yaml:"Stage"`
}

type config struct {
	AWSRegion       string
	S3Bucket        string
	S3FilePrefix    string
	CognitoClientID string
	CognitoPoolID   string
	CognitoRegion   string
	GraphqlURI      string
	Stage           StageEnvironment
}
