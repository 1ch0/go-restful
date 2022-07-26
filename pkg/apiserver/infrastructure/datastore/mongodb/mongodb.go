package mongodb

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/x/bsonx"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/1ch0/go-restful/pkg/apiserver/infrastructure/datastore"
	"github.com/1ch0/go-restful/pkg/apiserver/utils/log"
)

type mongodb struct {
	client   *mongo.Client
	database string
}

// PrimaryKey primary key
const PrimaryKey = "_name"

// New new mongodb datastore instance
func New(ctx context.Context, cfg datastore.Config) (datastore.DataStore, error) {
	if !strings.HasPrefix(cfg.URL, "mongodb://") {
		cfg.URL = fmt.Sprintf("mongodb://%s", cfg.URL)
	}
	clientOpts := options.Client().ApplyURI(cfg.URL)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	log.Logger.Infof("connected to mongodb: %s", hidePass(cfg.URL))
	m := &mongodb{
		client:   client,
		database: cfg.Database,
	}
	return m, err
}

func hidePass(str string) string {
	reg := regexp.MustCompile(`(^mongodb://.+?:)(.+)(@.+$)`)
	return reg.ReplaceAllString(str, `${1}xxx${3}`)
}

// Add add data model
func (m *mongodb) Add(ctx context.Context, entity datastore.Entity) error {
	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}
	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}
	entity.SetCreateTime(time.Now())
	if err := m.Get(ctx, entity); err != nil {
		return datastore.ErrRecordExist
	}
	model, err := converToMap(entity)
	if err != nil {
		return datastore.ErrEntityInvalid
	}
	model[PrimaryKey] = entity.PrimaryKey()
	collection := m.client.Database(m.database).Collection(entity.TableName())
	_, err = collection.InsertOne(ctx, model)
	if err != nil {
		return datastore.NewDBError(err)
	}
	return nil
}

// BatchAdd will adds batched entities to database, Name() and TableName() can't return zero value.
func (m *mongodb) BatchAdd(ctx context.Context, entities []datastore.Entity) error {
	notRollback := make(map[string]int)
	for i, saveEntity := range entities {
		if err := m.Add(ctx, saveEntity); err != nil {
			if errors.Is(err, datastore.ErrRecordExist) {
				notRollback[saveEntity.PrimaryKey()] = 1
			}
			for _, deleteEntity := range entities[:i] {
				if _, exit := notRollback[deleteEntity.PrimaryKey()]; !exit {
					if err := m.Delete(ctx, deleteEntity); err != nil {
						if !errors.Is(err, datastore.ErrRecordNotExist) {
							log.Logger.Errorf("rollback delete entity failure %w", err)
						}
					}
				}
			}
			return datastore.NewDBError(fmt.Errorf("save entities occur error, %w", err))
		}
	}
	return nil
}

// Put will update entity to database, Name() and TableName() can't return zero value.
func (m *mongodb) Put(ctx context.Context, entity datastore.Entity) error {
	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}
	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}
	entity.SetUpdateTime(time.Now())
	collection := m.client.Database(m.database).Collection(entity.TableName())
	_, err := collection.UpdateOne(ctx, makeNameFilter(entity.PrimaryKey()), makeEntityUpdate(entity))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return datastore.ErrRecordNotExist
		}
		return datastore.NewDBError(err)
	}
	return nil
}

// Delete entity from database, Name() and TableName() can't return zero value.
func (m *mongodb) Delete(ctx context.Context, entity datastore.Entity) error {
	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}
	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}
	if err := m.Get(ctx, entity); err != nil {
		return err
	}
	// check entity is exist
	collection := m.client.Database(m.database).Collection(entity.TableName())
	// delete at most one document in which the "name" field is "Bob" or "bob"
	// specify the SetCollation option to provide a collation that will ignore case for string comparisons
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	_, err := collection.DeleteOne(ctx, makeNameFilter(entity.PrimaryKey()), opts)
	if err != nil {
		log.Logger.Errorf("delete document failure %w", err)
		return datastore.NewDBError(err)
	}
	return nil
}

// Get entity from database, Name() and TableName() can't return zero value.
func (m *mongodb) Get(ctx context.Context, entity datastore.Entity) error {
	if entity.PrimaryKey() == "" {
		return datastore.ErrPrimaryEmpty
	}
	if entity.TableName() == "" {
		return datastore.ErrTableNameEmpty
	}
	collection := m.client.Database(m.database).Collection(entity.TableName())
	if err := collection.FindOne(ctx, makeNameFilter(entity.PrimaryKey())).Decode(entity); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return datastore.ErrRecordNotExist
		}
		return datastore.NewDBError(err)
	}
	return nil
}

// List entities from database, TableName() can't return zero value, if no matches, it will return a zero list without error.
func (m *mongodb) List(ctx context.Context, entity datastore.Entity, op *datastore.ListOptions) ([]datastore.Entity, error) {
	if entity.TableName() == "" {
		return nil, datastore.ErrTableNameEmpty
	}
	collection := m.client.Database(m.database).Collection(entity.TableName())
	// bson.D{{}} specifies 'all documents'
	filter := bson.D{}
	if entity.Index() != nil {
		for k, v := range entity.Index() {
			filter = append(filter, bson.E{
				Key:   strings.ToLower(k),
				Value: v,
			})
		}
	}
	if op != nil {
		filter = _applyFilterOptions(filter, op.FilterOptions)
	}
	var findOptions options.FindOptions
	if op != nil && op.PageSize > 0 && op.Page > 0 {
		findOptions.SetSkip(int64(op.PageSize * (op.Page - 1)))
		findOptions.SetLimit(int64(op.PageSize))
	}
	if op != nil && len(op.SortBy) > 0 {
		_d := bson.D{}
		for _, sortOp := range op.SortBy {
			key := strings.ToLower(sortOp.Key)
			if key == "createtime" || key == "updatetime" {
				key = "basemodel." + key
			}
			_d = append(_d, bson.E{Key: key, Value: int(sortOp.Order)})
		}
		findOptions.SetSort(_d)
	}
	cur, err := collection.Find(ctx, filter, &findOptions)
	if err != nil {
		return nil, datastore.NewDBError(err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Logger.Warnf("close mongodb cursor failure %s", err.Error())
		}
	}()
	var list []datastore.Entity
	for cur.Next(ctx) {
		item, err := datastore.NewEntity(entity)
		if err != nil {
			return nil, datastore.NewDBError(err)
		}
		if err := cur.Decode(item); err != nil {
			return nil, datastore.NewDBError(fmt.Errorf("decode entity failure %$w", err))
		}
		list = append(list, item)
	}
	if err := cur.Err(); err != nil {
		return nil, datastore.NewDBError(err)
	}
	return list, nil
}

func _applyFilterOptions(filter bson.D, filterOptions datastore.FilterOptions) bson.D {
	for _, queryOp := range filterOptions.Queries {
		filter = append(filter, bson.E{Key: strings.ToLower(queryOp.Key), Value: bsonx.Regex(".*"+queryOp.Query+".*", "s")})
	}
	for _, queryOp := range filterOptions.In {
		filter = append(filter, bson.E{Key: strings.ToLower(queryOp.Key), Value: bson.D{bson.E{Key: "$in", Value: queryOp.Values}}})
	}
	for _, queryOp := range filterOptions.IsNotExist {
		filter = append(filter, bson.E{Key: strings.ToLower(queryOp.Key), Value: bson.D{bson.E{Key: "$eq", Value: ""}}})
	}
	return filter
}

// Count entities from database, TableName() can't return zero value.
func (m *mongodb) Count(ctx context.Context, entity datastore.Entity, filterOptions *datastore.FilterOptions) (int64, error) {
	if entity.TableName() == "" {
		return 0, datastore.ErrTableNameEmpty
	}
	collection := m.client.Database(m.database).Collection(entity.TableName())
	filter := bson.D{}
	if entity.Index() != nil {
		for k, v := range entity.Index() {
			filter = append(filter, bson.E{
				Key:   strings.ToLower(k),
				Value: v,
			})
		}
	}
	if filterOptions != nil {
		filter = _applyFilterOptions(filter, *filterOptions)
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, datastore.NewDBError(err)
	}
	return count, nil
}

// IsExist Name() and TableName() can't return zero value.
func (m *mongodb) IsExist(ctx context.Context, entity datastore.Entity) (bool, error) {
	if entity.PrimaryKey() == "" {
		return false, datastore.ErrPrimaryEmpty
	}
	if entity.TableName() == "" {
		return false, datastore.ErrTableNameEmpty
	}
	entity.SetUpdateTime(time.Now())
	collection := m.client.Database(m.database).Collection(entity.TableName())
	err := collection.FindOne(ctx, makeNameFilter(entity.PrimaryKey())).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, datastore.ErrRecordNotExist
		}
		return false, datastore.NewDBError(err)
	}
	return true, nil
}

func makeNameFilter(name string) bson.D {
	return bson.D{
		{
			Key:   PrimaryKey,
			Value: name,
		},
	}
}

func makeEntityUpdate(entity interface{}) bson.M {
	return bson.M{"$set": entity}
}

func converToMap(model interface{}) (bson.M, error) {
	b, err := bson.Marshal(model)
	if err != nil {
		return nil, err
	}
	var re = make(bson.M)
	if err := bson.Unmarshal(b, &re); err != nil {
		return nil, err
	}
	return re, err
}
