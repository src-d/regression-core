package regression

import (
	"os/exec"
	"syscall"
)

// Server struct describes a daemon.
type Server struct {
	cmd *exec.Cmd
}

// NewServer creates a new Server struct.
func NewServer() *Server {
	return new(Server)
}

// Start executes a command in background.
func (s *Server) Start(name string, arg ...string) error {
	s.cmd = exec.Command(name, arg...)

	return s.cmd.Start()
}

// Stop kill the daemon.
func (s *Server) Stop() error {
	err := s.cmd.Process.Kill()
	if err != nil {
		return err
	}

	_ = s.cmd.Wait()
	return nil
}

// Alive checks if the process is still running.
func (s *Server) Alive() bool {
	if s.cmd == nil || s.cmd.Process == nil {
		return false
	}

	err := s.cmd.Process.Signal(syscall.Signal(0))
	return err == nil
}

// Rusage returns usage counters.
func (s *Server) Rusage() *syscall.Rusage {
	rusage, _ := s.cmd.ProcessState.SysUsage().(*syscall.Rusage)
	return rusage
}
