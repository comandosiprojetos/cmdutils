package cmdutils

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type ProcessMonitor struct {
	CmdName *string
	CmdArgs *[]string
	Process *os.Process
	Cmd     *exec.Cmd
	Output  *[]byte
	Err     error
}

type ProcessStateListenerImpl struct {
	monitor chan bool
}

func (processStateListenerImpl *ProcessStateListenerImpl) OnComplete(processMonitor *ProcessMonitor) {
	log.Println("on complete")
	log.Println("output is\n", string(*processMonitor.Output))
	processStateListenerImpl.monitor <- true
}

func (processStateListenerImpl *ProcessStateListenerImpl) OnError(processMonitor *ProcessMonitor, err error) {
	log.Panic("Error encountered", err)
	processStateListenerImpl.monitor <- true
}

type ProcessStateListener interface {
	OnComplete(processMonitor *ProcessMonitor)
	OnError(processMonitor *ProcessMonitor, err error)
}

func CriarForkProcessoWindows(pathProcesso string, argumento string) {
	processStateListenerImpl := &ProcessStateListenerImpl{make(chan bool)}
	fork(processStateListenerImpl, pathProcesso, argumento)

	<-processStateListenerImpl.monitor
}

func fork(processStateListener ProcessStateListener, cmdName string, cmdArgs ...string) {
	go func() {
		processMonitor := &ProcessMonitor{}
		args := strings.Join(cmdArgs, ",")
		command := exec.Command(cmdName, args)
		output, err := command.Output()
		if err != nil {
			processMonitor.Err = err
			processStateListener.OnError(processMonitor, err)
		}
		processMonitor.Output = &output
		processStateListener.OnComplete(processMonitor)
	}()
}
