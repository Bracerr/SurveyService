package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostgresDB struct {
	*pgxpool.Pool
}

type MongoDB struct {
	*mongo.Client
}

func NewPostgresDBWrapper(pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{Pool: pool}
}

func NewMongoDBWrapper(client *mongo.Client) *MongoDB {
	return &MongoDB{Client: client}
}
