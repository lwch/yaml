package yaml

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrNested nested more than 10 times
var ErrNested = errors.New("nested hierarchy is too high")

// Decode decode yaml file with include
func Decode(dir string, out interface{}) error {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()
	str, err := render(f, "", 0)
	if err != nil {
		return err
	}
	return yaml.Unmarshal([]byte(str), out)
}

// Render render yaml file with include
func Render(dir string) (string, error) {
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	f, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer f.Close()
	str, err := render(f, "", 0)
	if err != nil {
		return "", err
	}
	return str, nil
}

func render(f *os.File, prefix string, level int) (string, error) {
	if level > 10 {
		return "", ErrNested
	}
	br := bufio.NewReader(f)
	var buf bytes.Buffer
	handle := func(line string) error {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "#include") {
			str, err := replace(f.Name(),
				strings.TrimSpace(strings.TrimPrefix(trim, "#include")),
				prefix+strings.Repeat(" ", spaceCount(line)), level)
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(&buf, str)
			if err != nil {
				return err
			}
			return nil
		}
		_, err := fmt.Fprint(&buf, prefix+line)
		return err
	}
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = handle(line)
				if err != nil {
					return "", err
				}
				if !strings.HasSuffix(line, "\n") {
					buf.WriteRune('\n')
				}
				return buf.String(), nil
			}
			return "", err
		}
		err = handle(line)
		if err != nil {
			return "", err
		}
	}
}

func spaceCount(line string) int {
	for i, ch := range line {
		if ch != ' ' {
			return i
		}
	}
	return len(line)
}

func replace(self, include, prefix string, level int) (string, error) {
	if !filepath.IsAbs(include) {
		include = filepath.Join(filepath.Dir(self), include)
	}
	var files []string
	if _, err := os.Stat(include); err == nil {
		files = []string{include}
	} else {
		files, err = filepath.Glob(include)
		if err != nil {
			return "", err
		}
	}
	var buf bytes.Buffer
	_, err := fmt.Fprintln(&buf, prefix+"#include "+include)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file == self {
			continue
		}
		_, err = fmt.Fprintln(&buf, prefix+"#+++++ "+file+" +++++")
		if err != nil {
			return "", err
		}
		f, err := os.Open(file)
		if err != nil {
			return "", err
		}
		defer f.Close()
		str, err := render(f, prefix, level)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprint(&buf, str)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(&buf, prefix+"#----- "+file+" -----")
		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
