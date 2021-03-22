package service

import (
	"context"
	"sync"

	"github.com/cgoder/gsc/ffmpeg"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	//task chan
	taskCh = make(chan Message, 1)
)

const (
	taskStatusUnkonw = iota
	taskStatusTodo
	taskStatusDoing
	taskStatusDone
	taskStatusFail
	taskStatusCancel
)

type Contx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type Task struct {
	cmd     Message
	stats   Status
	srcInfo ffmpeg.FFProbeResponse
	ctx     Contx
}

type gscTask struct {
	m     sync.RWMutex
	tasks map[string]*Task
}

var taskMap = gscTask{tasks: make(map[string]*Task)}

func (t *gscTask) TaskAdd(msg Message) (string, error) {
	t.m.Lock()
	defer t.m.Unlock()

	for _, task := range t.tasks {
		if msg.Type == task.cmd.Type && msg.Input == task.cmd.Input && msg.Output == task.cmd.Output && msg.Payload == task.cmd.Payload {
			return "", ErrorTaskExisit
		}
	}
	var err error
	tid := uuid.Must(uuid.NewV4(), err).String()

	t.tasks[tid] = new(Task)
	t.tasks[tid].cmd = msg
	t.tasks[tid].stats.Progress = taskStatusTodo

	return tid, nil
}

func (t *gscTask) TaskDeleteByMsg(msg Message) error {
	t.m.Lock()
	defer t.m.Unlock()

	for tid, task := range t.tasks {
		if msg.Type == task.cmd.Type && msg.Input == task.cmd.Input && msg.Output == task.cmd.Output && msg.Payload == task.cmd.Payload {
			//TODO: stop ffmpeg task.
			delete(t.tasks, tid)
			return nil
		}
	}

	return ErrorTaskExisit
}

func (t *gscTask) TaskDelete(tid string) error {
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		delete(t.tasks, tid)
		return nil
	}

	return ErrorTaskNotFound
}

func (t *gscTask) TaskGet(msg Message) (string, error) {
	t.m.Lock()
	defer t.m.Unlock()

	for tid, task := range t.tasks {
		if msg.Input == task.cmd.Input && msg.Output == task.cmd.Output && msg.Payload == task.cmd.Payload {
			return tid, nil
		}
	}

	return "", ErrorTaskNotFound
}

func (t *gscTask) TaskStatusGet(tid string) (Status, error) {
	var st Status
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		st = t.tasks[tid].stats
		return st, nil
	}

	return st, ErrorTaskNotFound
}

func (t *gscTask) TaskStatusSet(tid string, stats Status) error {
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		t.tasks[tid].stats = stats
		return nil
	}

	return ErrorTaskUpdateFail
}

func (t *gscTask) TaskInfoGet(tid string) (ffmpeg.FFProbeResponse, error) {
	var info ffmpeg.FFProbeResponse
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		info = t.tasks[tid].srcInfo
		return info, nil
	}

	return info, ErrorTaskNotFound
}

func (t *gscTask) TaskInfoSet(tid string, info ffmpeg.FFProbeResponse) error {
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		t.tasks[tid].srcInfo = info
		return nil
	}

	return ErrorTaskUpdateFail
}

func (t *gscTask) TaskCtxGet(tid string) (Contx, error) {
	var ctx Contx
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		ctx = t.tasks[tid].ctx
		return ctx, nil
	}

	return ctx, ErrorTaskUpdateFail
}

func (t *gscTask) TaskCtxSet(tid string, ctx Contx) error {
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		t.tasks[tid].ctx = ctx
		return nil
	}

	return ErrorTaskUpdateFail
}

func handleMessages() {
	for msg := range taskCh {
		log.Infoln("Got cmd: ", JsonFormat(msg))
		switch msg.Type {
		case prefixStart:
			if tid, err := taskMap.TaskAdd(msg); err == nil {
				log.Debugln("task start. tid: ", tid)
				err := runProcess(tid, msg.Input, msg.Output, msg.Payload)
				if err != nil {
					taskMap.TaskDelete(tid)
					log.Errorln("process task fail! ", tid, err.Error())
				}
			} else {
				log.Errorln("add task fail. ", JsonFormat(msg))
			}
		case prefixStop, prefixCancel:
			if tid, err := taskMap.TaskGet(msg); err == nil {
				stopProcess(tid)
				log.Debugln("task remove success.", JsonFormat(msg))
			} else {
				log.Errorln("remove task fail. ", tid, err)
			}
		}
	}

}
