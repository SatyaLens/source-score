package pgsql

import (
	"context"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func NewClient(ctx context.Context, dsn string, config *gorm.Config) *Client {
	client := new(Client)
	DB, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatalf("failed to open connection with database: %s", err)
	}

	client.DB = DB

	return client
}

func (client *Client) SetAutoMigration(ctx context.Context, allModels []any) {
	err := client.DB.WithContext(ctx).AutoMigrate(allModels...)
	if err != nil {
		log.Fatalf("failed enable auto migration for all models :: %s", err)
	}
}

func (client *Client) Create(ctx context.Context, record any) *gorm.DB {
	return client.DB.WithContext(ctx).Create(record)
}

func (client *Client) Delete(ctx context.Context, record any) *gorm.DB {
	return client.DB.WithContext(ctx).Delete(record)
}

func (client *Client) FindFirst(ctx context.Context, record any) *gorm.DB {
	return client.DB.WithContext(ctx).First(record)
}

func (client *Client) FindAll(ctx context.Context, records any) *gorm.DB {
	return client.DB.WithContext(ctx).Find(records)
}

func (client *Client) Update(ctx context.Context, record any) *gorm.DB {
	return client.DB.WithContext(ctx).Save(record)
}
