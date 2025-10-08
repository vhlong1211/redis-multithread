package core

import (
	"TCPServer/constant"
	"errors"
	"fmt"
	"strconv"
	"syscall"
	"time"
)

func cmdPING(args []string) []byte {
	var res []byte
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}
	return res
}

func cmdSET(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SET' command"), false)
	}

	var key, value string
	var ttlMs int64 = -1

	key, value = args[0], args[1]
	if len(args) == 4 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		ttlMs = ttlSec * 1000
	}

	dictStore.Set(key, dictStore.NewObj(key, value, ttlMs))
	return constant.RespOk
}

func cmdGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}

	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}

	if dictStore.HasExpired(key) {
		return constant.RespNil
	}

	return Encode(obj.Value, false)
}

func cmdTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'TTL' command"), false)
	}
	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.TtlKeyNotExist
	}

	exp, isExpirySet := dictStore.GetExpiry(key)
	if !isExpirySet {
		return constant.TtlKeyExistNoExpire
	}

	remainMs := exp - uint64(time.Now().UnixMilli())
	if remainMs < 0 {
		return constant.TtlKeyNotExist
	}

	return Encode(int64(remainMs/1000), false)
}

func cmdEXPIRE(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'EXPIRE' command"), false)
	}

	key := args[0]
	ttlSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
	}

	// Use Get to check existence (Get removes expired keys)
	if dictStore.Get(key) == nil {
		return Encode(int64(0), false) // key does not exist
	}

	// Set expiry in milliseconds
	dictStore.SetExpiry(key, ttlSec*1000)
	return Encode(int64(1), false) // expiry set
}

func cmdDEL(args []string) []byte {
	if len(args) < 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'DEL' command"), false)
	}

	var deleted int64 = 0
	for _, k := range args {
		if dictStore.Del(k) {
			deleted++
		}
	}
	return Encode(deleted, false)
}

func cmdEXISTS(args []string) []byte {
	if len(args) < 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'EXISTS' command"), false)
	}

	var count int64 = 0
	for _, k := range args {
		// Use Get which returns nil for missing or expired keys
		if dictStore.Get(k) != nil {
			count++
		}
	}
	return Encode(count, false)
}

// ExecuteAndResponse given a Command, executes it and responses
func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	case "SET":
		res = cmdSET(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	case "EXPIRE":
		res = cmdEXPIRE(cmd.Args)
	case "DEL":
		res = cmdDEL(cmd.Args)
	case "EXISTS":
		res = cmdEXISTS(cmd.Args)
	case "ZADD":
		res = cmdZADD(cmd.Args)
	case "ZSCORE":
		res = cmdZSCORE(cmd.Args)
	case "ZRANK":
		res = cmdZRANK(cmd.Args)
	case "SADD":
		res = cmdSADD(cmd.Args)
	case "SREM":
		res = cmdSREM(cmd.Args)
	case "SMEMBERS":
		res = cmdSMEMBERS(cmd.Args)
	case "SISMEMBER":
		res = cmdSISMEMBER(cmd.Args)
	default:
		res = []byte(fmt.Sprintf("-CMD NOT FOUND\r\n"))
	}
	_, err := syscall.Write(connFd, res)
	return err
}
