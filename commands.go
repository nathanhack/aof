package aof

import "strings"

var validCommands = map[string]bool{
	"ACL":               true,
	"APPEND":            true,
	"ASKING":            true,
	"AUTH":              true,
	"BGREWRITEAOF":      true,
	"BGSAVE":            true,
	"BITCOUNT":          true,
	"BITFIELD":          true,
	"BITFIELD_RO":       true,
	"BITOP":             true,
	"BITPOS":            true,
	"BLPOP":             true,
	"BRPOP":             true,
	"BRPOPLPUSH":        true,
	"BLMOVE":            true,
	"LMPOP":             true,
	"BLMPOP":            true,
	"BZPOPMIN":          true,
	"BZPOPMAX":          true,
	"BZMPOP":            true,
	"CLIENT":            true,
	"CLUSTER":           true,
	"COMMAND":           true,
	"CONFIG":            true,
	"COPY":              true,
	"DBSIZE":            true,
	"DEBUG":             true,
	"DECR":              true,
	"DECRBY":            true,
	"DEL":               true,
	"DISCARD":           true,
	"DUMP":              true,
	"ECHO":              true,
	"EVAL":              true,
	"EVAL_RO":           true,
	"EVALSHA":           true,
	"EVALSHA_RO":        true,
	"EXEC":              true,
	"EXISTS":            true,
	"EXPIRE":            true,
	"EXPIREAT":          true,
	"EXPIRETIME":        true,
	"FAILOVER":          true,
	"FLUSHALL":          true,
	"FLUSHDB":           true,
	"GEOADD":            true,
	"GEOHASH":           true,
	"GEOPOS":            true,
	"GEODIST":           true,
	"GEORADIUS":         true,
	"GEORADIUSBYMEMBER": true,
	"GEOSEARCH":         true,
	"GEOSEARCHSTORE":    true,
	"GET":               true,
	"GETBIT":            true,
	"GETDEL":            true,
	"GETEX":             true,
	"GETRANGE":          true,
	"GETSET":            true,
	"HDEL":              true,
	"HELLO":             true,
	"HEXISTS":           true,
	"HGET":              true,
	"HGETALL":           true,
	"HINCRBY":           true,
	"HINCRBYFLOAT":      true,
	"HKEYS":             true,
	"HLEN":              true,
	"HMGET":             true,
	"HMSET":             true,
	"HSET":              true,
	"HSETNX":            true,
	"HRANDFIELD":        true,
	"HSTRLEN":           true,
	"HVALS":             true,
	"INCR":              true,
	"INCRBY":            true,
	"INCRBYFLOAT":       true,
	"INFO":              true,
	"LOLWUT":            true,
	"KEYS":              true,
	"LASTSAVE":          true,
	"LINDEX":            true,
	"LINSERT":           true,
	"LLEN":              true,
	"LPOP":              true,
	"LPOS":              true,
	"LPUSH":             true,
	"LPUSHX":            true,
	"LRANGE":            true,
	"LREM":              true,
	"LSET":              true,
	"LTRIM":             true,
	"MEMORY":            true,
	"MGET":              true,
	"MIGRATE":           true,
	"MODULE":            true,
	"MONITOR":           true,
	"MOVE":              true,
	"MSET":              true,
	"MSETNX":            true,
	"MULTI":             true,
	"OBJECT":            true,
	"PERSIST":           true,
	"PEXPIRE":           true,
	"PEXPIREAT":         true,
	"PEXPIRETIME":       true,
	"PFADD":             true,
	"PFCOUNT":           true,
	"PFMERGE":           true,
	"PING":              true,
	"PSETEX":            true,
	"PSUBSCRIBE":        true,
	"PUBSUB":            true,
	"PTTL":              true,
	"PUBLISH":           true,
	"PUNSUBSCRIBE":      true,
	"QUIT":              true,
	"RANDOMKEY":         true,
	"READONLY":          true,
	"READWRITE":         true,
	"RENAME":            true,
	"RENAMENX":          true,
	"RESET":             true,
	"RESTORE":           true,
	"ROLE":              true,
	"RPOP":              true,
	"RPOPLPUSH":         true,
	"LMOVE":             true,
	"RPUSH":             true,
	"RPUSHX":            true,
	"SADD":              true,
	"SAVE":              true,
	"SCARD":             true,
	"SCRIPT":            true,
	"SDIFF":             true,
	"SDIFFSTORE":        true,
	"SELECT":            true,
	"SET":               true,
	"SETBIT":            true,
	"SETEX":             true,
	"SETNX":             true,
	"SETRANGE":          true,
	"SHUTDOWN":          true,
	"SINTER":            true,
	"SINTERCARD":        true,
	"SINTERSTORE":       true,
	"SISMEMBER":         true,
	"SMISMEMBER":        true,
	"SLAVEOF":           true,
	"REPLICAOF":         true,
	"SLOWLOG":           true,
	"SMEMBERS":          true,
	"SMOVE":             true,
	"SORT":              true,
	"SORT_RO":           true,
	"SPOP":              true,
	"SRANDMEMBER":       true,
	"SREM":              true,
	"STRALGO":           true,
	"STRLEN":            true,
	"SUBSCRIBE":         true,
	"SUNION":            true,
	"SUNIONSTORE":       true,
	"SWAPDB":            true,
	"SYNC":              true,
	"PSYNC":             true,
	"TIME":              true,
	"TOUCH":             true,
	"TTL":               true,
	"TYPE":              true,
	"UNSUBSCRIBE":       true,
	"UNLINK":            true,
	"UNWATCH":           true,
	"WAIT":              true,
	"WATCH":             true,
	"ZADD":              true,
	"ZCARD":             true,
	"ZCOUNT":            true,
	"ZDIFF":             true,
	"ZDIFFSTORE":        true,
	"ZINCRBY":           true,
	"ZINTER":            true,
	"ZINTERCARD":        true,
	"ZINTERSTORE":       true,
	"ZLEXCOUNT":         true,
	"ZPOPMAX":           true,
	"ZPOPMIN":           true,
	"ZMPOP":             true,
	"ZRANDMEMBER":       true,
	"ZRANGESTORE":       true,
	"ZRANGE":            true,
	"ZRANGEBYLEX":       true,
	"ZREVRANGEBYLEX":    true,
	"ZRANGEBYSCORE":     true,
	"ZRANK":             true,
	"ZREM":              true,
	"ZREMRANGEBYLEX":    true,
	"ZREMRANGEBYRANK":   true,
	"ZREMRANGEBYSCORE":  true,
	"ZREVRANGE":         true,
	"ZREVRANGEBYSCORE":  true,
	"ZREVRANK":          true,
	"ZSCORE":            true,
	"ZUNION":            true,
	"ZMSCORE":           true,
	"ZUNIONSTORE":       true,
	"SCAN":              true,
	"SSCAN":             true,
	"HSCAN":             true,
	"ZSCAN":             true,
	"XINFO":             true,
	"XADD":              true,
	"XTRIM":             true,
	"XDEL":              true,
	"XRANGE":            true,
	"XREVRANGE":         true,
	"XLEN":              true,
	"XREAD":             true,
	"XGROUP":            true,
	"XREADGROUP":        true,
	"XACK":              true,
	"XCLAIM":            true,
	"XAUTOCLAIM":        true,
	"XPENDING":          true,
	"LATENCY":           true,
}

func SetValidCommand(names ...string) {
	for _, name := range names {
		validCommands[strings.ToUpper(name)] = true
	}
}

func UnsetValidCommand(names ...string) {
	for _, name := range names {
		validCommands[strings.ToUpper(name)] = false
	}
}

func ClearAllValidCommands() {
	for k := range validCommands {
		validCommands[k] = false
	}
}