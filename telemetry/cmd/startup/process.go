package startup

import (
	"context"
	"fmt"
	"go_agent/config"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ManagedProcess struct {
	config config.StartupProcessConfig
	cmd    *exec.Cmd
}

func StartEnabled(ctx context.Context, processesConfig config.StartupProcessesConfig, wg *sync.WaitGroup) ([]*ManagedProcess, error) {
	processes := make([]*ManagedProcess, 0, len(processesConfig.Processes))

	for _, processConfig := range processesConfig.Processes {
		if !processConfig.Enabled {
			continue
		}

		process, err := StartProcess(ctx, processConfig, wg)
		if err != nil {
			if processConfig.Required {
				return processes, err
			}
			log.Printf("Optional startup process %q failed to start: %v", processConfig.Name, err)
			continue
		}

		processes = append(processes, process)
	}

	return processes, nil
}

func StartProcess(ctx context.Context, processConfig config.StartupProcessConfig, wg *sync.WaitGroup) (*ManagedProcess, error) {
	if strings.TrimSpace(processConfig.Name) == "" {
		return nil, fmt.Errorf("startup process name is required")
	}

	if strings.TrimSpace(processConfig.Command) == "" {
		return nil, fmt.Errorf("startup process %q command is required", processConfig.Name)
	}

	cmd := exec.CommandContext(ctx, processConfig.Command, processConfig.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil && err != syscall.ESRCH {
			return err
		}
		return nil
	}
	cmd.WaitDelay = 5 * time.Second

	if strings.TrimSpace(processConfig.WorkingDir) != "" {
		cmd.Dir = processConfig.WorkingDir
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start startup process %q: %w", processConfig.Name, err)
	}

	process := &ManagedProcess{
		config: processConfig,
		cmd:    cmd,
	}

	log.Printf("Started startup process %q with pid %d", processConfig.Name, cmd.Process.Pid)

	wg.Add(1)
	go process.wait(ctx, wg)

	return process, nil
}

func (p *ManagedProcess) wait(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	err := p.cmd.Wait()
	if ctx.Err() != nil {
		return
	}

	if err != nil {
		log.Printf("Startup process %q exited unexpectedly with error: %v", p.config.Name, err)
		return
	}

	log.Printf("Startup process %q exited unexpectedly", p.config.Name)
}
