package main

import (
	"fmt"
	"os"
	"os/signal"
)

func isApplication(str string) bool {
	var stack []string
	delimeterIdx := 0
	for i := 0; i < len(str); i++ {
		if string(str[i]) == "(" {
			stack = append(stack, "(")
		} else if string(str[i]) == ")" {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 1 && string(str[i]) == " " {
			delimeterIdx = i
		}
	}
	fmt.Println(string(str[len(str)-1]))
	return string(str[0]) == "(" && string(str[len(str)-1]) == ")" && delimeterIdx != 0
}

func splitApplication(str string) (lhs string, rhs string) {
	var stack []string
	delimeterIdx := 0
	for i := 0; i < len(str); i++ {
		if string(str[i]) == "(" {
			stack = append(stack, "(")
		} else if string(str[i]) == ")" {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 1 && string(str[i]) == " " {
			delimeterIdx = i
		}
	}

	lhs = str[1:delimeterIdx]
	rhs = str[delimeterIdx+1 : len(str)-1]
	return lhs, rhs
}

func splitAbstraction(s string) (args string, body string) {
	foundExclamation := false
	foundDelimeter := false
	exclamationIdx := 0
	delimeterIdx := 0
	for i := 0; i < len(s); i++ {
		if foundDelimeter && foundExclamation {
			break
		}

		if string(s[i]) == "!" && !foundExclamation {
			foundExclamation = true
			exclamationIdx = i
			continue
		}

		if string(s[i]) == "." && !foundDelimeter {
			foundDelimeter = true
			delimeterIdx = i
		}
	}

	args = s[exclamationIdx:delimeterIdx]
	body = s[delimeterIdx+1:]
	return args, body
}

func isCharFree(c string, ex string) bool {
	found := false
	for i := 0; i < len(ex); i++ {
		if (string(ex[i]) == c && i > 2 && string(ex[i-2]) != ".") || (string(ex[i]) == c && string(ex[i+1]) != ".") {
			found = true
			return true
		}
	}

	return !found
}

func substitute(args string, body string, right string) string {
	if len(right) == 1 {
		if args == right {
			return body
		} else {
			return right
		}
	}

	return ""
}

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	// go func() {
	// 	<-quit
	// 	fmt.Println("ctrl-c pressed...")
	// 	close(quit)
	// 	os.Exit(0)
	// }()

	str := os.Args[1]

	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(isApplication(str))
	lhs, rhs := splitApplication(str)

	// fmt.Println("stdout: ", delimeterIdx)
	fmt.Println(lhs)
	fmt.Println(rhs)

	llhs, lrhs := splitApplication(lhs)
	fmt.Println(llhs)
	fmt.Println(lrhs)

	args, body := splitAbstraction(llhs)
	fmt.Println(args)
	fmt.Println(body)

	fmt.Println(isCharFree("b", "(!a.a a)"))

	// for {
	// 	// reader := bufio.NewReader(os.Stdin)
	// 	// str, err := reader.ReadString('\n')
	// 	str := os.Args[1]

	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }

	// 	fmt.Println(isApplication(str))
	// 	lhs, rhs := splitApplication(str)

	// 	// fmt.Println("stdout: ", delimeterIdx)
	// 	fmt.Println(lhs)
	// 	fmt.Println(rhs)
	// }
}
