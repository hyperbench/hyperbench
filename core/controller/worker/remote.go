package worker

import (
	"github.com/hyperbench/hyperbench/core/network/client"
)

// RemoteWorker is the agent of remote worker.
type RemoteWorker struct {
	*client.Client
}

// NewRemoteWorker create RemoteWorker.
func NewRemoteWorker(index int, url string) (*RemoteWorker, error) {
	c := client.NewClient(index, url)
	err := c.Init()
	if err != nil {
		return nil, err
	}
	return &RemoteWorker{
		Client: c,
	}, nil
}
