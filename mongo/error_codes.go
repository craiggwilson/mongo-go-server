package mongo

// Taken from the error codes at
// https://github.com/mongodb/mongo/blob/master/src/mongo/base/error_codes.err
const (
	CodeInternalError             = 1
	CodeBadValue                  = 2
	CodeNoSuchKey                 = 4
	CodeGraphContainsCycle        = 5
	CodeFailedToParse             = 9
	CodeUserNotFound              = 11
	CodeUnauthorized              = 13
	CodeTypeMismatch              = 14
	CodeProtocolError             = 17
	CodeAuthenticationFailed      = 18
	CodeIllegalOperation          = 20
	CodeInvalidBSON               = 22
	CodeNamespaceNotFound         = 26
	CodeNonExistentPath           = 29
	CodeRoleNotFound              = 31
	CodeFileStreamFailed          = 39
	CodeCursorNotFound            = 43
	CodeUserDataInconsistent      = 45
	CodeCommandNotFound           = 59
	CodeInvalidOptions            = 72
	CodeInvalidNamespace          = 73
	CodeCommandNotSupported       = 115
	CodeCommandFailed             = 125
	CodeCommandNotSupportedOnView = 166
	CodeInternalErrorNotSupported = 235
	CodeDuplicateKey              = 11000
	CodeInterrupted               = 11601
	CodeOutOfDiskSpace            = 14031
)

// Taken from the error code names at
// https://github.com/mongodb/mongo/blob/master/src/mongo/base/error_codes.err
func CodeToName(code int32) string {
	switch code {
	case 0:
		return "OK"
	case 1:
		return "InternalError"
	case 2:
		return "BadValue"
	case 3:
		return "OBSOLETE_DuplicateKey"
	case 4:
		return "NoSuchKey"
	case 5:
		return "GraphContainsCycle"
	case 6:
		return "HostUnreachable"
	case 7:
		return "HostNotFound"
	case 8:
		return "UnknownError"
	case 9:
		return "FailedToParse"
	case 10:
		return "CannotMutateObject"
	case 11:
		return "UserNotFound"
	case 12:
		return "UnsupportedFormat"
	case 13:
		return "Unauthorized"
	case 14:
		return "TypeMismatch"
	case 15:
		return "Overflow"
	case 16:
		return "InvalidLength"
	case 17:
		return "ProtocolError"
	case 18:
		return "AuthenticationFailed"
	case 19:
		return "CannotReuseObject"
	case 20:
		return "IllegalOperation"
	case 21:
		return "EmptyArrayOperation"
	case 22:
		return "InvalidBSON"
	case 23:
		return "AlreadyInitialized"
	case 24:
		return "LockTimeout"
	case 25:
		return "RemoteValidationError"
	case 26:
		return "NamespaceNotFound"
	case 27:
		return "IndexNotFound"
	case 28:
		return "PathNotViable"
	case 29:
		return "NonExistentPath"
	case 30:
		return "InvalidPath"
	case 31:
		return "RoleNotFound"
	case 32:
		return "RolesNotRelated"
	case 33:
		return "PrivilegeNotFound"
	case 34:
		return "CannotBackfillArray"
	case 35:
		return "UserModificationFailed"
	case 36:
		return "RemoteChangeDetected"
	case 37:
		return "FileRenameFailed"
	case 38:
		return "FileNotOpen"
	case 39:
		return "FileStreamFailed"
	case 40:
		return "ConflictingUpdateOperators"
	case 41:
		return "FileAlreadyOpen"
	case 42:
		return "LogWriteFailed"
	case 43:
		return "CursorNotFound"
	case 45:
		return "UserDataInconsistent"
	case 46:
		return "LockBusy"
	case 47:
		return "NoMatchingDocument"
	case 48:
		return "NamespaceExists"
	case 49:
		return "InvalidRoleModification"
	case 50:
		return "MaxTimeMSExpired"
	case 51:
		return "ManualInterventionRequired"
	case 52:
		return "DollarPrefixedFieldName"
	case 53:
		return "InvalidIdField"
	case 54:
		return "NotSingleValueField"
	case 55:
		return "InvalidDBRef"
	case 56:
		return "EmptyFieldName"
	case 57:
		return "DottedFieldName"
	case 58:
		return "RoleModificationFailed"
	case 59:
		return "CommandNotFound"
	case 60:
		return "OBSOLETE_DatabaseNotFound"
	case 61:
		return "ShardKeyNotFound"
	case 62:
		return "OplogOperationUnsupported"
	case 63:
		return "StaleShardVersion"
	case 64:
		return "WriteConcernFailed"
	case 65:
		return "MultipleErrorsOccurred"
	case 66:
		return "ImmutableField"
	case 67:
		return "CannotCreateIndex"
	case 68:
		return "IndexAlreadyExists"
	case 69:
		return "AuthSchemaIncompatible"
	case 70:
		return "ShardNotFound"
	case 71:
		return "ReplicaSetNotFound"
	case 72:
		return "InvalidOptions"
	case 73:
		return "InvalidNamespace"
	case 74:
		return "NodeNotFound"
	case 75:
		return "WriteConcernLegacyOK"
	case 76:
		return "NoReplicationEnabled"
	case 77:
		return "OperationIncomplete"
	case 78:
		return "CommandResultSchemaViolation"
	case 79:
		return "UnknownReplWriteConcern"
	case 80:
		return "RoleDataInconsistent"
	case 81:
		return "NoMatchParseContext"
	case 82:
		return "NoProgressMade"
	case 83:
		return "RemoteResultsUnavailable"
	case 84:
		return "DuplicateKeyValue"
	case 85:
		return "IndexOptionsConflict"
	case 86:
		return "IndexKeySpecsConflict"
	case 87:
		return "CannotSplit"
	case 88:
		return "SplitFailed_OBSOLETE"
	case 89:
		return "NetworkTimeout"
	case 90:
		return "CallbackCanceled"
	case 91:
		return "ShutdownInProgress"
	case 92:
		return "SecondaryAheadOfPrimary"
	case 93:
		return "InvalidReplicaSetConfig"
	case 94:
		return "NotYetInitialized"
	case 95:
		return "NotSecondary"
	case 96:
		return "OperationFailed"
	case 97:
		return "NoProjectionFound"
	case 98:
		return "DBPathInUse"
	case 100:
		return "UnsatisfiableWriteConcern"
	case 101:
		return "OutdatedClient"
	case 102:
		return "IncompatibleAuditMetadata"
	case 103:
		return "NewReplicaSetConfigurationIncompatible"
	case 104:
		return "NodeNotElectable"
	case 105:
		return "IncompatibleShardingMetadata"
	case 106:
		return "DistributedClockSkewed"
	case 107:
		return "LockFailed"
	case 108:
		return "InconsistentReplicaSetNames"
	case 109:
		return "ConfigurationInProgress"
	case 110:
		return "CannotInitializeNodeWithData"
	case 111:
		return "NotExactValueField"
	case 112:
		return "WriteConflict"
	case 113:
		return "InitialSyncFailure"
	case 114:
		return "InitialSyncOplogSourceMissing"
	case 115:
		return "CommandNotSupported"
	case 116:
		return "DocTooLargeForCapped"
	case 117:
		return "ConflictingOperationInProgress"
	case 118:
		return "NamespaceNotSharded"
	case 119:
		return "InvalidSyncSource"
	case 120:
		return "OplogStartMissing"
	case 121: // Only for the document validator on collections
		return "DocumentValidationFailure"
	case 122:
		return "OBSOLETE_ReadAfterOptimeTimeout"
	case 123:
		return "NotAReplicaSet"
	case 124:
		return "IncompatibleElectionProtocol"
	case 125:
		return "CommandFailed"
	case 126:
		return "RPCProtocolNegotiationFailed"
	case 127:
		return "UnrecoverableRollbackError"
	case 128:
		return "LockNotFound"
	case 129:
		return "LockStateChangeFailed"
	case 130:
		return "SymbolNotFound"
	case 132:
		return "OBSOLETE_ConfigServersInconsistent"
	case 133:
		return "FailedToSatisfyReadPreference"
	case 134:
		return "ReadConcernMajorityNotAvailableYet"
	case 135:
		return "StaleTerm"
	case 136:
		return "CappedPositionLost"
	case 137:
		return "IncompatibleShardingConfigVersion"
	case 138:
		return "RemoteOplogStale"
	case 139:
		return "JSInterpreterFailure"
	case 140:
		return "InvalidSSLConfiguration"
	case 141:
		return "SSLHandshakeFailed"
	case 142:
		return "JSUncatchableError"
	case 143:
		return "CursorInUse"
	case 144:
		return "IncompatibleCatalogManager"
	case 145:
		return "PooledConnectionsDropped"
	case 146:
		return "ExceededMemoryLimit"
	case 147:
		return "ZLibError"
	case 148:
		return "ReadConcernMajorityNotEnabled"
	case 149:
		return "NoConfigMaster"
	case 150:
		return "StaleEpoch"
	case 151:
		return "OperationCannotBeBatched"
	case 152:
		return "OplogOutOfOrder"
	case 153:
		return "ChunkTooBig"
	case 154:
		return "InconsistentShardIdentity"
	case 155:
		return "CannotApplyOplogWhilePrimary"
	case 156:
		return "OBSOLETE_NeedsDocumentMove"
	case 157:
		return "CanRepairToDowngrade"
	case 158:
		return "MustUpgrade"
	case 159:
		return "DurationOverflow"
	case 160:
		return "MaxStalenessOutOfRange"
	case 161:
		return "IncompatibleCollationVersion"
	case 162:
		return "CollectionIsEmpty"
	case 163:
		return "ZoneStillInUse"
	case 164:
		return "InitialSyncActive"
	case 165:
		return "ViewDepthLimitExceeded"
	case 166:
		return "CommandNotSupportedOnView"
	case 167:
		return "OptionNotSupportedOnView"
	case 168:
		return "InvalidPipelineOperator"
	case 169:
		return "CommandOnShardedViewNotSupportedOnMongod"
	case 170:
		return "TooManyMatchingDocuments"
	case 171:
		return "CannotIndexParallelArrays"
	case 172:
		return "TransportSessionClosed"
	case 173:
		return "TransportSessionNotFound"
	case 174:
		return "TransportSessionUnknown"
	case 175:
		return "QueryPlanKilled"
	case 176:
		return "FileOpenFailed"
	case 177:
		return "ZoneNotFound"
	case 178:
		return "RangeOverlapConflict"
	case 179:
		return "WindowsPdhError"
	case 180:
		return "BadPerfCounterPath"
	case 181:
		return "AmbiguousIndexKeyPattern"
	case 182:
		return "InvalidViewDefinition"
	case 183:
		return "ClientMetadataMissingField"
	case 184:
		return "ClientMetadataAppNameTooLarge"
	case 185:
		return "ClientMetadataDocumentTooLarge"
	case 186:
		return "ClientMetadataCannotBeMutated"
	case 187:
		return "LinearizableReadConcernError"
	case 188:
		return "IncompatibleServerVersion"
	case 189:
		return "PrimarySteppedDown"
	case 190:
		return "MasterSlaveConnectionFailure"
	case 191:
		return "OBSOLETE_BalancerLostDistributedLock"
	case 192:
		return "FailPointEnabled"
	case 193:
		return "NoShardingEnabled"
	case 194:
		return "BalancerInterrupted"
	case 195:
		return "ViewPipelineMaxSizeExceeded"
	case 197:
		return "InvalidIndexSpecificationOption"
	case 198:
		return "OBSOLETE_ReceivedOpReplyMessage"
	case 199:
		return "ReplicaSetMonitorRemoved"
	case 200:
		return "ChunkRangeCleanupPending"
	case 201:
		return "CannotBuildIndexKeys"
	case 202:
		return "NetworkInterfaceExceededTimeLimit"
	case 203:
		return "ShardingStateNotInitialized"
	case 204:
		return "TimeProofMismatch"
	case 205:
		return "ClusterTimeFailsRateLimiter"
	case 206:
		return "NoSuchSession"
	case 207:
		return "InvalidUUID"
	case 208:
		return "TooManyLocks"
	case 209:
		return "StaleClusterTime"
	case 210:
		return "CannotVerifyAndSignLogicalTime"
	case 211:
		return "KeyNotFound"
	case 212:
		return "IncompatibleRollbackAlgorithm"
	case 213:
		return "DuplicateSession"
	case 214:
		return "AuthenticationRestrictionUnmet"
	case 215:
		return "DatabaseDropPending"
	case 216:
		return "ElectionInProgress"
	case 217:
		return "IncompleteTransactionHistory"
	case 218:
		return "UpdateOperationFailed"
	case 219:
		return "FTDCPathNotSet"
	case 220:
		return "FTDCPathAlreadySet"
	case 221:
		return "IndexModified"
	case 222:
		return "CloseChangeStream"
	case 223:
		return "IllegalOpMsgFlag"
	case 224:
		return "QueryFeatureNotAllowed"
	case 225:
		return "TransactionTooOld"
	case 226:
		return "AtomicityFailure"
	case 227:
		return "CannotImplicitlyCreateCollection"
	case 228:
		return "SessionTransferIncomplete"
	case 229:
		return "MustDowngrade"
	case 230:
		return "DNSHostNotFound"
	case 231:
		return "DNSProtocolError"
	case 232:
		return "MaxSubPipelineDepthExceeded"
	case 233:
		return "TooManyDocumentSequences"
	case 234:
		return "RetryChangeStream"
	case 235:
		return "InternalErrorNotSupported"
	case 236:
		return "ForTestingErrorExtraInfo"
	case 237:
		return "CursorKilled"
	case 238:
		return "NotImplemented"
	case 239:
		return "SnapshotTooOld"
	case 240:
		return "DNSRecordTypeMismatch"
	case 241:
		return "ConversionFailure"
	case 242:
		return "CannotCreateCollection"
	case 243:
		return "IncompatibleWithUpgradedServer"
	case 244:
		return "NOT_YET_AVAILABLE_TransactionAborted"
	case 245:
		return "BrokenPromise"
	case 246:
		return "SnapshotUnavailable"
	case 247:
		return "ProducerConsumerQueueBatchTooLarge"
	case 248:
		return "ProducerConsumerQueueEndClosed"
	case 249:
		return "StaleDbVersion"
	case 250:
		return "StaleChunkHistory"
	case 251:
		return "NoSuchTransaction"
	case 252:
		return "ReentrancyNotAllowed"
	case 253:
		return "FreeMonHttpInFlight"
	case 254:
		return "FreeMonHttpTemporaryFailure"
	case 255:
		return "FreeMonHttpPermanentFailure"
	case 256:
		return "TransactionCommitted"
	case 257:
		return "TransactionTooLarge"
	case 258:
		return "UnknownFeatureCompatibilityVersion"
	case 259:
		return "KeyedExecutorRetry"
	case 260:
		return "InvalidResumeToken"
	case 261:
		return "TooManyLogicalSessions"
	case 262:
		return "ExceededTimeLimit"
	case 263:
		return "OperationNotSupportedInTransaction"
	case 264:
		return "TooManyFilesOpen"
	case 265:
		return "OrphanedRangeCleanUpFailed"
	case 266:
		return "FailPointSetFailed"
	case 267:
		return "PreparedTransactionInProgress"
	case 268:
		return "CannotBackup"
	case 269:
		return "DataModifiedByRepair"
	case 270:
		return "RepairedReplicaSetNode"
	case 271:
		return "JSInterpreterFailureWithStack"
	case 272:
		return "MigrationConflict"
	case 273:
		return "ProducerConsumerQueueProducerQueueDepthExceeded"
	case 274:
		return "ProducerConsumerQueueConsumed"
	case 275:
		return "ExchangePassthrough"
	case 276:
		return "IndexBuildAborted"
	case 277:
		return "AlarmAlreadyFulfilled"
	case 278:
		return "UnsatisfiableCommitQuorum"
	case 279:
		return "ClientDisconnect"
	case 280:
		return "ChangeStreamFatalError"
	case 283:
		return "WouldChangeOwningShard"
	case 284:
		return "ForTestingErrorExtraInfoWithExtraInfoInNamespace"
	case 285:
		return "IndexBuildAlreadyInProgress"
	case 286:
		return "ChangeStreamHistoryLost"

	// Error codes 4000-8999 are reserved.

	// Non-sequential error codes (for compatibility only
	case 9001:
		return "SocketException"
	case 9996:
		return "OBSOLETE_RecvStaleConfig"
	case 10107:
		return "NotMaster"
	case 10003:
		return "CannotGrowDocumentInCappedNamespace"
	case 10334:
		return "BSONObjectTooLarge"
	case 11000:
		return "DuplicateKey"
	case 11600:
		return "InterruptedAtShutdown"
	case 11601:
		return "Interrupted"
	case 11602:
		return "InterruptedDueToStepDown"
	case 12586:
		return "BackgroundOperationInProgressForDatabase"
	case 12587:
		return "BackgroundOperationInProgressForNamespace"
	case 13436:
		return "NotMasterOrSecondary"
	case 13435:
		return "NotMasterNoSlaveOk"
	case 13334:
		return "ShardKeyTooBig"
	case 13388:
		return "StaleConfig"
	case 13297:
		return "DatabaseDifferCase"
	case 13104:
		return "OBSOLETE_PrepareConfigsFailed"
	case 14031:
		return "OutOfDiskSpace"
	case 17280:
		return "KeyTooLong"
	default:
		return ""
	}
}
