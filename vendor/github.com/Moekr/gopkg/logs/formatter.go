package logs

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type formatter struct {
	host string
}

func newFormatter() *formatter {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	return &formatter{
		host: host,
	}
}

func (f *formatter) Format(e *logrus.Entry) ([]byte, error) {
	ss := []string{
		strings.ToUpper(e.Level.String()),
		e.Time.Format("2006-01-02 15:04:05.000"),
		path.Base(e.Caller.File) + ":" + strconv.Itoa(e.Caller.Line),
		f.host,
		e.Message,
	}
	return []byte(strings.Join(ss, " ") + "\n"), nil
}
