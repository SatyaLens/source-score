package cnpg

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func NewClient(dbURL string, config *gorm.Config) *Client {
    client := new(Client)
    db, err := gorm.Open(postgres.Open(dbURL), config)
    if err != nil {
        log.Fatalf("failed to open connection with database:%s :: %s", dbURL, err)
    }

    client.DB = db

    return client
}