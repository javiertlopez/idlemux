package mongodb

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DB struct {
	mongo  *mongo.Database
	logger *logrus.Logger
}

func New(
	l *logrus.Logger,
	m *mongo.Database,
) *DB {
	return &DB{
		mongo:  m,
		logger: l,
	}
}
