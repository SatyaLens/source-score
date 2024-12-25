package cnpg

import (
	"context"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(ctx context.Context, dbURL string, config *gorm.Config) *Client {
	client := new(Client)
	db, err := gorm.Open(postgres.Open(dbURL), config)
	if err != nil {
		log.Fatalf("failed to open connection with database:%s :: %s", dbURL, err)
	}

	client.db = db

	return client
}

func (client *Client) SetAutoMigration(ctx context.Context, allModels []interface{}) {
	err := client.db.AutoMigrate(allModels...)
	if err != nil {
		log.Fatalf("failed enable auto migration for all models :: %s", err)
	}
}

func (client *Client) Create(ctx context.Context, record interface{}) *gorm.DB {
	return client.db.Create(record)
}

func (client *Client) Delete(ctx context.Context, record interface{}) *gorm.DB {
	return client.db.Delete(record)
}

func (client *Client) FindByPrimaryKey(ctx context.Context, record interface{}, primaryKey interface{}) *gorm.DB {
	return client.db.First(record, primaryKey)
}

func (client *Client) Update(ctx context.Context, record interface{}) *gorm.DB {
	return client.db.Save(record)
}