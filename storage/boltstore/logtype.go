package boltstore

import (
	"github.com/hashicorp/raft"
)

func LogTypeToPB(logType raft.LogType) LogType {
	switch logType {
	case raft.LogCommand:
		return LogType_LogCommand
	case raft.LogNoop:
		return LogType_LogNoop
	case raft.LogAddPeerDeprecated:
		return LogType_LogAddPeerDeprecated
	case raft.LogRemovePeerDeprecated:
		return LogType_LogRemovePeerDeprecated
	case raft.LogBarrier:
		return LogType_LogBarrier
	case raft.LogConfiguration:
		return LogType_LogConfiguration
	default:
		return LogType_LogCommand
	}
}

func LogTypeFromPB(pbLogType LogType) raft.LogType {
	switch pbLogType {
	case LogType_LogCommand:
		return raft.LogCommand
	case LogType_LogNoop:
		return raft.LogNoop
	case LogType_LogAddPeerDeprecated:
		return raft.LogAddPeerDeprecated
	case LogType_LogRemovePeerDeprecated:
		return raft.LogRemovePeerDeprecated
	case LogType_LogBarrier:
		return raft.LogBarrier
	case LogType_LogConfiguration:
		return raft.LogConfiguration
	default:
		return raft.LogCommand
	}
}
