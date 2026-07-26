package views

import "errors"

type failWriter struct {
	target int
	count  int
}

func (w *failWriter) Write(p []byte) (n int, err error) {
	if w.count >= w.target {
		return 0, errors.New("simulated custom error")
	}
	w.count++
	return len(p), nil
}
