package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
)

func getFileSizeString(fileInfo os.FileInfo) string {
	if fileInfo.Size() == 0 {
		return "(empty)"
	}

	return "(" + strconv.FormatInt(fileInfo.Size(), 10) + "b)"
}

func dirTree(out *bytes.Buffer, path string, printFiles bool) error {
	return dirTreeHelper(out, path, printFiles, 0, new(bytes.Buffer))
}

func getNewPrefixBuffer(index, totalCount int, prefixBuffer *bytes.Buffer) (*bytes.Buffer, error) {
	newPrefixBuffer := new(bytes.Buffer)
	_, err := newPrefixBuffer.Write(prefixBuffer.Bytes())
	if err != nil {
		return nil, err
	}
	if index+1 == totalCount {
		_, err = newPrefixBuffer.WriteString("\t")
	} else {
		_, err = newPrefixBuffer.WriteString("│\t")
	}
	if err != nil {
		return nil, err
	}
	return newPrefixBuffer, nil
}

func filter(files []os.FileInfo, predicate func(os.FileInfo) bool) []os.FileInfo {
	result := make([]os.FileInfo, 0, len(files)/10)
	for _, file := range files {
		if predicate(file) {
			result = append(result, file)
		}
	}
	return result
}

func dirTreeHelper(
	out *bytes.Buffer,
	path string,
	printFiles bool,
	depth int,
	prefixBuffer *bytes.Buffer,
) (err error) {
	defer func() {
		if err != nil {
			out.Reset()
		}
	}()

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	files = filter(files, func(file os.FileInfo) bool {
		return file.IsDir() || printFiles
	})

	for index, file := range files {
		marker := "├───"
		if index+1 == len(files) {
			marker = "└───"
		}
		if _, err = out.Write(prefixBuffer.Bytes()); err != nil {
			return
		}
		if _, err = out.WriteString(marker + file.Name()); err != nil {
			return
		}
		if file.IsDir() {
			newPrefixBuffer, err := getNewPrefixBuffer(index, len(files), prefixBuffer)
			if err != nil {
				return err
			}
			if _, err = out.WriteString("\n"); err != nil {
				return err
			}
			if err = dirTreeHelper(out, path+"/"+file.Name(), printFiles, depth+1, newPrefixBuffer); err != nil {
				return err
			}
		} else if printFiles {
			if _, err = out.WriteString(" " + getFileSizeString(file) + "\n"); err != nil {
				return
			}
		}
	}

	return nil
}

func main() {
	buffer := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(buffer, path, printFiles)
	if err != nil {
		panic(err.Error())
	} else {
		os.Stdout.Write(buffer.Bytes())
	}
}
