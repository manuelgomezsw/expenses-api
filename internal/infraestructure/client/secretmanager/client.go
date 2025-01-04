package secretmanager

import (
	gcpsecretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	clientSecretManager *gcpsecretmanager.Client
	once                sync.Once
)

func init() {
	once.Do(func() {
		var err error

		clientSecretManager, err = gcpsecretmanager.NewClient(context.Background())
		if err != nil {
			log.Fatalf("Error creating Secret Manager client: %v", err)
		}
	})
}

func GetValue(secretName string) (string, error) {
	secretPath := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", os.Getenv("PROJECT_ID"), secretName)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}
	result, err := clientSecretManager.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("error getting secret: %v", err)
	}

	return string(result.Payload.Data), nil
}
