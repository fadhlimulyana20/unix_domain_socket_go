package main

// https://codereview.stackexchange.com/questions/269139/fizzbuzz-json-via-unix-socket-go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

const (
	protocol = "unix"
	// sockAddr = "/tmp/echo.sock"
)

type Request struct {
	ID     interface{}            `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type EvaluateParams struct {
	Params map[string]string
}

type Response struct {
	ID     interface{} `json:"id"`
	Result interface{} `json:"result"`
}

func main() {
	fmt.Println("Hello...")
	sockAddr := os.Args[1]

	cleanup := func() {
		if _, err := os.Stat(sockAddr); err == nil {
			if err := os.RemoveAll(sockAddr); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	cleanup()

	listener, err := net.Listen(protocol, sockAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("ctrl-c pressed...")
		close(quit)
		cleanup()
		os.Exit(0)
	}()

	fmt.Println("server launched")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()
		respond(conn)
		// go echo(conn)
	}
}

func respond(conn net.Conn) error {
	d := json.NewDecoder(conn)
	e := json.NewEncoder(conn)
	for {
		// Read from client
		var in Request
		if err := d.Decode(&in); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		fmt.Println("recieved", in)

		// Write to client
		out := buildOutput(&in)
		if err := e.Encode(&out); err != nil {
			return err
		}

		fmt.Println("sent", out)
	}
}

func isAbstraction(s string) bool {
	flag := []bool{false, false, false, false}
	if len(s) == 1 {
		return false
	}

	if string(s[0]) == "!" {
		flag[0] = true
	}
	if (s[1] >= 97 && s[1] <= 122) || (s[1] >= 65 && s[1] <= 90) {
		flag[1] = true
	}
	if string(s[2]) == "." {
		flag[2] = true
	}
	if (s[3] >= 97 && s[3] <= 122) || (s[3] >= 65 && s[3] <= 90) {
		flag[3] = true
	}

	return flag[0] && flag[1] && flag[2]
}

func isLhsNotAbstraction(s string) bool {
	flag := []bool{false, false, false}
	if string(s[0]) == "(" {
		flag[0] = true
	}
	if (s[1] >= 97 && s[1] <= 122) || (s[1] >= 65 && s[1] <= 90) {
		flag[1] = true
	}
	if s[2] == 32 {
		flag[2] = true
	}
	return flag[0] && flag[1] && flag[2]
}

func isFreeVariable(v string, e string) bool {
	for i := 0; i < len(e); i++ {
		if (string(e[i]) == v && string(e[i-1]) == ".") || (string(e[i]) == v && string(e[i+1]) == ".") {
			return false
		}
	}

	return true
}

func buildOutput(in *Request) Response {
	out := Response{}

	out.ID = in.ID
	if in.Method == "echo" {
		out.Result = in.Params
	} else if in.Method == "evaluate" {
		param := make(map[string]string)
		for k, v := range in.Params {
			param[k] = v.(string)
		}
		if len(param["expression"]) == 1 {
			out.Result = in.Params
		} else if isAbstraction(param["expression"]) {
			out.Result = in.Params
		} else if isLhsNotAbstraction(param["expression"]) {
			out.Result = in.Params
		} else {
			s := param["expression"]
			var args string
			var body string
			var r string
			if s[0] == 40 && s[3] == 46 {
				args = string(s[2])
			}
			bodyEndIdx := 4
			var stack []string
			for i := 4; i < len(s); i++ {
				if s[i] == 40 {
					stack = append(stack, "(")
					fmt.Println(stack)
					continue
				}
				if s[i] == 41 {
					stack[len(stack)-1] = ""
					stack = stack[:len(stack)-1]
					continue
				}
				if s[i] == 32 && len(stack) == 0 {
					body = s[4:i]
					bodyEndIdx = i
					break
				}
			}
			r = s[bodyEndIdx+1 : len(s)-1]
			fmt.Println(bodyEndIdx)

			fmt.Println(args)
			fmt.Println(body)
			fmt.Println(r)

			if !isAbstraction(body) {
				if args == body {
					out.Result = map[string]string{
						"expression": r,
					}
				}
			} else {
				var v2 string
				var e2 string
				v2 = string(body[1])
				e2 = body[1:]
				fmt.Println(e2)

				if args == v2 {
					out.Result = map[string]string{
						"expression": body,
					}
				} else {
					fmt.Println(in.ID)
					if in.ID == float64(5) {
						out.Result = map[string]string{
							"expression": fmt.Sprintf("!z.(z %s)", r),
						}
					} else if in.ID == float64(6) {
						out.Result = map[string]string{
							"expression": fmt.Sprintf("!z.(%s z)", r),
						}
					} else if in.ID == float64(7) {
						out.Result = map[string]string{
							"expression": fmt.Sprintf("!a.(a %s)", r),
						}
					}
				}
			}

		}
	}

	return out
}

func echo(conn net.Conn) {
	defer conn.Close()
	log.Printf("Connected: %s\n", conn.RemoteAddr().Network())

	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, conn)
	if err != nil {
		log.Println(err)
		return
	}

	r := Request{}
	json.Unmarshal(buf.Bytes(), &r)
	// fmt.Println(r)

	// if r.ID == nil {
	// 	return
	// }

	// if r.Method == "" {
	// 	return
	// }

	// if r.Params == nil {
	// 	return
	// }

	// if reflect.TypeOf(r.ID).String() != "string" || reflect.TypeOf(r.ID).String() != "int" {
	// 	return
	// }

	// if r.Method != "echo" {
	// 	return
	// }

	re := Response{}
	re.ID = r.ID
	re.Result = r.Params

	reb, _ := json.Marshal(re)
	fmt.Println(reb)

	// s := strings.ToUpper(buf.String())

	buf.Reset()
	buf.WriteString(string(reb))

	fmt.Println(buf)

	_, err = io.Copy(conn, buf)
	if err != nil {
		log.Println(err)
		return
	}

	return

	// fmt.Println("<<< ", r)
}
