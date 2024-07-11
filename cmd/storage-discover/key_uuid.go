package main

import (
	"io"
	"os"

	"github.com/google/uuid"
)

var zeroUUID = uuid.UUID{}

func loadUUIDKey() (string, error) {
	var id uuid.UUID
	if f, err := os.Open("eaas-id.txt"); err == nil {
		defer f.Close()
		keyBytes, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		id, err = uuid.ParseBytes(keyBytes)
		if err != nil {
			return "", err
		}
	}
	if id == zeroUUID {
		var err error
		id, err = uuid.NewUUID()
		if err != nil {
			return "", err
		}
		keyBytes, err := id.MarshalBinary()
		if err != nil {
			return "", err
		}
		os.WriteFile("eaas-id.txt", keyBytes, 0o600)
	}
	return id.String(), nil
}
