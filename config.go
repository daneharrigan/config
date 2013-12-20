// Package config parses and exposes read-only access to INI files
package config

import (
	"os"
	"bufio"
	"errors"
	"strings"
	"io"
)

var (
	NotFound  = errors.New("not found")
	Malformed = errors.New("malformed config file")
)

type Config struct {
	content map[string]map[string]string
}

// Generates a new config struct containing the contents of the INI file. If the
// path provided does not point to a file or the file cannot be properly parsed,
// the `Malformed` error is returned.
func New(path string) (c *Config, err error) {
	c = new(Config)
	c.content, err = parse(path)
	return c, err
}

// Returns the value found under the specified section and key. If the value
// does not exist, the `NotFound` error is returned.
func (c *Config) Get(s, k string) (string, error) {
	if _, ok := c.content[s]; ok {
		if _, ok := c.content[s][k]; ok {
			return c.content[s][k], nil
		}
	}

	return "", NotFound
}

func parse(path string) (map[string]map[string]string, error) {
	content := make(map[string]map[string]string)

	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return content, err
	}

	r := bufio.NewReader(f)
	var section, key, value string

	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return content, nil
			}
			return content, err
		}

		switch b {
		case '\n': // skip
		case ';': // comment
			_, err = r.ReadString('\n')
		case '[':
			// clear
			key = key[:0]
			value = value[:0]

			if section, err = readSection(r); err == nil {
				section = trim(section)
			}
		default:
			if key, err = readKey(r); err == nil {
				key = trim(key)
			}

			if value, err = readValue(r); err == nil {
				value = trim(value)
			}
		}

		if err != nil {
			return content, err
		}

		if section != "" && key != "" && value != "" {
			if _, ok := content[section]; !ok {
				content[section] = make(map[string]string)
			}

			content[section][key] = value
		}
	}
}

func readSection(r *bufio.Reader) (string, error) {
	var v []byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}

		switch b {
		case ']':
			_, err = r.ReadBytes('\n')
			return string(v), err
		default:
			v = append(v, b)
		}
	}

	return "", NotFound
}




func readKey(r *bufio.Reader) (string, error) {
	var v []byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}

		switch b {
		case '=':
			return string(v), nil
		default:
			v = append(v, b)
		}
	}

	return "", Malformed
}

func readValue(r *bufio.Reader) (string, error) {
	var v []byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}

		switch b {
		case '\n':
			return string(v), err
		case ';': // comment
			_, err = r.ReadBytes('\n')
			return string(v), err
		default:
			v = append(v, b)
		}
	}

	return "", Malformed
}

func trim(v string) string {
	v = strings.Trim(v, " ")
	return strings.Trim(v, "\t")
}
