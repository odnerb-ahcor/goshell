package goshell

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/chzyer/readline"
)

type params struct {
	value string
}

type Context struct {
	Rl     *readline.Instance
	Params map[string]*params
	RED    string
	GREN   string
	YELLOW string
	BLUE   string
	RESET  string
}

func newParams(value string) *params {
	return &params{
		value: value,
	}
}

func newContext(rl *readline.Instance) *Context {
	return &Context{
		Rl:     rl,
		RED:    RED,
		GREN:   "\033[92m",
		YELLOW: "\033[93m",
		BLUE:   "\033[94m",
		RESET:  "\033[0m",
		Params: make(map[string]*params),
	}
}

func (c *Context) PrintRed(format string) error {
	_, err := fmt.Printf("%s%s%s", RED, format, RESET)
	return err
}

func (c *Context) PrintGreen(format string) error {
	_, err := fmt.Printf("%s%s%s", GREN, format, RESET)
	return err
}

func (c *Context) PrintYellow(format string) error {
	_, err := fmt.Printf("%s%s%s", YELLOW, format, RESET)
	return err
}

func (c *Context) PrintBlue(format string) error {
	_, err := fmt.Printf("%s%s%s", BLUE, format, RESET)
	return err
}

func (c *Context) ScanInt(text string) int {
	return 0
}

func (p *params) ToString() string {
	return p.value
}

func (p *params) ToInt() (int, error) {
	value, err := strconv.Atoi(p.value)
	if err != nil {
		return value, errors.New("Invalid params " + p.value + ": expectative integer value")
	}
	return value, nil
}

func (p *params) ToDate() (time.Time, error) {
	dateFormat := "2006-01-02"

	return time.Parse(dateFormat, p.value)
}
