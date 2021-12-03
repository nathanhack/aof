package aof

import (
	"bufio"
	"bytes"
	"math/rand"
	"reflect"
	"testing"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ{}!@#$%^&*(),./<>?;':\"[]{}\\|`~1234567890-=_+`")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomArgs(count, size int) []string {
	if count == 0 {
		return nil
	}
	result := make([]string, count)

	for i := range result {
		result[i] = randStr(rand.Intn(size))
	}
	return result
}

func TestRandomWriteRead(t *testing.T) {
	for command := range validCommands {
		expected := &Command{
			Name:      command,
			Arguments: randomArgs(5, 20),
		}

		var bb bytes.Buffer
		wbb := bufio.NewWriter(&bb)
		err := WriteCommand(expected, wbb)
		if err != nil {
			t.Fatalf("expected no error but found %v on command: %#v", err, expected)
		}

		rbb := bufio.NewReader(bytes.NewBuffer(bb.Bytes()))
		actual, bs, err := ReadCommand(rbb)
		if err != nil {
			t.Fatalf("expected no error but found %v on command: %#v", err, expected)
		}

		if !reflect.DeepEqual(bs, bb.Bytes()) {
			t.Fatalf("expected %v but found %v", bb.Bytes(), bs)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("expected equal but found %#v != %#v", expected, actual)
		}
	}
}
