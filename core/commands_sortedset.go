package core

import (
	"TCPServer/constant"
	"TCPServer/data_structure"
	"errors"
	"fmt"
	"strconv"
)

func cmdZADD(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZADD' command"), false)
	}
	key := args[0]
	scoreIndex := 1

	numScoreEleArgs := len(args) - scoreIndex
	if numScoreEleArgs%2 == 1 || numScoreEleArgs == 0 {
		return Encode(errors.New(fmt.Sprintf("(error) Wrong number of (score, member) arg: %d", numScoreEleArgs)), false)
	}

	zset, exist := zsetStore[key]
	if !exist {
		zset = data_structure.CreateZSet()
		zsetStore[key] = zset
	}

	count := 0
	for i := scoreIndex; i < len(args); i += 2 {
		member := args[i+1]
		score, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return Encode(errors.New("(error) Score must be floating point number"), false)
		}
		ret := zset.Add(score, member)
		if ret != 1 {
			return Encode(errors.New("error when adding element"), false)
		}
		count++
	}
	return Encode(count, false)
}

func cmdZSCORE(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZSCORE' command"), false)
	}
	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}
	ret, score := zset.GetScore(member)
	if ret == -1 {
		return constant.RespNil
	}
	return Encode(fmt.Sprintf("%f", score), false)
}

func cmdZRANK(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZRANK' command"), false)
	}
	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}
	rank, _ := zset.GetRank(member, false)
	return Encode(rank, false)
}
