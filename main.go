package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

type Server struct {
	Location string
}

type Servers struct {
	Servers []Server
}

func main() {

	// Register command-line flags.
	port := flag.String("port", "800", "port to broadcast server on")
	jsonFile := flag.String("p", "", "file to parse")
	flag.Parse()

	if *jsonFile == "" {
		log.Fatal(fmt.Errorf("You need a file to parse, use -p"))
	}
	data, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var servers Servers
	err = json.Unmarshal(data, &servers)
	if err != nil {
		log.Fatal(err)
	}

	packageDec := "package main\n"

	imports := importBlock{
		imports: []string{
			"net/http",
			"fmt",
			"html",
		},
	}

	main := mainBlock{
		lines: []string{
			`fmt.Println("listening on port ` + *port + `")`,
		},
	}

	var handlerFunctions []handlerBlock
	for _, value := range servers.Servers {
		handlerName := formatFunctionName(value.Location) + "handler"
		main.lines = append(
			main.lines,
			`http.HandleFunc("`+value.Location+`", `+handlerName+`)`,
		)
		handler := handlerBlock{
			name: handlerName,
		}
		handlerFunctions = append(handlerFunctions, handler)
	}

	main.lines = append(main.lines, "http.ListenAndServe(\":"+*port+"\", nil)")

	code := packageDec + imports.codeGen() + main.codeGen()
	for _, value := range handlerFunctions {
		code = code + value.codeGen()
	}
	// module.codeGen()

	err = buildAndRun(code)
	if err != nil {
		fmt.Println(code)
		log.Fatal(err)
	}

}

func formatFunctionName(url string) string {
	validUrlCharachters := []string{"0", "1", "2", "3", "4",
		"5", "6", "7", "8", "9", "-", ".", "_", "~", ":",
		"/", "?", "#", "[", "]", "@", "!", "$", "&", "'",
		"(", ")", "*", "+", ",", ";", "=", "/"}
	for _, value := range validUrlCharachters {
		url = strings.Replace(url, value, "", -1)
	}
	return url
}

func buildAndRun(code string) (err error) {
	err = ioutil.WriteFile("gencode.go", []byte(code), 0777)

	output, err := runCommand("gofmt", "gencode.go")
	if err != nil {
		return err
	}
	color.Green("Generated Code:")
	fmt.Println(string(output))
	err = ioutil.WriteFile("gencode.go", output, 0777)
	if err != nil {
		return err
	}

	output, err = runCommand("go", "build", "gencode.go")
	if err != nil {
		return err
	}

	output, err = runCommand("rm", "gencode.go")
	if err != nil {
		return err
	}

	color.Green("Running Binary:")
	cmd := exec.Command("./gencode")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("rm", "gencode")
	err = cmd.Run()
	if err != nil {
		return err
	}
	return
}

func runCommand(name string, arg ...string) (output []byte, err error) {
	cmd := exec.Command(name, arg...)
	var o bytes.Buffer
	var e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e
	err = cmd.Run()
	if err != nil {
		return output, fmt.Errorf(string(e.Bytes()), err)
	}
	return o.Bytes(), err

}

type mainBlock struct {
	lines []string
}

func (i *mainBlock) codeGen() string {
	return "func main() {\n" +
		strings.Join(i.lines, "\n") +
		"\n}\n\n"
}

type importBlock struct {
	imports []string
}

func (i *importBlock) codeGen() string {
	return "import ( \n\"" +
		strings.Join(i.imports, "\"\n\"") +
		"\"\n)\n"
}

type handlerBlock struct {
	name  string
	lines []string
}

func (h *handlerBlock) codeGen() string {
	return "func " + h.name +
		"(w http.ResponseWriter, r *http.Request) {\n" +
		`fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))` +
		"}\n\n"

}
