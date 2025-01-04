package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"sync"
)

var (
	clientFirestore *firestore.Client
	once            sync.Once
)

func init() {
	once.Do(func() {
		var err error

		clientFirestore, err = firestore.NewClient(context.Background(), "quotes-api-100")
		if err != nil {
			log.Fatalf("Error creating Firestore client: %v", err)
		}
	})
}

func GetValue(key string) (string, error) {
	doc, err := clientFirestore.Collection("quotes-api").Doc("config").Get(context.Background())
	if err != nil {
		return "", err
	}

	config := doc.Data()
	return config[key].(string), nil
}
