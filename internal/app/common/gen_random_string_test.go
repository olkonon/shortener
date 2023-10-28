package common

import (
	"github.com/sirupsen/logrus"
	"io"
	"testing"
)

func init() {
	logrus.SetOutput(io.Discard)
}

func TestGenRandomString(t *testing.T) {
	expectedLen := 16
	uniqMap := make(map[string]bool)
	for i := 0; i < 1000000; i++ {
		tmp := GenRandomString(expectedLen)
		if len(tmp) != expectedLen {
			t.Fail()
		}

		if _, ok := uniqMap[tmp]; ok {
			t.Fail()
		} else {
			uniqMap[tmp] = true
		}

	}
}
