package main

import (
	"bufio"
	"fmt"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type MessageMetadata struct {
	Messages []struct {
		Name   string
		Type   int
		Fields []struct {
			Name     string
			JsonName string `yaml:"json-name"`
			Type     string
		}
		Contexts []string
	}
}

func main() {

	messageMetadata := MessageMetadata{}

	data, _ := ioutil.ReadFile("./generate/messages.yaml")
	err := yaml.Unmarshal(data, &messageMetadata)
	if err != nil {
		panic(err)
	}

	GenerateMessageTypes(messageMetadata, os.Args[1])
}

func GenerateMessageTypes(messageMetadata MessageMetadata, path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	_, _ = fmt.Fprintln(w, "################################################################################")
	_, _ = fmt.Fprintln(w, "# This code is generated by `go run ./generate/godotmessages`. DON'T TOUCH IT! #")
	_, _ = fmt.Fprintln(w, "################################################################################")
	_, _ = fmt.Fprintf(w, "\n")

	_, _ = fmt.Fprintln(w, "extends Node\n")
	for _, message := range messageMetadata.Messages {
		_, _ = fmt.Fprintf(w, "const %s = %d\n",
			strings.ToUpper(strcase.ToSnake(message.Name)), message.Type)
	}
	_ = w.Flush()
}
