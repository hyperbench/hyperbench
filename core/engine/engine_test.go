package engine

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go test1(ctx)
	go test2(ctx)
	cancel()
	fmt.Println("main1")
	cancel()
	fmt.Println("main2")
	time.Sleep(time.Second)
	fmt.Println("main3")
}

func test1(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("test1")
			return
		}
	}
}

func test2(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("test2")
			return
		}
	}
}
