package mongotest

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

func NewOpMsg(cmd bsoncore.Document) []byte {
	_, msg := wiremessage.AppendHeaderStart(nil, 1, 1, wiremessage.OpMsg)
	msg = wiremessage.AppendMsgFlags(msg, 0)
	msg = wiremessage.AppendMsgSectionType(msg, wiremessage.SingleDocument)
	msg = append(msg, cmd...)
	msg = bsoncore.UpdateLength(msg, 0, int32(len(msg)))
	return msg
}

func NewDocument(elems ...interface{}) []byte {
	if len(elems)%2 != 0 {
		panic("must have an even number of elems")
	}

	d := primitive.D{}
	for i := 0; i < len(elems); i += 2 {
		d = append(d, primitive.E{Key: elems[i].(string), Value: elems[i+1]})
	}

	doc, err := bson.Marshal(d)
	if err != nil {
		panic(err)
	}

	return doc
}
