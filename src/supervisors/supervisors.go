package supervisors

import "time"

const cWorkerStarted string = "Started"
const cWorkerPaniced string = "Paniced"
const cWorkerExited string = "Exited"

//type generalFunc func(...interface{}) interface{}
type generalFunc func()
type logFunc func(...interface{})

type Worker struct {
	name          string
	lastErrMsg    interface{}
	reportingChan chan string
}

type Supervisor struct {
	name       string
	workerList map[string]Worker
	exited     bool
	channel    chan string
	log        logFunc
}

func (s *Supervisor) monitor() {
	for !s.exited {
		for _, v := range s.workerList {
			select {
			case workerReport := <-v.reportingChan:
				s.log(v.name, workerReport)
				if workerReport == cWorkerPaniced || workerReport == cWorkerExited {
					delete(s.workerList, v.name)
				}

			default:
			}
		}

		select {
		case command := <-s.channel:
			if command == "exit" {
				s.exited = true
				s.log(s.name + " exit")
			}
		default:
		}
		time.Sleep(time.Second * 1)
	}

}

func New(supervisorName string, fn logFunc) *Supervisor {
	s := &Supervisor{
		name:       supervisorName,
		workerList: make(map[string]Worker),
		exited:     false,
		channel:    make(chan string),
		log:        fn,
	}

	go func() {
		s.monitor()
	}()

	return s
}

func (s Supervisor) Exit() {
	s.channel <- "exit"
}

func (s Supervisor) StartWorker(name string, fn generalFunc, params ...interface{}) {

	worker := s.makeWorker(name)
	go func() {
		worker.reportingChan <- cWorkerStarted
		//fn(params)
		fn()

		defer func() {
			if r := recover(); r != nil {
				worker.reportingChan <- cWorkerPaniced
				worker.lastErrMsg = r
			}
			worker.reportingChan <- cWorkerExited
		}()
	}()
}

func (s Supervisor) makeWorker(workerName string) Worker {
	channel := make(chan string)
	w := Worker{
		name:          workerName,
		reportingChan: channel,
	}
	s.workerList[workerName] = w
	return w
}
