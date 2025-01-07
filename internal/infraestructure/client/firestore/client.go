package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"sync"
)

const (
	projectID      = "quotes-api-100"
	collectionName = "expenses-api"
	documentName   = "config"
)

var (
	clientFirestore *firestore.Client
	once            sync.Once
)

func init() {
	once.Do(func() {
		var err error

		clientFirestore, err = firestore.NewClient(context.Background(), projectID)
		if err != nil {
			log.Fatalf("Error creating Firestore client: %v", err)
		}
	})
}

func GetValue(key string) (string, error) {
	doc, err := clientFirestore.Collection(collectionName).Doc(documentName).Get(context.Background())
	if err != nil {
		return "", err
	}

	config := doc.Data()
	return config[key].(string), nil
}
