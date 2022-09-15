package unit_service

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"testing"
	"ukgm/internal/app/config/env"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"

	"firebase.google.com/go/auth"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func MongoDeclareForTesting() (repo *mongorepository.MongoRepo, err error) {
	dbUsername := env.Get("DB_USERNAME")
	dbPassword := env.Get("DB_PASSWORD")
	dbHost := env.Get("DB_HOST")
	dbPort := env.Get("DB_PORT")
	dbName := env.Get("DB_DATABASE")
	dbURL := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", dbUsername, url.QueryEscape(dbPassword), dbHost, dbPort)
	repo, err = mongorepository.NewMongoRepo(dbURL, dbName)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

func TestStoreUnit(t *testing.T) {
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(repo)

	// var uid string = "n8xgf00q7udFuCX7jbTIhtyf2gk2"
	// firebaseApi := firebase_api.NewFirebaseApi()
	// firebaseToken := firebaseApi.GenerateToken(uid)
	// tokenID, err := firebaseApi.SignInWithCustomToken(firebaseToken)
	// if err != nil {
	// 	t.Log(err)
	// }

	// authToken, err := firebaseApi.EncodeToken(tokenID)
	// if err != nil {
	// 	t.Log(err)
	// }

	srv := New(repo)

	t.Run("success store unit", func(t *testing.T) {
		unit := entities.Unit{
			DistrictID: "631dd39ef9c8ba488679be28",
			Name:       "unit 2 test",
		}

		res, err := srv.StoreData(unit, &auth.Token{})
		t.Log(err)
		log.Println(err)
		assert.Equal(t, err, nil)
		t.Log(res)
	})
}

func TestTransactionCommit(t *testing.T) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	var client = repo.Client
	var collection *mongo.Collection
	var ctx = context.Background()
	var id = primitive.NewObjectID()
	var doc = bson.M{"_id": id, "hometown": "Atlanta 2", "year": int32(1998)}
	var result *mongo.UpdateResult
	var session mongo.Session
	var update = bson.D{{Key: "$set", Value: bson.D{{Key: "year", Value: int32(2000)}}}}

	// if client, err = repo.Client; err != nil {
	// 	t.Fatal(err)
	// }

	defer client.Disconnect(ctx)
	collection = client.Database("ukgm_dev").Collection("tests")
	if _, err = collection.InsertOne(ctx, doc); err != nil {
		t.Fatal(err)
	}

	if session, err = client.StartSession(); err != nil {
		t.Fatal(err)
	}
	if err = session.StartTransaction(); err != nil {
		t.Fatal(err)
	}
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		id = primitive.NewObjectID()
		doc = bson.M{"_id": id, "hometown": "Atlanta 2", "year": int32(1998)}
		if _, err = collection.InsertOne(ctx, doc); err != nil {
			t.Fatal(err)
		}

		result, err = collection.UpdateOne(sc, bson.M{"_id": id}, update)

		var err = fmt.Errorf("update failed")

		if err != nil {
			log.Fatal(err)
		}

		if result.MatchedCount != 1 || result.ModifiedCount != 1 {
			// t.Fatal("replace failed, expected 1 but got", result.MatchedCount)
			log.Fatal("replace failed, expected 1 but got", result.MatchedCount)
		}

		err = session.CommitTransaction(sc)

		if err != nil {
			log.Fatal(err)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	session.EndSession(ctx)

	// var v bson.M
	// if err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&v); err != nil {
	// 	t.Fatal(err)
	// }
	// if v["year"] != int32(2000) {
	// 	t.Log(v)
	// 	t.Fatal("expected 2000 but got", v["year"])
	// }

	// res, _ := collection.DeleteOne(ctx, bson.M{"_id": id})
	// if res.DeletedCount != 1 {
	// 	t.Fatal("delete failed, expected 1 but got", res.DeletedCount)
	// }
}
