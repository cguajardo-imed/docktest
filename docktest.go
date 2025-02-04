package docktest

import (
	"crypto/md5"
	"fmt"
	"os/exec"
	"time"
)

type ContainerConfig struct {
	TestName      *string           // eg. TestUserCreate
	ImageName     string            // eg. redis, redis:latest, redis:alpine
	LocalPort     int               // eg. 6379, this is the port that the container will use to expose the image to you.
	ContainerPort int               // eg. 6379, this is the internal port of the container.
	Environment   map[string]string // eg. {"REDIS_PASSWORD": "123456"}
	Command       string            // eg. --bind_ip_all or sh -c "redis-server --bind_ip_all"
}

type Container interface {
	Stop()
	IsRunning() bool
	Reload()

	// Getters
	GetName() string
	GetLocalPort() int
}

type ContainerData struct {
	Name      string
	LocalPort int
}

func StartContainer(config ContainerConfig) *ContainerData {
	hash := md5.Sum([]byte(
		fmt.Sprintf("%s-%v--%v", config.ImageName, time.Now(), &config.TestName),
	))
	name := fmt.Sprintf("docktest-%x", hash)

	cmd := exec.Command("docker", "run", "-d", "--rm",
		"-p", fmt.Sprintf("%d:%d", config.LocalPort, config.ContainerPort),
		"--name", name, config.Command)

	if len(config.Environment) > 0 {
		for k, v := range config.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("-e %s=%s", k, v))
		}
	}

	cmd.Args = append(cmd.Args, config.ImageName)

	err := cmd.Start()
	if err != nil {
		Error("Error::Start:", err.Error())
		return nil
	}
	_ = cmd.Wait()

	cd := &ContainerData{Name: name, LocalPort: config.LocalPort}
	if cd.IsRunning() {
		Success(fmt.Sprintf("Container %s started successfully", name))
	}

	return cd
}

func (c ContainerData) Stop() {
	cmd := exec.Command("docker", "stop", c.Name)
	err := cmd.Run()
	if err != nil {
		Error("Error::Stop:", err.Error())
	} else {
		Success(fmt.Sprintf("Container %s stopped", c.Name))
	}
}

func (c ContainerData) IsRunning() bool {
	cmd := exec.Command("docker", "inspect", c.Name)
	err := cmd.Run()
	if err != nil {
		return false
	}
	Info(fmt.Sprintf("Container %s is running", c.Name))
	return true
}

func (c ContainerData) Reload() {
	cmd := exec.Command("docker", "restart", c.Name)
	err := cmd.Run()
	if err != nil {
		Error("Error::Reload:", fmt.Sprintf("can't restart the container %s", c.Name))
	} else {
		Info(fmt.Sprintf("Container %s reloaded", c.Name))
	}
}

func (c ContainerData) GetName() string {
	return c.Name
}

func (c ContainerData) GetLocalPort() int {
	return c.LocalPort
}
