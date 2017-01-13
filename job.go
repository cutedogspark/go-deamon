package gworker

import (
	"fmt"
	"time"
	"strings"
	"os/exec"
	"errors"
)
var (
	ErrCmdNotFound              = errors.New("exec binary not exist")
	ErrCmdExecFail              = errors.New("exec fail")
)

const layout = "2006-01-02 15:04:05"

type JobType string

const (
	JobTypeCMD JobType = "cmd"
)

type Job struct {
	server  string
	jobType JobType
	context string
	delay   time.Duration
	runTime time.Time
	stopped chan bool
}

// Create a new job with the time interval.
func NewJob() Job {
	return Job{
		jobType: JobTypeCMD, // default command line
		context: "",
		delay:   0 * time.Second,
		stopped: make(chan bool, 1),
	}
}


func command(str string) (error) {
	fmt.Println("cmd: ", str)
	strArray := strings.Split(str, " ")
	binary, lookErr := exec.LookPath(strArray[0])
	if lookErr != nil {
		return ErrCmdNotFound
	}
	args := strArray[1:]

	cmd := exec.Command(binary, args...)
	data, err := cmd.Output()
	if err != nil {
		return ErrCmdExecFail
	}
	fmt.Println(string(data))
	return nil
}



func (s *Job) Start() {
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				if s.runTime.Before(time.Now().UTC()) {
					fmt.Println("Run At:", time.Now().UTC().Format(layout))
					command(s.context)
					s.stopped <- true
				}
			case <-s.stopped:
				return
			}
		}
	}()
}

func (s *Job) Stop() {
	s.stopped <- true
}

func (s *Job) At(str string) bool{
	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println("@@ Job At Error :", err)
		return false
	}
	s.runTime = t
	return true
}

func (s *Job) SetDelay(t time.Duration) {
	s.delay = t
}

func (s *Job) GetDelay() time.Duration {
	return s.delay
}

func (s *Job) SetJobType(t JobType){
	s.jobType = t
}

func (s *Job) GetJobType() time.Duration {
	return s.delay
}

func (s *Job) SetJobContext(t string){
	s.context = t
}

func (s *Job) GetJobContext() string {
	return s.context
}

