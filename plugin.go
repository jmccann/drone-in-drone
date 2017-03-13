package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// default docker registry
	defaultRegistry = "https://index.docker.io/v1/"
)

type (
	// Daemon defines Docker daemon parameters.
	Daemon struct {
		Registry      string   // Docker registry
		Mirror        string   // Docker registry mirror
		Insecure      bool     // Docker daemon enable insecure registries
		StorageDriver string   // Docker daemon storage driver
		StoragePath   string   // Docker daemon storage path
		Disabled      bool     // DOcker daemon is disabled (already running)
		Debug         bool     // Docker daemon started in debug mode
		Bip           string   // Docker daemon network bridge IP address
		DNS           []string // Docker daemon dns server
		MTU           string   // Docker daemon mtu setting
		IPv6          bool     // Docker daemon IPv6 networking
		Experimental  bool     // Docker daemon enable experimental mode
	}

	// Exec defines drone exec build parameters.
	DroneExec struct {
		Context    string   // Docker build context
		Args       []string // Docker build args
	}

	// Plugin defines the Docker plugin parameters.
	Plugin struct {
		DroneExec DroneExec // Drone exec configuration
		Daemon    Daemon    // Docker daemon configuration
	}
)

// Exec executes the plugin step
func (p Plugin) Exec() error {
	// start the Docker daemon server
	if !p.Daemon.Disabled {
		cmd := commandDaemon(p.Daemon)
		if p.Daemon.Debug {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			cmd.Stdout = ioutil.Discard
			cmd.Stderr = ioutil.Discard
		}
		go func() {
			trace(cmd)
			cmd.Run()
		}()
	}

	// poll the docker daemon until it is started. This ensures the daemon is
	// ready to accept connections before we proceed.
	for i := 0; i < 15; i++ {
		cmd := commandInfo()
		err := cmd.Run()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}

	var cmds []*exec.Cmd
	cmds = append(cmds, commandDockerPrune())          // cleanup docker
	cmds = append(cmds, commandVersion())              // docker version
	cmds = append(cmds, commandInfo())                 // docker info
	cmds = append(cmds, commandDroneExec(p.DroneExec)) // drone exec

	// execute all commands in batch mode.
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

const dockerExe = "/usr/local/bin/docker"
const dockerdExe = "/usr/local/bin/dockerd"
const droneExe = "/bin/drone"

// helper function to create the drone exec command
func commandDroneExec(droneExec DroneExec) *exec.Cmd {
	args := []string{
		"exec",
	}

	cmd := exec.Command(droneExe, args...)
	fmt.Printf("Setting Dir to: %s\n", droneExec.Context)
	cmd.Dir = droneExec.Context

	return cmd
}

// helper function to create the docker info command.
func commandVersion() *exec.Cmd {
	return exec.Command(dockerExe, "version")
}

// helper function to create the docker info command.
func commandInfo() *exec.Cmd {
	return exec.Command(dockerExe, "info")
}

// helper function to get a proxy value from the environment.
//
// assumes that the upper and lower case versions of are the same.
func getProxyValue(key string) string {
	value := os.Getenv(key)

	if len(value) > 0 {
		return value
	}

	return os.Getenv(strings.ToUpper(key))
}

// helper function to create the docker daemon command.
func commandDaemon(daemon Daemon) *exec.Cmd {
	args := []string{"-g", daemon.StoragePath}

	if daemon.StorageDriver != "" {
		args = append(args, "-s", daemon.StorageDriver)
	}
	if daemon.Insecure && daemon.Registry != "" {
		args = append(args, "--insecure-registry", daemon.Registry)
	}
	if daemon.IPv6 {
		args = append(args, "--ipv6")
	}
	if len(daemon.Mirror) != 0 {
		args = append(args, "--registry-mirror", daemon.Mirror)
	}
	if len(daemon.Bip) != 0 {
		args = append(args, "--bip", daemon.Bip)
	}
	for _, dns := range daemon.DNS {
		args = append(args, "--dns", dns)
	}
	if len(daemon.MTU) != 0 {
		args = append(args, "--mtu", daemon.MTU)
	}
	if daemon.Experimental {
		args = append(args, "--experimental")
	}
	return exec.Command(dockerdExe, args...)
}

func commandDockerPrune() *exec.Cmd {
	return exec.Command(dockerExe, "system", "prune", "-f")
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
