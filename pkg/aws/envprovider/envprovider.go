package envprovider

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// EnvProviderName provides a name of Env provider
const EnvProviderName = "DiggerEnvProvider"

var (
	// ErrAccessKeyIDNotFound is returned when the AWS Access Key ID can't be
	// found in the process's environment.
	ErrAccessKeyIDNotFound = awserr.New("EnvAccessKeyNotFound", "DIGGER_AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY not found in environment", nil)

	// ErrSecretAccessKeyNotFound is returned when the AWS Secret Access Key
	// can't be found in the process's environment.
	ErrSecretAccessKeyNotFound = awserr.New("EnvSecretNotFound", "DIGGER_AWS_SECRET_ACCESS_KEY or AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY not found in environment", nil)
)

// A EnvProvider retrieves credentials from the environment variables of the
// running process. Environment credentials never expire.
//
// Environment variables used:
//
// * Access Key ID:     DIGGER_AWS_ACCESS_KEY_ID or DIGGER_AWS_ACCESS_KEY, AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
//
// * Secret Access Key: DIGGER_AWS_SECRET_ACCESS_KEY or DIGGER_AWS_SECRET_KEY, AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
type EnvProvider struct {
	retrieved bool
}

// Retrieve retrieves the keys from the environment.
func (e *EnvProvider) Retrieve() (credentials.Value, error) {
	e.retrieved = false

	//assign id from env vars
	idEnvVars := []string{"DIGGER_AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID"}
	id, err := assignEnv(idEnvVars)
	if err != nil {
		return credentials.Value{ProviderName: EnvProviderName}, ErrAccessKeyIDNotFound
	}

	//assign secret from env vars
	secretEnvVars := []string{"DIGGER_AWS_SECRET_ACCESS_KEY", "AWS_SECRET_KEY"}
	secret, err := assignEnv(secretEnvVars)
	if err != nil {
		return credentials.Value{ProviderName: EnvProviderName}, ErrSecretAccessKeyNotFound
	}

	e.retrieved = true
	return credentials.Value{
		AccessKeyID:     id,
		SecretAccessKey: secret,
		SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
		ProviderName:    EnvProviderName,
	}, nil
}

// Assign first non-nil env var
func assignEnv(envVars []string) (string, error) {
	var v string
	for _, envVar := range envVars {
		if value, ok := os.LookupEnv(envVar); ok {
			v = value
			return v, nil
		}
	}
	return "", errors.New("not found")
}

// IsExpired returns if the credentials have been retrieved.
func (e *EnvProvider) IsExpired() bool {
	return !e.retrieved
}
