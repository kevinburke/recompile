package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	shellwords "github.com/mattn/go-shellwords"
	"golang.org/x/sync/errgroup"
)

var mu sync.Mutex

func main() {
	command := flag.String("command", "echo", "Command to run")
	extension := flag.String("extension", "", "Extension for new files (defaults to current file's extension)")
	outDir := flag.String("out-dir", "", "Output directory for new files (defaults to current file's directory)")
	flag.Parse()
	args := flag.Args()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmdArgs, err := shellwords.Parse(*command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing --command: %v\n", err)
		os.Exit(2)
	}
	group, errctx := errgroup.WithContext(ctx)
	buf := new(bytes.Buffer)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		group.Go(func() error {
			outFile := arg
			if *outDir != "" {
				_, file := filepath.Split(outFile)
				outFile = filepath.Join(*outDir, file)
			}
			if *extension != "" {
				ext := filepath.Ext(arg)
				newExt := *extension
				if newExt[0] != '.' {
					newExt = "." + newExt
				}
				outFile = outFile[0:len(outFile)-len(ext)] + newExt
			}
			// todo get fancy with printing out command duration?
			fmt.Printf("%s %q > %s\n", *command, arg, outFile)
			cmd := exec.CommandContext(errctx, cmdArgs[0], append(cmdArgs[1:], arg)...)
			// TODO write to temporary file and rename
			f, err := os.Create(outFile)
			if err != nil {
				return err
			}
			cmd.Stdout = f
			cmd.Stderr = buf
			return cmd.Run()
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		io.Copy(os.Stderr, buf)
		os.Exit(2)
	}
}
