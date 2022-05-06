package handlers

import (
	"log"
	"time"
)

type Events interface {
	Start(info CStart)
	StartResult(info CStart, err error)
	Shutdown(info CShutdown)
	ShutdownResult(info CShutdown, err error)
	CPeriod(info CPeriod, t time.Time, enabled bool)
	CPeriodNext(period CPeriod, d time.Duration)
	CPeriodExec(period CPeriod, input interface{})
	CPeriodResult(period CPeriod, err error)
	CPeriodExit(period CPeriod)
	CConsume(info CConsume, input interface{})
	CConsumeNext(info CConsume)
	CConsumeExec(info CConsume, input interface{})
	CConsumeResult(info CConsume, err error)
	CConsumeExit(info CConsume)
	PConsume(info PConsume, input interface{})
	PConsumeResult(info PConsume, err error)
	PConsumeExit(info PConsume)
	PBlockingStart(info PBlocking)
	PBlockingResult(info PBlocking, err error)
}

type NullEvents struct{}

func (NullEvents) Start(info CStart)                               {}
func (NullEvents) StartResult(info CStart, err error)              {}
func (NullEvents) Shutdown(info CShutdown)                         {}
func (NullEvents) ShutdownResult(info CShutdown, err error)        {}
func (NullEvents) CPeriod(info CPeriod, t time.Time, enabled bool) {}
func (NullEvents) CPeriodNext(period CPeriod, d time.Duration)     {}
func (NullEvents) CPeriodExec(period CPeriod, input interface{})   {}
func (NullEvents) CPeriodResult(period CPeriod, err error)         {}
func (NullEvents) CPeriodExit(period CPeriod)                      {}
func (NullEvents) CConsume(info CConsume, input interface{})       {}
func (NullEvents) CConsumeNext(info CConsume)                      {}
func (NullEvents) CConsumeExec(info CConsume, input interface{})   {}
func (NullEvents) CConsumeResult(info CConsume, err error)         {}
func (NullEvents) CConsumeExit(info CConsume)                      {}
func (NullEvents) PConsume(info PConsume, input interface{})       {}
func (NullEvents) PConsumeResult(info PConsume, err error)         {}
func (NullEvents) PConsumeExit(info PConsume)                      {}
func (NullEvents) PBlockingStart(info PBlocking)                   {}
func (NullEvents) PBlockingResult(info PBlocking, err error)       {}

type LogEvents struct {
}

func (LogEvents) Start(info CStart) {
	log.Printf("Start")
}

func (LogEvents) StartResult(info CStart, err error) {
	log.Printf("StartResult")
}

func (LogEvents) Shutdown(info CShutdown) {
	log.Printf("Shutdown")
}

func (LogEvents) ShutdownResult(info CShutdown, err error) {
	log.Printf("ShutdownResult")
}

func (LogEvents) CPeriod(info CPeriod, t time.Time, enabled bool) {
	log.Printf("CPeriod")
}

func (LogEvents) CPeriodNext(period CPeriod, d time.Duration) {
	log.Printf("CPeriodNext")
}

func (LogEvents) CPeriodExec(period CPeriod, input interface{}) {
	log.Printf("CPeriodExec")
}

func (LogEvents) CPeriodResult(period CPeriod, err error) {
	log.Printf("CPeriodResult")
}

func (LogEvents) CPeriodExit(period CPeriod) {
	log.Printf("CPeriodExit")
}

func (LogEvents) CConsume(info CConsume, input interface{}) {
	log.Printf("CConsume")
}
func (LogEvents) CConsumeNext(info CConsume) {
	log.Printf("CConsumeNext")
}

func (LogEvents) CConsumeExec(info CConsume, input interface{}) {
	log.Printf("CConsumeExec")
}

func (LogEvents) CConsumeResult(info CConsume, err error) {
	log.Printf("CConsumeResult")
}

func (LogEvents) CConsumeExit(info CConsume) {
	log.Printf("CConsumeExit")
}

func (LogEvents) PConsume(info PConsume, input interface{}) {
	log.Printf("PConsume")
}

func (LogEvents) PConsumeResult(info PConsume, err error) {
	log.Printf("PConsumeResult")
}

func (LogEvents) PConsumeExit(info PConsume) {
	log.Printf("PConsumeExit")
}

func (LogEvents) PBlockingStart(info PBlocking) {
	log.Printf("PBlockingStart")
}

func (LogEvents) PBlockingResult(info PBlocking, err error) {
	log.Printf("PBlockingResult")
}
