package docktest

import (
	"crypto/md5"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type ContainerConfig struct {
	TestName      *string           // eg. TestUserCreate
	ImageName     string            // eg. redis, redis:latest, redis:alpine
	LocalPort     int               // eg. 6379, this is the port that the container will use to expose the image to you.
	ContainerPort int               // eg. 6379, this is the internal port of the container.
	Environment   map[string]string // eg. {"REDIS_PASSWORD": "123456"}
}

type Container interface {
	Stop()
	IsRunning() bool
	Reload()
}

type ContainerData struct {
	Name string
}

func StartContainer(config ContainerConfig) *ContainerData {
	hash := md5.Sum([]byte(
		fmt.Sprintf("%s-%v--%v", config.ImageName, time.Now(), &config.TestName),
	))
	name := fmt.Sprintf("docktest-%x", hash)

	cmd := exec.Command("docker", "run", "-d", "--rm",
		"-p", fmt.Sprintf("%d:%d", config.LocalPort, config.ContainerPort),
		"--name", name)

	if len(config.Environment) > 0 {
		for k, v := range config.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("-e %s=%s", k, v))
		}
	}

	cmd.Args = append(cmd.Args, config.ImageName)

	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("Container started: %s", name))

	return &ContainerData{Name: name}
}

func (c ContainerData) Stop() {
	cmd := exec.Command("docker", "stop", c.Name)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("Container %s stopped", c.Name))
}

func (c ContainerData) IsRunning() bool {
	cmd := exec.Command("docker", "inspect", c.Name)
	err := cmd.Run()
	if err != nil {
		return false
	}
	log.Println(fmt.Sprintf("Container %s is running", c.Name))
	return true
}

func (c ContainerData) Reload() {
	cmd := exec.Command("docker", "restart", c.Name)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("Container %s reloaded", c.Name))
}
