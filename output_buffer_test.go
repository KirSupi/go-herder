package herder

import (
	"strconv"
	"testing"
)

func TestOutputBuffer(t *testing.T) {
	t.Run("TestUnlimitedOutputBuffer", TestUnlimitedOutputBuffer)
	t.Run("TestLimitedOutputBuffer", TestLimitedOutputBuffer)
	t.Run("TestLimitedOutputBuffer", TestDisabledOutputBuffer)
}

func TestUnlimitedOutputBuffer(t *testing.T) {
	b := newOutputBuffer(-1)
	var previousLen int
	for i := 0; i <= 10; i++ {
		previousLen = len(b.buffer)
		iString := strconv.Itoa(i)
		n, err := b.Write([]byte(iString))
		if err != nil {
			t.Errorf("b.Write(s) err != nil: %s", err.Error())
		}
		if n != len(iString) {
			t.Errorf("n, err := b.Write(s); n != len(s); should be %d, got %d", len(iString), n)
		}
		if len(b.buffer) <= previousLen {
			t.Errorf("len(b.buffer) <= previousLen")
		}
	}
}

func TestDisabledOutputBuffer(t *testing.T) {
	b := newOutputBuffer(0)
	var previousLen int
	for i := 0; i <= 10; i++ {
		previousLen = len(b.buffer)
		iString := strconv.Itoa(i)
		n, err := b.Write([]byte(iString))
		if err != nil {
			t.Errorf("b.Write(s) err != nil: %s", err.Error())
		}
		if n != 0 {
			t.Errorf("n, err := b.Write(s); n != len(s); should be 0, got %d", n)
		}
		if len(b.buffer) != previousLen {
			t.Errorf("len(b.buffer) != previousLen")
		}
	}
}

func TestLimitedOutputBuffer(t *testing.T) {
	maxLen := 10
	b := newOutputBuffer(maxLen)
	for i := 0; i <= maxLen+10; i++ {
		iString := strconv.Itoa(i)
		n, err := b.Write([]byte(iString))
		if err != nil {
			t.Errorf("b.Write(s) err != nil: %s", err.Error())
		}
		if n != len(iString) {
			t.Errorf("n, err := b.Write(s); n != len(s); should be %d, got %d", len(iString), n)
		}
		if len(b.buffer) > maxLen {
			t.Errorf("len(b.buffer) > maxLen; should be %d, got %d", maxLen, len(b.buffer))
		}
	}
	var tooLong string
	for i := 0; i < maxLen+10; i++ {
		tooLong += strconv.Itoa(i)
	}
	n, err := b.Write([]byte(tooLong))
	if err != nil {
		t.Errorf("b.Write(tooLong) err != nil: %s", err.Error())
	}
	if n == len(tooLong) {
		t.Errorf("n, err := b.Write(tooLong); n != len(tooLong); should be %d (maxLen), got %d", maxLen, n)
	}
	if len(b.buffer) > maxLen {
		t.Errorf("len(b.buffer) > maxLen; should be %d, got %d", maxLen, len(b.buffer))
	}
}
