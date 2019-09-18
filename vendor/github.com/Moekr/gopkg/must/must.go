package must

import "io"

func Close(c io.Closer) {
	_ = c.Close()
}
