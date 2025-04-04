package snippetd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

const containerdSocket = "/run/containerd/containerd.sock"
const snippetNamespace = "snippetd"

type CodeRuntime struct {
	client *containerd.Client
}

func NewCodeRuntime() (*CodeRuntime, error) {
	client, err := containerd.New(containerdSocket)
	if err != nil {
		return nil, err
	}

	runtime := &CodeRuntime{client: client}

	// create the namespace if it does not exist
	err = runtime.createNamespace()
	if err != nil {
		return nil, err
	}

	// remove all container leftovers
	err = runtime.removeContainers()
	if err != nil {
		return nil, err
	}

	return runtime, nil
}

// creates the namespace if it does not exist
func (runtime *CodeRuntime) createNamespace() error {
	ctx := context.Background()
	namespaces, err := runtime.client.NamespaceService().List(ctx)
	if err != nil {
		return err
	}

	for _, ns := range namespaces {
		if ns == snippetNamespace {
			return nil // Namespace already exists
		}
	}

	// Create the namespace if it doesn't exist
	err = runtime.client.NamespaceService().Create(ctx, snippetNamespace, map[string]string{})
	if err != nil {
		return err
	}

	return nil
}

// removes any containers all containers from the namespace
func (runtime *CodeRuntime) removeContainers() error {
	ctx := namespaces.WithNamespace(context.Background(), snippetNamespace)
	containers, err := runtime.client.Containers(ctx)
	if err != nil {
		return err
	}

	for _, container := range containers {
		err = container.Delete(ctx, containerd.WithSnapshotCleanup)
		if err != nil {
			return err
		}
	}

	return nil
}

func (runtime *CodeRuntime) Execute(executionUuid, sourceCode string, config RuntimeConfig) CodeExecution {
	ctx := namespaces.WithNamespace(context.Background(), snippetNamespace)

	// pull the required container for this runtime
	imageUrl := "docker.io/library/" + config.Container
	containerImage, err := runtime.client.Pull(ctx, imageUrl, containerd.WithPullUnpack)
	if err != nil {
		return CodeExecution{StandardError: fmt.Sprintf("Failed to pull image: %v", err)}
	}

	// create a temporary directory with the provided uuid
	tempDir, err := os.MkdirTemp("", executionUuid)
	if err != nil {
		return CodeExecution{StandardError: fmt.Sprintf("Failed to create temp dir: %v", err)}
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// write the source code to the temp directory
	sourcePath := filepath.Join(tempDir, config.FileName)
	err = os.WriteFile(sourcePath, []byte(sourceCode), 0777)
	if err != nil {
		return CodeExecution{StandardError: fmt.Sprintf("Failed to write source code to file: %v", err)}
	}

	scriptPath := filepath.Join(tempDir, "exec.sh")
	err = os.WriteFile(scriptPath, []byte(config.RunScript), 0777)
	if err != nil {
		return CodeExecution{StandardError: fmt.Sprintf("Failed to write run script to file: %v", err)}
	}

	containerSpec, err := oci.GenerateSpecWithPlatform(ctx, runtime.client, config.Container, nil)
	if err != nil {
		return CodeExecution{StandardError: fmt.Sprintf("Failed to generate container spec: %v", err)}
	}
	containerSpec.Process.Args = []string{"/bin/sh", "/app/exec.sh"}
	containerSpec.Process.Cwd = "/app"

	container, err := runtime.client.NewContainer(ctx, executionUuid,
		containerd.WithNewSnapshot(executionUuid, containerImage),
		containerd.WithSpec(containerSpec))

	if err != nil {
		return CodeExecution{StandardError: fmt.Sprintf("Failed to create container: %v", err)}
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
	// Create pipes to capture stdout and stderr
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return CodeExecution{StandardError: fmt.Sprintf("Failed to create stdout pipe: %v", err)}
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		stdoutReader.Close()
		stdoutWriter.Close()
		return CodeExecution{StandardError: fmt.Sprintf("Failed to create stderr pipe: %v", err)}
	}

	// Configure container I/O with the pipes
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStreams(nil, stdoutWriter, stderrWriter)))
	if err != nil {
		stdoutReader.Close()
		stdoutWriter.Close()
		stderrReader.Close()
		stderrWriter.Close()
		return CodeExecution{StandardError: fmt.Sprintf("Failed to create task: %v", err)}
	}
	defer task.Delete(ctx)

	// Start reading from the pipes in goroutines
	var stdout, stderr string
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go func() {
		data, err := io.ReadAll(stdoutReader)
		if err == nil {
			stdout = string(data)
		}
		close(stdoutDone)
	}()

	go func() {
		data, err := io.ReadAll(stderrReader)
		if err == nil {
			stderr = string(data)
		}
		close(stderrDone)
	}()

	pid := task.Pid()
	fmt.Printf("Container has process id: %d\n", pid)

	err = task.Start(ctx)
	if err != nil {
		// Close writers to unblock readers
		stdoutWriter.Close()
		stderrWriter.Close()
		// Wait for readers to finish
		<-stdoutDone
		<-stderrDone
		// Close readers
		stdoutReader.Close()
		stderrReader.Close()
		return CodeExecution{StandardError: fmt.Sprintf("Task start error: %v", err)}
	}

	status, err := task.Wait(ctx)
	if err != nil {
		// Close writers to unblock readers
		stdoutWriter.Close()
		stderrWriter.Close()
		// Wait for readers to finish
		<-stdoutDone
		<-stderrDone
		// Close readers
		stdoutReader.Close()
		stderrReader.Close()
		return CodeExecution{StandardError: fmt.Sprintf("Task wait error: %v", err)}
	}

	exitStatus := <-status

	// Close writers to signal EOF to readers
	stdoutWriter.Close()
	stderrWriter.Close()

	// Wait for the reading goroutines to finish
	<-stdoutDone
	<-stderrDone

	// Close readers
	stdoutReader.Close()
	stderrReader.Close()

	return CodeExecution{
		Uuid:           executionUuid,
		ExitCode:       int(exitStatus.ExitCode()),
		StandardOutput: stdout,
		StandardError:  stderr,
	}
}

func (runtime *CodeRuntime) Close() {
	if runtime.client != nil {
		runtime.client.Close()
	}
}
