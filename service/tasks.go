package service

import (
	"context"
	"sync"

	"github.com/cgoder/gsc/common"
	"github.com/cgoder/gsc/ffmpeg"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	//task chan
	taskCh = make(chan Message, 1)

	taskMap = gscTask{tasks: make(map[string]*Task)}
)

type TaskStatusType int

const (
	taskStatusTodo TaskStatusType = iota
	taskStatusDoing
	taskStatusDone
	taskStatusCancel
	taskStatusFail
)

// type Contx struct {
// 	ctx    context.Context
// 	cancel context.CancelFunc
// }

type FFInfo struct {
	ffmpeg  *ffmpeg.FFmpeg
	ffProbe ffmpeg.FFProbeResponse
	ctx     context.Context
	cancel  context.CancelFunc
}
type Task struct {
	ID string
	// Statu  TaskStatusType
	cmd    Message
	ff     FFInfo
	metric TaskMetric
}

type gscTask struct {
	m     sync.RWMutex
	tasks map[string]*Task
}

func (t *gscTask) TaskAdd(msg Message) (string, error) {
	t.m.Lock()
	defer t.m.Unlock()

	for _, task := range t.tasks {
		if msg.Type == task.cmd.Type && msg.Input == task.cmd.Input && msg.Output == task.cmd.Output && msg.Payload == task.cmd.Payload {
			return "", ErrorTaskExisit
		}
	}

	var task Task
	var err error
	tid := uuid.Must(uuid.NewV4(), err).String()
	task.ID = tid

	t.tasks[tid] = &task
	t.tasks[tid].cmd = msg
	// t.tasks[tid].Statu = taskStatusTodo

	metricMap.MetricsAdd(tid)

	return tid, nil
}

func (t *gscTask) TaskDeleteByMsg(msg Message) error {
	t.m.Lock()
	defer t.m.Unlock()

	for tid, task := range t.tasks {
		if msg.Type == task.cmd.Type && msg.Input == task.cmd.Input && msg.Output == task.cmd.Output && msg.Payload == task.cmd.Payload {
			//TODO: stop ffmpeg task.
			metricMap.MetricsRemove(tid)
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
		metricMap.MetricsRemove(tid)
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

func (t *gscTask) TaskFFGet(tid string) (FFInfo, error) {
	var ff FFInfo
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		ff = t.tasks[tid].ff
		return ff, nil
	}

	return ff, ErrorTaskUpdateFail
}

func (t *gscTask) TaskFFSet(tid string, ff FFInfo) error {
	t.m.Lock()
	defer t.m.Unlock()

	if _, ok := t.tasks[tid]; ok {
		t.tasks[tid].ff = ff
		return nil
	}

	return ErrorTaskUpdateFail
}

func HandleTaskMessages() {
	//collect all task metric.
	go MetricsCollect()

	//read msg from client, add task and exec.
	for msg := range taskCh {
		log.Infoln("Got cmd: ", common.JsonFormat(msg))
		switch msg.Type {
		case prefixStart:
			if tid, err := taskMap.TaskAdd(msg); err == nil {
				log.Debugln("task start. tid: ", tid)
				err := runProcess(tid, msg.Input, msg.Output, msg.Payload)
				if err != nil {
					if err := taskMap.TaskDelete(tid); err != nil {
						log.Errorln("task delete fail! ", tid)
					}
					log.Errorln("process task fail! ", tid, err.Error())
				}
			} else {
				log.Errorln("add task fail. ", common.JsonFormat(msg))
			}

		case prefixStop, prefixCancel:
			if tid, err := taskMap.TaskGet(msg); err == nil {
				stopProcess(tid)
				if err := taskMap.TaskDelete(tid); err != nil {
					log.Errorln("task delete fail! ", tid)
				}
				log.Debugln("task remove success. tid: ", tid)
			} else {
				log.Errorln("remove task fail. ", tid, err)
			}
		}
	}

}
