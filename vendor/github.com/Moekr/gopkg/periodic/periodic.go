package periodic

import (
	"sync"
	"time"
)

const (
	FixedDelay  DurationMode = 1 // d = 3 |1000|1111000|11000|
	FixedRate   DurationMode = 2 // d = 3 |100|111|210|
	MinInterval DurationMode = 3 // d = 3 |100|1111|110|
)

type DurationMode int

type periodic struct {
	fn     func()
	durFn  func() time.Duration
	mode   DurationMode
	once   *sync.Once
	stopCh chan struct{}
}

func NewPeriodic(fn func(), durFn func() time.Duration, mode DurationMode) *periodic {
	return &periodic{
		fn:     fn,
		durFn:  durFn,
		mode:   mode,
		once:   &sync.Once{},
		stopCh: make(chan struct{}, 1),
	}
}

func NewStaticPeriodic(fn func(), duration time.Duration, mode DurationMode) *periodic {
	return NewPeriodic(fn, func() time.Duration { return duration }, mode)
}

func (p *periodic) Start() {
	p.once.Do(p.start)
}

func (p *periodic) start() {
	switch p.mode {
	case FixedDelay:
		go p.startFixedDelay()
	case FixedRate:
		go p.startFixedRate()
	case MinInterval:
		go p.startMinInterval()
	}
}

func (p *periodic) startFixedDelay() {
	timer := time.NewTimer(p.durFn())
	defer timer.Stop()
	for {
		p.fn()
		timer.Reset(p.durFn())
		select {
		case <-timer.C:
		case <-p.stopCh:
			p.stopCh <- struct{}{}
			return
		}
	}
}

func (p *periodic) startFixedRate() {
	duration := p.durFn()
	ticker := time.NewTicker(duration)
	defer func() {
		ticker.Stop()
	}()
	go p.fn()
	for {
		select {
		case <-ticker.C:
			go p.fn()
			if newDuration := p.durFn(); newDuration != duration {
				ticker.Stop()
				ticker = time.NewTicker(newDuration)
				duration = newDuration
			}
		case <-p.stopCh:
			p.stopCh <- struct{}{}
			return
		}
	}
}

func (p *periodic) startMinInterval() {
	timer := time.NewTimer(p.durFn())
	defer timer.Stop()
	for {
		timer.Reset(p.durFn())
		p.fn()
		select {
		case <-p.stopCh:
			p.stopCh <- struct{}{}
			return
		case <-timer.C:
		}
	}
}

func (p *periodic) Stop() chan struct{} {
	p.stopCh <- struct{}{}
	return p.stopCh
}
