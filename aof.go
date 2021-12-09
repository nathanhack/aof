package aof

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	arraysPrefix        = "*"
	simpleStringsPrefix = "+"
	bulkStringPrefix    = "$"
	errorsPrefix        = "-"
	integerPrefix       = ":"
)

type Command struct {
	Name      string
	Arguments []string
}

//ReadCommand reads a command from the bufio.Reader.  It returns the Command and the bytes used to
// create the Command.  In the case of an error it will return the current bytes read and the error.
func ReadCommand(reader *bufio.Reader) (*Command, []byte, error) {
	//for all command we expect it to contain an array with the first element
	// to be the key word (bulk string)
	buf := &bytes.Buffer{}
	str, err := readline(reader, buf)
	if err != nil {
		return nil, buf.Bytes(), err
	}

	if len(str) <= 1 {
		return nil, buf.Bytes(), errors.New("corrupt data beginning of command too small")
	}

	if string(str[0]) != arraysPrefix {
		return nil, buf.Bytes(), fmt.Errorf("command not found, expected an array instead found %v", str)
	}

	count, err := strconv.Atoi(str[1:])
	if err != nil {
		return nil, buf.Bytes(), fmt.Errorf("array value not")
	}

	if count < 1 {
		return nil, buf.Bytes(), fmt.Errorf("command required at least one value to be present")
	}

	//now time to read the key word
	command, err := readBulkString(reader, buf)
	if err != nil {
		return nil, buf.Bytes(), err
	}

	if _, has := validCommands[strings.ToUpper(command)]; !has {
		return nil, buf.Bytes(), fmt.Errorf("expected to find a command instead found %v", command)
	}

	args := make([]string, 0)
	for i := 1; i < count; i++ {
		arg, err := readBulkString(reader, buf)
		if err != nil {
			return nil, buf.Bytes(), fmt.Errorf("command %v had an error %v", command, err)
		}
		args = append(args, arg)
	}

	return &Command{
		Name:      command,
		Arguments: args,
	}, buf.Bytes(), nil
}

func readBulkString(reader *bufio.Reader, buf *bytes.Buffer) (string, error) {
	str, err := readline(reader, buf)
	if err != nil {
		return "", err
	}
	if len(str) == 1 {
		return "", errors.New("corrupt data beginning of bulk string too small")
	}

	//this should be a bulk string
	if string(str[0]) != bulkStringPrefix {
		return "", errors.New("bulk string requires '$' as prefix")
	}

	size, err := strconv.Atoi(str[1:])
	if err != nil {
		return "", fmt.Errorf("the size of the bulk string was not a parsable: %v", err)
	}

	str, err = readline(reader, buf)
	if err != nil {
		return "", err
	}

	// for the case readline returned something smaller than the size
	// as is the case for values containing \r\n
	for len(str) < size && err == nil {
		str += "\r\n" // we add back what the readline removed
		tmp, err := readline(reader, buf)
		if err != nil {
			return "", err
		}
		str += tmp
	}

	if len(str) != size {
		return "", fmt.Errorf("the size of the bulk string was not equal expected %v but found %v", size, len(str))
	}

	return str, nil
}

func readline(reader *bufio.Reader, buf *bytes.Buffer) (string, error) {
	var err error
	var str string
	for !strings.HasSuffix(str, "\r\n") && !errors.Is(err, io.EOF) {
		tmp, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return "", err
			} else {
				if len(str)+len(tmp) == 0 {
					return "", err
				}
			}
		}
		str += tmp
		buf.WriteString(tmp)
	}

	if !strings.HasSuffix(str, "\r\n") {
		return "", fmt.Errorf("bad formatted expected suffix \"\\r\\n\" for '%v'", str)
	}

	//remove the suffix
	return strings.TrimSuffix(str, "\r\n"), nil
}

//WriteCommand will write a command to the writer via one call to writer.Write().  The caller is
// responsable for calling the writer's flush() if needed.
func WriteCommand(command *Command, writer *bufio.Writer) error {
	buf := &bytes.Buffer{}

	err := writeArrayPrefix(1+len(command.Arguments), buf)
	if err != nil {
		return err
	}

	err = writeBulkString(command.Name, buf)
	if err != nil {
		return err
	}

	for i := 0; i < len(command.Arguments); i++ {
		err = writeBulkString(command.Arguments[i], buf)
		if err != nil {
			return err
		}
	}

	bs := buf.Bytes()
	n, err := writer.Write(bs)
	if err != nil {
		return err
	}

	if n != len(bs) {
		return fmt.Errorf("expected %v bytes written found %v", len(bs), n)
	}

	return nil
}

func writeArrayPrefix(count int, writer *bytes.Buffer) error {
	arrayStr := fmt.Sprintf("%v%v\r\n", arraysPrefix, count)
	n, err := writer.WriteString(arrayStr)
	if err != nil {
		return err
	}
	if len(arrayStr) != n {
		return fmt.Errorf("write failed expected to write %v bytes but found %v", len(arrayStr), n)
	}
	return nil
}

func writeBulkString(str string, writer *bytes.Buffer) error {
	bulkStr := fmt.Sprintf("%v%v\r\n%v\r\n", bulkStringPrefix, len(str), str)

	n, err := writer.WriteString(bulkStr)
	if err != nil {
		return err
	}
	if len(bulkStr) != n {
		return fmt.Errorf("write failed expected to write %v bytes but found %v", len(bulkStr), n)
	}
	return nil
}
