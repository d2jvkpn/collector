package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/d2jvkpn/collector/pkg/wrap"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestMongo(t *testing.T) {
	var (
		err    error
		bin    json.RawMessage
		data   map[string]any
		client *mongo.Client
	)

	if client, err = wrap.MongoClient(_TestConfig, "mongodb"); err != nil {
		t.Fatal(err)
	}

	coll := client.Database("test_collector").Collection("tests")

	bin = json.RawMessage(`{"hello": "world"}`)
	data = map[string]any{"name": "id0001", "data": bin}
	result, err := coll.InsertOne(_TestCtx, data)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("~~~ result:", result)
	// ~~~ result: &{ObjectID("646c957b75b54fd59e4339de")}

	bts, _ := json.Marshal(data)
	fmt.Printf("~~~ bts: %s, err: %v\n", bts, err)
	// ~~~ bts: {"data":{"hello":"world"},"name":"id0001"}, err: <nil>

	bin = json.RawMessage(`{"hello": "world`)
	data = map[string]any{"name": "id0001", "data": bin}
	bts, err = json.Marshal(data)
	fmt.Printf("~~~ bts: %s, err: %v\n", bts, err)
	// ~~~ bts: , err: json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input
}

/*
db.tests.findOne()
{
	"_id" : ObjectId("646c957b75b54fd59e4339de"),
	"name" : "id0001",
	"data" : BinData(0,"eyJoZWxsbyI6ICJ3b3JsZCJ9")
}
*/
