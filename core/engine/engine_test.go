package engine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	engine1 := NewEngine(BaseEngineConfig{
		Rate:     1,
		Duration: time.Millisecond * 500,
	})
	assert.NotNil(t, engine1)

	engine2 := NewEngine(BaseEngineConfig{
		Rate:     101,
		Duration: time.Second * 1,
	})
	assert.NotNil(t, engine2)

	engine1.Run(func() {})

	engine1.Close()

	engine2.Run(func() {})
}

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
