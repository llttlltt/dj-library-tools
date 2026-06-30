package sys

import (
	"os"
	"os/exec"
)

// FileSystem defines the interface for common OS filesystem operations.
type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (*os.File, error)
}

// OSFileSystem is the real implementation of FileSystem using the 'os' package.
type OSFileSystem struct{}

func (f OSFileSystem) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
func (f OSFileSystem) MkdirAll(path string, perm os.FileMode) error { return os.MkdirAll(path, 0755) }
func (f OSFileSystem) Create(name string) (*os.File, error) { return os.Create(name) }

// Runner defines the interface for running external commands.
type Runner interface {
	LookPath(file string) (string, error)
	Run(name string, arg ...string) ([]byte, error)
}

// RealRunner is the real implementation of Runner using 'os/exec'.
type RealRunner struct{}

func (r RealRunner) LookPath(file string) (string, error) { return exec.LookPath(file) }
func (r RealRunner) Run(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}
