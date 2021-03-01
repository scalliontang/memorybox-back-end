package main

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

const (
	BUCKET_NAME = "Bucket_Name"
)

func saveToGCS(r io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	object := client.Bucket(BUCKET_NAME).Object(objectName)
	wc := object.NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	return attrs.MediaLink, nil
}
