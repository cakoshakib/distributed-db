package boltstore

import (
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func LogToPB(log *raft.Log) ([]byte, error) {
	pbLog := &Log{
		Index:      log.Index,
		Term:       log.Term,
		Type:       LogTypeToPB(log.Type),
		Data:       log.Data,
		Extensions: log.Extensions,
		AppendedAt: timestamppb.New(log.AppendedAt),
	}
	return proto.Marshal(pbLog)
}

func LogFromPB(data []byte) (*raft.Log, error) {
	pbLog := Log{}
	if err := proto.Unmarshal(data, &pbLog); err != nil {
		return nil, err
	}

	raftLog := &raft.Log{
		Index:      pbLog.Index,
		Term:       pbLog.Term,
		Type:       LogTypeFromPB(pbLog.Type),
		Data:       pbLog.Data,
		Extensions: pbLog.Extensions,
		AppendedAt: pbLog.AppendedAt.AsTime(),
	}
	return raftLog, nil
}
