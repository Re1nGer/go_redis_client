package main

import (
	"fmt"
	"log"
	"time"

	r "github.com/re1nger/go_redis_client"
)

func main() {
	c, err := r.NewClient("localhost", 15358,
		r.WithSocketConnectTimeout(time.Duration(20*time.Second)),
		r.WithSocketTimeout(time.Duration(20*time.Second)))

	//c, err := r.NewClient("localhost", 6379)

	if err != nil {
		fmt.Printf("error while connecting to redis %s", err)
		return
	}

	//arr, err := c.HExpire("no-key", 20, []string{"field1", "field2"}, r.WithNX())

	//arr1, err1 := c.Set("mynewkey", "field1") //guess gotta offer map as input -> ticket ?

	//arr4, err4 := c.SetWithOptions("keyexpire", "val", r.NewSetOpts().WithEX(200))

	//arr2, err2 := c.Get("mynewkey")

	//arr3, err3 := c.HGetAll("mykey")

	/* 	arr1, _ := c.Set("test", "vallllll")
	   	arr, err := c.Get("test")
	   	arr2, err2 := c.Exists("123", "3123", "124123")
	   	arr3, err3 := c.SetWithOptions("test1", "val1", NewSetOpts().WithNX().WithEX(20)) */

	//fmt.Println("response", arr1, arr, arr2, err2, err, arr3, err3)

	//arr1, err1 := c.RPush("mylist", "a", "b", "c", "d", "1", "2", "3", "4", "3", "3", "3")

	//arr2, err2 := c.LPos("mylist", "3", r.WithCountLPos(3), r.WithRankLPos(2))

	//fmt.Println("response", arr1, err1, arr2, err2)

	// Test 1: Basic MULTI/EXEC
	fmt.Println("Test 1: Basic MULTI/EXEC")
	if err := testBasicTransaction(c); err != nil {
		log.Printf("Test 1 failed: %v", err)
	} else {
		fmt.Println("Test 1 passed")
	}

	// Test 2: WATCH for optimistic locking
	fmt.Println("\nTest 2: WATCH for optimistic locking")
	if err := testWatch(c); err != nil {
		log.Printf("Test 2 failed: %v", err)
	} else {
		fmt.Println("Test 2 passed")
	}

	// Test 3: Transaction Pipeline
	fmt.Println("\nTest 3: Transaction Pipeline")
	if err := testTransactionPipeline(c); err != nil {
		log.Printf("Test 3 failed: %v", err)
	} else {
		fmt.Println("Test 3 passed")
	}
}

func testBasicTransaction(client *r.RedisClient) error {
	err := client.Multi()
	if err != nil {
		return fmt.Errorf("MULTI failed: %v", err)
	}

	_, err = client.Do("SET", "key1", "value1")
	if err != nil {
		client.Discard()
		return fmt.Errorf("SET failed: %v", err)
	}

	_, err = client.Do("SET", "key2", "value2")
	if err != nil {
		client.Discard()
		return fmt.Errorf("SET failed: %v", err)
	}

	results, err := client.Exec()
	if err != nil {
		return fmt.Errorf("EXEC failed: %v", err)
	}

	fmt.Printf("Transaction results: %v\n", results)

	// Verify the results
	val1, err := client.Get("key1")
	if err != nil || val1 != "value1" {
		return fmt.Errorf("unexpected value for key1: %v", val1)
	}

	val2, err := client.Get("key2")
	if err != nil || val2 != "value2" {
		return fmt.Errorf("unexpected value for key2: %v", val2)
	}

	return nil
}

func testWatch(client *r.RedisClient) error {
	// Set initial value
	_, err := client.Set("watched_key", "initial")
	if err != nil {
		return fmt.Errorf("failed to set initial value: %v", err)
	}

	err = client.Watch("watched_key")
	if err != nil {
		return fmt.Errorf("WATCH failed: %v", err)
	}

	// Simulate another client changing the value
	_, err = client.Set("watched_key", "changed")
	if err != nil {
		return fmt.Errorf("failed to change value: %v", err)
	}

	err = client.Multi()
	if err != nil {
		return fmt.Errorf("MULTI failed: %v", err)
	}

	_, err = client.Do("SET", "watched_key", "transaction_value")
	if err != nil {
		client.Discard()
		return fmt.Errorf("SET failed: %v", err)
	}

	results, err := client.Exec()
	if err != nil {
		return fmt.Errorf("EXEC failed: %v", err)
	}

	if results != nil {
		return fmt.Errorf("transaction should have failed due to WATCH, but it succeeded")
	}

	// Verify the value hasn't changed
	val, err := client.Get("watched_key")
	if err != nil || val != "changed" {
		return fmt.Errorf("unexpected value for watched_key: %v", val)
	}

	fmt.Println("WATCH test successful: transaction was aborted as expected")
	return nil
}

func testTransactionPipeline(client *r.RedisClient) error {
	pipeline := client.TxPipeline()

	pipeline.Queue("SET", "pipeline_key1", "pipeline_value1")
	pipeline.Queue("INCR", "pipeline_counter")
	pipeline.Queue("GET", "pipeline_key1")

	results, err := pipeline.Exec()
	if err != nil {
		return fmt.Errorf("pipeline EXEC failed: %v", err)
	}

	fmt.Printf("Pipeline results: %v\n", results)

	// Verify the results
	val1, err := client.Get("pipeline_key1")
	if err != nil || val1 != "pipeline_value1" {
		return fmt.Errorf("unexpected value for pipeline_key1: %v", val1)
	}

	counter, err := client.Get("pipeline_counter")
	if err != nil || counter != "1" {
		return fmt.Errorf("unexpected value for pipeline_counter: %v", counter)
	}

	return nil
}
