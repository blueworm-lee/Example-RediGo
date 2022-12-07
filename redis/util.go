package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func Ping() error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

func GetAllKeys() ([]string, error) {
	conn := Pool.Get()
	defer conn.Close()

	sResults, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		return nil, fmt.Errorf("KEYS *. error=%v", err)
	}

	return sResults, err
}

func ShowQueueValues(strKey string, nStart, nEnd int) ([]string, error) {
	conn := Pool.Get()
	defer conn.Close()

	//redis.String(conn.Do("LRANGE", strKey, nStart, nEnd))

	arrsResult, err := redis.Strings(conn.Do("LRANGE", strKey, nStart, nEnd))
	if err != nil {
		return nil, fmt.Errorf("LRANGE %v %v %v, error=%v", strKey, nStart, nEnd, err)
	}

	return arrsResult, nil

}

func GetQueueLen(strKey string) (int, error) {
	conn := Pool.Get()
	defer conn.Close()

	nSize, err := redis.Int(conn.Do("LLEN", strKey))
	if err != nil {
		return -1, fmt.Errorf("LLEN %v, error=%v", strKey, err)
	}

	return nSize, err

}

func PushQueue(sKey string, iValue interface{}) error {
	conn := Pool.Get()
	defer conn.Close()

	err := conn.Send("LPUSH", sKey, iValue)
	if err != nil {
		return fmt.Errorf("lpush Key=%s, Value=%s, err=%v", sKey, iValue, err)
	}
	return err
}

func PopQueue(sKey string) (string, error) {
	conn := Pool.Get()
	defer conn.Close()

	sResult, err := redis.String(conn.Do("RPOP", sKey))
	if err == redis.ErrNil {
		return sResult, nil
	} else if err != nil {
		return "", fmt.Errorf("RPOP Key=%s, err=%v", sKey, err)
	}

	return sResult, nil

}

// TimeOut must be shorter than DB connection timeout
func PopQueueBlock(sKey string, nTimeOut int) ([]string, error) {
	conn := Pool.Get()
	defer conn.Close()

	//	sResult, err := redis.String(redis.DoWithTimeout(conn, time.Second*time.Duration(nTimeOut), "BRPOP", sKey))
	arrsResult, err := redis.Strings(conn.Do("BRPOP", sKey, 1))
	if err == redis.ErrNil {
		return arrsResult, nil

	} else if err != nil {
		return nil, fmt.Errorf("BRPOP Key=%s, err=%v", sKey, err)
	}

	return arrsResult, nil

}
