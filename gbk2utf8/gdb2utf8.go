package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"os/exec"
	"bytes"
	"io/ioutil"
)

func main() {
	dir := "."
	pattern := "*"
	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			usage()
			return
		}
		pattern = os.Args[1]
	}

	files, err := LookupFiles(dir, pattern)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	
	allok := true
	for _, file := range files {
		err := ConvGBK2UTF8(file)
		if err != nil {
			fmt.Printf("convert [%v] from GBK to UTF8 failed. %v\n", err.Error())
			allok = false
		}
	}
	
	if allok {
		fmt.Printf("All convert OK\n")
	}
}

func ConvGBK2UTF8(file string) error {
	//TODO use go buildin function to iconv
	cmd := exec.Command("iconv", "-f", "gbk", "-t", "utf-8", file)
    cmd.Env = os.Environ()
    stdoutput := new(bytes.Buffer)
    erroutput := new(bytes.Buffer)
    cmd.Stdout = stdoutput
    cmd.Stderr = erroutput
    err := cmd.Run()
    if err != nil {
    	return errors.New(err.Error() + " " + erroutput.String())
    }
    
    bakcmd := exec.Command("mv", file, file + ".bak")
    err = bakcmd.Run()
    if err != nil {
    	return errors.New(err.Error() + " " + erroutput.String())
    }
    
    tty, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC, 0755)
	defer func() {
		tty.Close()
	}()
	
	ioutil.WriteFile(file, stdoutput.Bytes(), 0755)
	return nil
}

func LookupFiles(dir, pattern string) ([]string, error) {
	var files []string = make([]string, 0, 5)

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if ok, err := filepath.Match(pattern, f.Name()); err != nil {
			return err
		} else if ok {
			files = append(files, path)
		}
		return nil
	})

	if len(files) == 0 {
		return files, errors.New("Not found any files")
	}

	return files, err
}

func usage() {
	fmt.Printf("usage : %v <pattern>\n", os.Args[0])
	fmt.Printf("Example 1: %v *.cc\n", os.Args[0])
	fmt.Printf("Example 2: %v *\n", os.Args[0])
	fmt.Printf("Example 3: %v dir/*\n", os.Args[0])
}