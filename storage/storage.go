package storage

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"time"
)

type Storage interface {
	Add(key, value string) error
	Get(key string) (string, error)
	Close()
}

type firestoreStorage struct {
	client  *firestore.Client
	context context.Context
}

func CreateStorage() (Storage, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return firestoreStorage{
		client:  client,
		context: ctx, // TODO: Context should not be saved in structs...
	}, nil
}

func (storage firestoreStorage) Add(key, value string) error {
	data := map[string]interface{}{
		"time":  time.Now(),
		"value": value,
	}
	_, err := storage.client.Collection(key).NewDoc().Set(storage.context, data)
	return err
}

func (storage firestoreStorage) Get(key string) (string, error) {
	return "", errors.New("Implement me!")
}

func (storage firestoreStorage) Close() {
	storage.client.Close()
}
