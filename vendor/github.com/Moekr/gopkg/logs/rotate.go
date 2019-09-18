package logs

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
)

type rotateHook struct {
	logsPath string
	curDay   int
	curOut   io.WriteCloser
	lock     *sync.Mutex
}

func newRotateHook(logsPath string) *rotateHook {
	return &rotateHook{
		logsPath: logsPath,
		lock:     &sync.Mutex{},
	}
}

func (h *rotateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *rotateHook) Fire(e *logrus.Entry) error {
	nowDay := e.Time.Day()
	if nowDay == h.curDay {
		return nil
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	if nowDay == h.curDay {
		return nil
	}
	return h.fire(e, nowDay)
}

func (h *rotateHook) fire(e *logrus.Entry, nowDay int) error {
	if err := os.MkdirAll(h.logsPath, 0777); err != nil {
		return err
	}
	p := path.Join(h.logsPath, e.Time.Format("2006-01-02")+".log")
	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	w := newWriter(f)
	e.Logger.Out = w
	if h.curOut != nil {
		_ = h.curOut.Close()
	}
	h.curDay, h.curOut = nowDay, w
	return nil
}
