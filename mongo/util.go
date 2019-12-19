package mongo

import (
	"encoding/binary"
	"hash/crc32"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// WriteSingleDocument is a helper for writing a document to a ResponseWriter.
func WriteSingleDocument(resp MessageResponseWriter, req *MessageRequest, doc bsoncore.Document) error {
	_, replyTo, _, opCode, msg, _ := wiremessage.ReadHeader(req.Message)

	var rmsg wiremessage.WireMessage

	switch opCode {
	case wiremessage.OpMsg:
		flags, _, _ := wiremessage.ReadMsgFlags(msg)
		if (flags & wiremessage.MoreToCome) == 0 {
			_, rmsg = wiremessage.AppendHeaderStart(rmsg, wiremessage.NextRequestID(), replyTo, wiremessage.OpMsg)
			rmsg = wiremessage.AppendMsgFlags(rmsg, flags&wiremessage.ChecksumPresent)
			rmsg = wiremessage.AppendMsgSectionType(rmsg, wiremessage.SingleDocument)
			rmsg = append(rmsg, doc...)
			if flags&wiremessage.ChecksumPresent != 0 {
				rmsg = append(rmsg, 0, 0, 0, 0)
			}
			rmsg = bsoncore.UpdateLength(rmsg, 0, int32(len(rmsg)))
			if flags&wiremessage.ChecksumPresent != 0 {
				checksum := crc32.Checksum(rmsg[0:len(rmsg)-4], crc32.MakeTable(crc32.Castagnoli))
				binary.LittleEndian.PutUint32(rmsg[len(rmsg)-4:], checksum)
			}
		}
	default:
		_, rmsg = wiremessage.AppendHeaderStart(rmsg, wiremessage.NextRequestID(), replyTo, wiremessage.OpReply)
		rmsg = wiremessage.AppendReplyFlags(rmsg, 0)
		rmsg = wiremessage.AppendReplyCursorID(rmsg, 0)
		rmsg = wiremessage.AppendReplyStartingFrom(rmsg, 0)
		rmsg = wiremessage.AppendReplyNumberReturned(rmsg, 1)
		rmsg = append(rmsg, doc...)
		rmsg = bsoncore.UpdateLength(rmsg, 0, int32(len(rmsg)))
	}

	return resp.WriteMessage(rmsg)
}
