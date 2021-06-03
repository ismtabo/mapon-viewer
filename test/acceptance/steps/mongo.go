package steps

import (
	"context"
	"encoding/json"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoSteps adds golium steps to interactuate with Mongo databases.
type MongoSteps struct{}

// InitializeSteps add steps additional godog steps.
func (s MongoSteps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Retrieve HTTP session
	var client *mongo.Client
	var database *mongo.Database
	var lastSingleResult *mongo.SingleResult
	// Initialize the steps
	scenCtx.Step(`^the mongodb URI "([^"]+)"$`, func(u string) error {
		uri := golium.ValueAsString(ctx, u)
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return err
		}
		return nil
	})
	scenCtx.Step(`^the mongodb database "([^"]+)"$`, func(d string) error {
		if client == nil {
			return errors.New("failed using database. empty client reference")
		}
		db := golium.ValueAsString(ctx, d)
		database = client.Database(db)
		return nil
	})
	scenCtx.Step(`^I create mongodb document in collection "([^"]+)" with properties$`, func(c string, t *godog.Table) error {
		if database == nil {
			return errors.New("failed using collection. empty database reference")
		}
		c = golium.ValueAsString(ctx, c)
		collection := database.Collection(c)
		data, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return err
		}
		if _, err := collection.InsertOne(ctx, data); err != nil {
			return err
		}
		return nil
	})
	scenCtx.Step(`^I create mongodb document in collection "([^"]+)" with properties$`, func(c string, t *godog.Table) error {
		if database == nil {
			return errors.New("failed using collection. empty database reference")
		}
		c = golium.ValueAsString(ctx, c)
		collection := database.Collection(c)
		data, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return err
		}
		if _, err := collection.InsertOne(ctx, data); err != nil {
			return err
		}
		return nil
	})
	scenCtx.Step(`^I create mongodb document in collection "([^"]+)" with properties$`, func(c string, t *godog.Table) error {
		if database == nil {
			return errors.New("failed using collection. empty database reference")
		}
		c = golium.ValueAsString(ctx, c)
		collection := database.Collection(c)
		data, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return err
		}
		if _, err := collection.InsertOne(ctx, data); err != nil {
			return err
		}
		return nil
	})
	scenCtx.Step(`^I search a mongodb single result document in collection "([^"]+)" with filter$`, func(c string, t *godog.Table) error {
		if database == nil {
			return errors.New("failed using collection. empty database reference")
		}
		c = golium.ValueAsString(ctx, c)
		collection := database.Collection(c)
		data, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return err
		}
		filter := bson.M{}
		for k, v := range data {
			filter[k] = v
		}
		if lastSingleResult = collection.FindOne(ctx, data); lastSingleResult.Err() != nil {
			err := lastSingleResult.Err()
			if err != mongo.ErrNoDocuments {
				return err
			}
		}
		return nil
	})
	scenCtx.Step(`^the mongodb search result should be empty$`, func() error {
		if lastSingleResult == nil {
			return errors.New("failed verifying empty mongo search result. empty last single result reference")
		}
		if err := lastSingleResult.Err(); err != nil {
			if err != mongo.ErrNoDocuments {
				return err
			}
			return nil
		}
		data := map[string]interface{}{}
		if err := lastSingleResult.Decode(data); err != nil {
			return err
		}
		return errors.Errorf("failed verifying empty mongo search result. actual result: %+v", data)
	})
	scenCtx.Step(`^the mongodb search result should be have the following JSON properties$`, func(t *godog.Table) error {
		if lastSingleResult == nil {
			return errors.New("failed verifying empty mongo search result. empty last single result reference")
		}
		if err := lastSingleResult.Err(); err != nil {
			if err != mongo.ErrNoDocuments {
				return err
			}
			return errors.Errorf("failed verifying empty mongo search result. no results found")
		}
		expectedProps, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return err
		}
		data := map[string]interface{}{}
		if err := lastSingleResult.Decode(data); err != nil {
			return err
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		actualProps := golium.NewMapFromJSONBytes(bytes)
		for key, expectedValue := range expectedProps {
			value := actualProps.Get(key)
			if value != expectedValue {
				return errors.Errorf("mismatch of json property '%s': expected '%s', actual '%s'", key, expectedValue, value)
			}
		}
		return nil
	})
	scenCtx.AfterScenario(func(sc *godog.Scenario, err error) {
		if client != nil {
			client.Disconnect(ctx)
		}
	})
	return ctx
}
