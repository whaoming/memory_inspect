package main

import (
	"os/exec"
	"fmt"
	"bufio"
	"io"
	"strings"
	"regexp"
)

type Process struct {
	pid int
	cpu float64
}

var values map[string]map[string]string
var keys []string
var isLsatEmptyLine bool

func main() {
	cmd := exec.Command("top")
	isLsatEmptyLine = false
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	i := 0
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		ParseCommandLine(line)
		i++
	}

}

func ParseCommandLine(line string) {
	if strings.Compare(line, "\n") == 0 {
		values = make(map[string]map[string]string)
		var tmp []string
		keys = tmp
		isLsatEmptyLine = true
		return
	}
	if values == nil {
		return
	}
	words := strings.Fields(line)
	if isLsatEmptyLine {
		i := 0
		for {
			if i > len(words)-1 {
				break
			}
			keys = append(keys, words[i])
			i++
		}
		isLsatEmptyLine = false
		return
	}
	regex := "[0x30,0x39]"
	isProcess, _ := regexp.MatchString(regex, words[0])
	if isProcess {
		i := 0
		pid := words[0]
		if values[pid] == nil {
			values[pid] = make(map[string]string)
		}
		for {
			if i > len(words)-1 {
				break
			}
			if i < len(keys){
				key := keys[i]
				values[pid][key] = words[i]
			}
			i++
		}
	}
	isLsatEmptyLine = false
}
