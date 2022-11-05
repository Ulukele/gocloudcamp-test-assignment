package data

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

type Manager struct {
	client *mongo.Client
}

type MongoConnectionConfig struct {
	User     string
	Password string
	Host     string
	Options  map[string]string
}

type ConfigEntity struct {
	Service string            `bson:"service" json:"service" binding:"required"`
	Version int               `bson:"version" json:"version"`
	Data    map[string]string `bson:"data" json:"data" binding:"required"`
}

func constructOptionsString(options map[string]string) (string, error) {
	if len(options) == 0 {
		return "", nil
	}

	b := new(bytes.Buffer)
	_, err := fmt.Fprint(b, "?")
	if err != nil {
		return "", err
	}

	for key, value := range options {
		_, err := fmt.Fprintf(b, "%s=%s&", key, value)
		if err != nil {
			return "", err
		}
	}

	return strings.TrimRight(b.String(), "&"), nil
}

func NewManager(config MongoConnectionConfig) (*Manager, error) {
	optionsString, err := constructOptionsString(config.Options)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s/%s",
		config.User,
		config.Password,
		config.Host,
		optionsString)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	log.Printf("Connected to %s", uri)

	return &Manager{client: client}, nil
}

func (m *Manager) EndManagerSession() error {
	log.Print("Disconnect mongo")
	return m.client.Disconnect(context.TODO())
}

// CRUD operations

func (m *Manager) ReadAllConfigsVersions(service string) ([]ConfigEntity, error) {
	coll := m.client.Database("configs").Collection("configs")
	filter := bson.D{{"service", service}}
	opts := options.Find().SetSort(bson.D{{"version", -1}})

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return []ConfigEntity{}, err
	}
	var configs []ConfigEntity

	if err := cursor.All(context.TODO(), &configs); err != nil {
		return []ConfigEntity{}, err
	}

	return configs, nil
}

func (m *Manager) ReadConfigWithVersion(service string, version int) (ConfigEntity, error) {
	coll := m.client.Database("configs").Collection("configs")
	filter := bson.D{{"service", service}, {"version", version}}

	var config ConfigEntity
	if err := coll.FindOne(context.TODO(), filter).Decode(&config); err != nil {
		return ConfigEntity{}, err
	}

	return config, nil
}

func (m *Manager) CreateConfig(config ConfigEntity) error {
	coll := m.client.Database("configs").Collection("configs")
	_, err := coll.InsertOne(context.TODO(), config)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) ReadConfig(service string) (ConfigEntity, error) {
	coll := m.client.Database("configs").Collection("configs")
	filter := bson.D{{"service", service}}
	opts := options.FindOne().SetSort(bson.D{{"version", -1}})

	var config ConfigEntity
	if err := coll.FindOne(context.TODO(), filter, opts).Decode(&config); err != nil {
		return ConfigEntity{}, err
	}

	return config, nil
}

func (m *Manager) UpdateConfig(config ConfigEntity) (ConfigEntity, error) {

	configOld, err := m.ReadConfig(config.Service)
	if err != nil {
		return ConfigEntity{}, err
	}

	config.Version = configOld.Version + 1

	return config, m.CreateConfig(config)
}

func (m *Manager) DeleteConfig(service string) ([]ConfigEntity, error) {

	configs, err := m.ReadAllConfigsVersions(service)
	if err != nil {
		return []ConfigEntity{}, nil
	}

	if len(configs) == 0 {
		return []ConfigEntity{}, fmt.Errorf("no such service: %s", service)
	}

	coll := m.client.Database("configs").Collection("configs")
	filter := bson.D{{"service", service}}

	_, err = coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		return []ConfigEntity{}, nil
	}
	return configs, nil
}
