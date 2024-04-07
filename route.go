package goshell

import (
	"regexp"
	"strings"

	"github.com/chzyer/readline"
)

type router struct {
	Tree *tree
	Rl   *readline.Instance
}

func newRouter(rl *readline.Instance) *router {
	return &router{
		Rl:   rl,
		Tree: newTree(),
	}
}

func (r *router) AddRouter(router string, action func(c *Context), flag string) {
	prompt, params := r.splitParams(router)
	r.Tree.AddNode(strings.TrimSpace(prompt), action, flag, params)
}

func (r *router) GetAutoComplete() *readline.PrefixCompleter {
	complete := make([]readline.PrefixCompleterInterface, 0)
	r.autoComplete(r.Tree.Node, &complete)

	return readline.NewPrefixCompleter(complete...)
}

func (r *router) autoComplete(node *node, complete *[]readline.PrefixCompleterInterface) {
	if node != nil {
		r.autoComplete(node.Left, complete)
		subCom := make([]readline.PrefixCompleterInterface, 0)
		r.autoComplete(node.Tree.Node, &subCom)
		*complete = append(*complete, readline.PcItem(node.Prompt, subCom...))
		r.autoComplete(node.Right, complete)
	}
}

func (r *router) Execute(router string) {
	c := newContext(r.Rl)
	action, params := r.Tree.FindNode(router)

	if action == nil {
		c.PrintRed(router + " not found\n")
		return
	}

	if action.Data == nil {
		c.PrintYellow("usage: " + router + " <command>\n")
		return
	}

	if params != nil {
		for i, param := range action.Params {
			c.Params[param] = newParams(params[i])
		}
	}

	action.Data(c)
}

func (r *router) splitParams(prompt string) (string, []string) {
	pos := strings.Index(prompt, "[")

	if pos != -1 {
		teste := regexp.MustCompile(`\[([^\]]+)[^\s{0}]\]`)
		params := teste.FindAllString(prompt, -1)

		for i, param := range params {
			str := strings.ReplaceAll(param, "[", "")
			str = strings.ReplaceAll(str, "]", "")

			params[i] = str
		}

		line := strings.Split(prompt, "")
		prompt := strings.Join(line[:pos], "")
		return prompt, params
	}

	return prompt, nil
}
