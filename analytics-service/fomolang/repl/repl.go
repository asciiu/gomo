package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/asciiu/gomo/analytics-service/fomolang/evaluator"
	"github.com/asciiu/gomo/analytics-service/fomolang/lexer"
	"github.com/asciiu/gomo/analytics-service/fomolang/parser"
)

const PROMPT = ">> "

//func Start(in io.Reader, out io.Writer) {
//	io.WriteString(out, MONKEY_FACE)
//	scanner := bufio.NewScanner(in)
//
//	for {
//		fmt.Printf(PROMPT)
//		scanned := scanner.Scan()
//		if !scanned {
//			return
//		}
//
//		line := scanner.Text()
//		l := lexer.New(line)
//		p := parser.New(l)
//
//		program := p.ParseProgram()
//		if len(p.Errors()) != 0 {
//			printParserErrors(out, p.Errors())
//			continue
//		}
//
//		io.WriteString(out, program.String())
//		io.WriteString(out, "\n")
//	}
//}

func Start(in io.Reader, out io.Writer) {
	io.WriteString(out, MONKEY_FACE)
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
