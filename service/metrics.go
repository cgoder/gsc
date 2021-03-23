package service

import (
	"sync"
	"time"
)

type TaskMetric struct {
	ID      string
	Statu   TaskStatusType `json:"statu"`
	Percent float64        `json:"percent"`
	Speed   string         `json:"speed"`
	FPS     float64        `json:"fps"`
	Err     string         `json:"err,omitempty"`
}

type gscTaskMetrics struct {
	m       sync.RWMutex
	metrics map[string]TaskMetric
}

var (
	metricMap   = gscTaskMetrics{metrics: make(map[string]TaskMetric)}
	metricIntvl = time.Duration(1)
)

func MetricsCollect() {

	tick := time.NewTicker(time.Second * metricIntvl)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			//TODO: write to clients
			// log.Debugln("send metrics. cnt:", len(metricMap.metrics))
			// metricMap.m.Lock()
			for _, mt := range metricMap.metrics {
				// if mt.Statu != taskStatusDoing {
				// 	continue
				// }
				sendMetric2Clients(mt)
			}
			// metricMap.m.Unlock()

		default:
			time.Sleep(time.Millisecond * 10)
			//TODO:check task status. clean metric.
			// for id, _ := range metricMap.metrics {
			// 	if sts, err := taskMap.TaskStatusGet(id); err == nil {
			// 		if sts != taskStatusDoing {
			// 			delete(metricMap.metrics, id)
			// 		}
			// 	}
			// }
		}

	}

}

func (tmt *gscTaskMetrics) MetricsGet(tid string) (TaskMetric, error) {
	tmt.m.Lock()
	defer tmt.m.Unlock()

	mt, ok := tmt.metrics[tid]
	if ok {
		return mt, nil
	}

	return mt, ErrorTaskNotFound
}

func (tmt *gscTaskMetrics) MetricsSet(tid string, mt TaskMetric) error {
	tmt.m.Lock()
	defer tmt.m.Unlock()

	tmt.metrics[tid] = mt

	return nil
}

func (tmt *gscTaskMetrics) MetricsAdd(tid string) error {
	tmt.m.Lock()
	defer tmt.m.Unlock()

	var mt TaskMetric
	tmt.metrics[tid] = mt
	return nil
}

func (tmt *gscTaskMetrics) MetricsRemove(tid string) error {
	tmt.m.Lock()
	defer tmt.m.Unlock()

	_, ok := tmt.metrics[tid]
	if ok {
		delete(tmt.metrics, tid)
		return nil
	}

	return ErrorTaskNotFound
}
