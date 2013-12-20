package config

import (
	"os"
	"log"
	"bufio"
	"errors"
	"strings"
)

var (
	NotFound  = errors.New("not found")
	Malformed = errors.New("malformed config file")
)

type Config struct {
	content map[string]map[string]string
}

func (c *Config) Get(s, k string) (string, error) {
	if _, ok := c.content[s]; ok {
		if _, ok := c.content[s][k]; ok {
			return c.content[s][k], nil
		}
	}

	return "", NotFound
}

func New(path string) (*Config, error) {
	c := new(Config)
	c.content, err := parse(path)
	return c, err
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
			return content, err
		}

		switch b {
		case '\n': // skip
		case '#':
			_, err = r.ReadString('\n')
		case '[':
			// clear
			key = key[:0]
			value = value[:0]

			if section, err = readSection(r); err == nil {
				section = strings.Trim(section, " ")
			}
		default:
			if key, err = readKey(r); err == nil {
				key = strings.Trim(key, " ")
			}

			if value, err = readValue(r); err == nil {
				value = strings.Trim(value, " ")
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

func readdSection(r *bufio.Reader) (string, error) {
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
		case '#':
			_, err = r.ReadBytes('\n')
			return string(v), err
		default:
			v = append(v, b)
		}
	}

	return "", Malformed
}
