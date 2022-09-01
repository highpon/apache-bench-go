package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type apachBecnhArgs struct {
	n            int
	c            int
	t            int
	A            string
	P            string
	X            string
	url          string
	goRoutineNum int
}

var abargs apachBecnhArgs

func init() {
	flag.IntVar(&abargs.n, "n", 0, "requests")
	flag.IntVar(&abargs.c, "c", 0, "concurrency")
	flag.IntVar(&abargs.t, "t", 0, "timelimit")
	flag.StringVar(&abargs.A, "A", "", "Add Basic WWW Authentication, the attribute are a colon separated username and password.")
	flag.StringVar(&abargs.P, "P", "", "Add Basic Proxy Authentication, the attributes are a colon separated username and password.")
	flag.StringVar(&abargs.X, "X", "", "proxy:port   Proxyserver and port number to use")
	flag.StringVar(&abargs.url, "url", "", "apache bench url")
	flag.IntVar(&abargs.goRoutineNum, "goroutineNum", 0, "goroutine number")
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup

	wg.Add(abargs.goRoutineNum)
	for i := 0; i < abargs.goRoutineNum; i++ {
		abargs := abargs
		go func(abargs apachBecnhArgs) {
			out, err := runApachBench(abargs)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out)
			wg.Done()
		}(abargs)
	}
	wg.Wait()
}

func runApachBench(abargs apachBecnhArgs) (io.Writer, error) {
	cmdString := []string{}
	if abargs.n > 0 {
		cmdString = append(cmdString, "-n", strconv.Itoa(abargs.n))
	}

	if abargs.c > 0 {
		cmdString = append(cmdString, "-c", strconv.Itoa(abargs.c))
	}

	if abargs.t > 0 {
		cmdString = append(cmdString, "-t", strconv.Itoa(abargs.t))
	}

	if abargs.A != "" {
		cmdString = append(cmdString, "-A", abargs.A)
	}

	if abargs.P != "" {
		cmdString = append(cmdString, "-P", abargs.A)
	}

	if abargs.X != "" {
		cmdString = append(cmdString, "-X", abargs.A)
	}

	if abargs.url == "" {
		return nil, errors.New("url cannnot empty")
	}

	cmd := exec.Command("ab", strings.Join(cmdString, " "), abargs.url)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, "cmd.Run")
	}

	return cmd.Stdout, nil
}
