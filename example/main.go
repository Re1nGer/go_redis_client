package main

import (
	"fmt"

	r "github.com/re1nger/go_redis_client"
)

func main() {
	c, err := r.NewClient("localhost", 6379)

	if err != nil {
		fmt.Printf("error while connecting to redis %s", err)
		return
	}

	/* 	arr1, _ := c.Set("test", "vallllll")
	   	arr, err := c.Get("test")
	   	arr2, err2 := c.Exists("123", "3123", "124123")
	   	arr3, err3 := c.SetWithOptions("test1", "val1", NewSetOpts().WithNX().WithEX(20)) */

	arr4, err4 := c.Get("test1")

	//fmt.Println("response", arr1, arr, arr2, err2, err, arr3, err3)

	fmt.Println("response", arr4, err4)
}
