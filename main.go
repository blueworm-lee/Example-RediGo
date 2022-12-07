package main

import (
	"fmt"
	"redis/redigo/redis"
	"strconv"
	"time"
)

func Push(strKeyName string, nNum int) {
	for i := 0; i < nNum; i++ {
		if err := redis.PushQueue(strKeyName, strconv.Itoa(i)); err != nil {
			fmt.Println("err: ", err.Error())
		} else {
			fmt.Println("Push OK")
		}
	}
}

func main() {
	var strKeyName = "jplee:test"

	//Push to Queue - LPUSH
	Push(strKeyName, 5)

	// - KEYS * (Warning, Redis can be blocked)
	if saResults, err := redis.GetAllKeys(); err != nil {
		fmt.Println("err: ", err.Error())
	} else {
		for _, sKey := range saResults {
			fmt.Println("Key=", sKey)
		}
	}

	//Queue Len(Size) - LLEN
	if nSize, err := redis.GetQueueLen(strKeyName); err != nil {
		fmt.Println("err=", err.Error())
	} else {
		fmt.Println("Key=", strKeyName, ", Size=", nSize)
	}

	//Queue Len(Size) - LRANGE
	if arrsResult, err := redis.ShowQueueValues(strKeyName, 0, -1); err != nil {
		fmt.Println("err=", err.Error())
	} else {
		for _, value := range arrsResult {
			fmt.Println("List. Value:", value)
		}
	}

	//Pop from Queue - RPOP
	for {
		sResult, err := redis.PopQueue(strKeyName)
		if err != nil {
			fmt.Println("err: ", err.Error())
		} else {
			if sResult == "" {
				break
			} else {
				fmt.Println("Poped. Key: ", strKeyName, ", Value: ", sResult)
			}

		}
		time.Sleep(time.Second * 1)
	}

	//Push to Queue - LPUSH
	go Push(strKeyName, 5)

	// - BRPOP
	for {
		sResult, err := redis.PopQueueBlock(strKeyName, 2)
		if err != nil {
			fmt.Println("err: ", err.Error())
		} else {
			fmt.Println("Poped. Key: ", strKeyName, ", Value: ", sResult)
		}
	}

}
