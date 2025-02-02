package cnpg

import (
	"context"
	"log"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func NewClient(ctx context.Context, dbURL string, config *gorm.Config) *Client {
	client := new(Client)
	DB, err := gorm.Open(postgres.Open(dbURL), config)
	if err != nil {
		log.Fatalf("failed to open connection with database:%s :: %s", dbURL, err)
	}

	client.DB = DB

	return client
}

func (client *Client) SetAutoMigration(ctx context.Context, allModels []interface{}) {
	err := client.DB.AutoMigrate(allModels...)
	if err != nil {
		log.Fatalf("failed enable auto migration for all models :: %s", err)
	}
}

func (client *Client) Create(ctx context.Context, record interface{}) *gorm.DB {
	return client.DB.Create(record)
}

func (client *Client) Delete(ctx context.Context, record interface{}) *gorm.DB {
	return client.DB.Delete(record)
}

func (client *Client) FindFirst(ctx context.Context, record interface{}) *gorm.DB {
	return client.DB.First(record)
}

func (client *Client) Update(ctx context.Context, record interface{}) *gorm.DB {
	return client.DB.Save(record)
}
