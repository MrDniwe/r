package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func NewRepository(client *mongo.Client, l *log.Logger) *ArcticleRepo {
	// пингуем на всякий случай
	err := client.Ping(context.Background(), nil)
	if err != nil {
		l.Fatal(err)
	}
	return &ArcticleRepo{client, l}
}

type ArcticleRepo struct {
	db *mongo.Client
	L  *log.Logger
}