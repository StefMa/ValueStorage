package storage

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"time"
)

type Storage interface {
	Add(key, value string) error
	Get(key string) (map[time.Time]string, error)
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

func (storage firestoreStorage) Get(key string) (map[time.Time]string, error) {
	docs, err := storage.client.Collection(key).DocumentRefs(storage.context).GetAll()
	if err != nil {
		return nil, nil
	}

	allData := make(map[time.Time]string)
	for _, doc := range docs {
		docsnap, err := doc.Get(storage.context)
		if err != nil {
			continue
		}
		data := docsnap.Data()
		time := data["time"].(time.Time)
		value := data["value"].(string)
		allData[time] = value
	}
	return allData, nil
}

func (storage firestoreStorage) Close() {
	storage.client.Close()
}
