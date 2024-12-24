package main

import (
	"time"

	"github.com/cguajardo-imed/docktest"
)

func main() {
	container := docktest.StartContainer(docktest.ContainerConfig{
		ImageName:     "redis",
		LocalPort:     1414,
		ContainerPort: 6379,
		Environment: map[string]string{
			"REDIS_PASSWORD": "123456",
		},
	})

	time.Sleep(time.Second * 5)
	isRunning := container.IsRunning()
	if isRunning {
		container.Reload()
	}

	time.Sleep(time.Second * 10)
	container.Stop()
}
