package goshell

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

type Shell struct {
	Rl     *readline.Instance
	router *router
	active bool
}

func New(prompt string) *Shell {
	rl, _ := readline.NewEx(&readline.Config{
		Prompt: prompt + GREN + "» " + RESET,
	})

	shell := &Shell{
		Rl:     rl,
		router: newRouter(rl),
		active: true,
	}
	shell.addExit()
	return shell
}

func (s *Shell) Close() {
	s.Rl.Close()
}

func (s *Shell) Start() {
	s.Rl.Config.AutoComplete = s.router.GetAutoComplete()
	count := 1

	for s.active {
		print("\n")
		opcao, err := s.Rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if count >= 2 {
					os.Exit(1)
				}
				count++
				fmt.Printf("%s%s%s\n", YELLOW, "Input Ctrl-c once more to exit", RESET)
			}
			continue
		}
		count = 1
		opcao = strings.TrimSpace(opcao)
		s.router.Execute(opcao)
	}
}

func (s *Shell) addExit() {
	s.Router("exit", func(ct *Context) {
		s.active = false
		s.Rl.Close()
	}, "exit the program")
}

func (s *Shell) Router(router string, action func(c *Context), flag string) {
	s.router.AddRouter(router, action, flag)
}
