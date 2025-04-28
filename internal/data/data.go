package data

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
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
	nats  *nats.Conn
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	ctx := context.Background()

	log.Debug("Connecting to MongoDB")
	var mongoDb *mongo.Database
	mongoClient, mongoCleanup, mongoErr := MongoClient(ctx, c, log)
	if mongoErr != nil {
		log.Error(mongoErr, "failed to connect to MongoDB")
	} else {
		dbname := "orka"
		if c.Database.Db != "" {
			dbname = c.Database.Db
		}
		mongoDb = MongoConn(ctx, mongoClient, log, dbname)
		log.Debug("Connected to MongoDB")
	}

	log.Debug("Connecting to Redis")
	redisClient, redisCleanup, redisErr := RedisConn(ctx, c, log)
	if redisErr != nil {
		log.Error(redisErr, "failed to connect to Redis")
	} else {
		log.Debug("Connected to Redis")
	}

	log.Debug("Connecting to NATS")
	var natsClient *nats.Conn
	natsClient, natsCleanup, natsErr := NatsConn(ctx, c, log)
	if natsErr != nil {
		log.Error(natsErr, "failed to connect to NATS")
	} else {
		log.Debug("Connected to NATS")
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
		if natsErr == nil {
			err := natsCleanup()
			if err != nil {
				log.Error(err, "failed to close NATS")
			}
		}
	}

	log.Debug("Connected to all data sources")
	return &Data{
		mongo: mongoClient,
		redis: redisClient,
		db:    mongoDb,
		nats:  natsClient,
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
	_, err := MongoIndexes(ctx, db)
	if err != nil {
		log.Error(err, "failed to create indexes")
	}
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

func NatsConn(ctx context.Context, c *conf.Data, log *log.Helper) (*nats.Conn, func() error, error) {
	if c.GetNats() == nil || c.GetNats().GetAddr() == "" {
		return nil, nil, errors.New("NATS configuration is missing")
	}

	var opts []nats.Option

	if c.GetNats().GetUsername() != "" && c.GetNats().GetPassword() != "" {
		log.Debug("Found NATS credentials, using them")
		opts = append(opts, nats.UserInfo(c.Nats.Username, c.Nats.Password))
	}

	if c.GetNats().GetSubject() != "" {
		log.Debug("Found NATS subject, using it")
		opts = append(opts, nats.Name(c.Nats.Subject))
	}

	natsClient, natsErr := nats.Connect(c.Nats.Addr, opts...)
	if natsErr != nil {
		log.Error(natsErr, "failed to connect to NATS")
		return nil, nil, natsErr
	}

	natsCleanup := func() error {
		return natsClient.Drain()
	}

	return natsClient, natsCleanup, nil
}
