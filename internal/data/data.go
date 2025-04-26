package data

import (
	"context"
	"errors"

	"github.com/orka-org/orkacore/internal/conf"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// DataProviderSet is data providers.
var DataProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	mongo *mongo.Client
	redis *redis.Client
	db    *mongo.Database
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	ctx := context.Background()

	var mongoDb *mongo.Database
	mongoClient, mongoCleanup, mongoErr := MongoClient(ctx, c, log)
	if mongoErr != nil {
		log.Error(mongoErr, "failed to connect to MongoDB")
	}
	dbname := "orka"
	if c.Database.Db != "" {
		dbname = c.Database.Db
	}
	mongoDb = MongoConn(ctx, mongoClient, log, dbname)

	redisClient, redisCleanup, redisErr := RedisConn(ctx, c, log)
	if redisErr != nil {
		log.Error(redisErr, "failed to connect to Redis")
	}

	cleanup := func() {
		log.Info("closing the data resources")
		if mongoErr == nil {
			err := mongoCleanup()
			if err != nil {
				log.Error(err, "failed to close MongoDB")
			}
		}
		if redisErr == nil {
			err := redisCleanup()
			if err != nil {
				log.Error(err, "failed to close Redis")
			}
		}
	}

	return &Data{
		mongo: mongoClient,
		redis: redisClient,
		db:    mongoDb,
	}, cleanup, nil
}

func MongoClient(ctx context.Context, c *conf.Data, log *log.Helper) (*mongo.Client, func() error, error) {
	if c.Database == nil || c.Database.Url == "" {
		err := errors.New("MongoDB configuration is missing")
		log.Error(err, "failed to connect to MongoDB")
		return nil, nil, err
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(c.Database.Url).SetServerAPIOptions(serverAPI)
	if c.GetDatabase().GetUsername() != "" && c.GetDatabase().GetPassword() != "" {
		opts.SetAuth(options.Credential{
			Username: c.Database.Username,
			Password: c.Database.Password,
		})
	}
	mongoClient, mongoErr := mongo.Connect(ctx, opts)
	if mongoErr != nil {
		log.Error(mongoErr, "failed to connect to MongoDB")
	}
	mongoCleanup := func() error {
		return mongoClient.Disconnect(ctx)
	}
	return mongoClient, mongoCleanup, mongoErr
}

func MongoConn(ctx context.Context, client *mongo.Client, log *log.Helper, dbname string) *mongo.Database {
	db := client.Database(dbname, &options.DatabaseOptions{})
	return db
}

func RedisConn(ctx context.Context, c *conf.Data, log *log.Helper) (*redis.Client, func() error, error) {
	if c.Redis == nil || c.Redis.Addr == "" {
		return nil, nil, errors.New("Redis configuration is missing")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
	})
	redisCleanup := func() error {
		return redisClient.Close()
	}

	return redisClient, redisCleanup, nil
}
