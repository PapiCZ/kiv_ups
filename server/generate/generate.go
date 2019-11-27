package main

import (
	"bufio"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type MessageMetadata struct {
	ContextOnConnect string `yaml:"context-on-connect"`
	Messages         []struct {
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

	GenerateMessageTypes(messageMetadata)
	GenerateActions(messageMetadata)
	GenerateGodotMessageTypes(messageMetadata)
}

func GenerateMessageTypes(messageMetadata MessageMetadata) {
	f := NewFile("protocol")
	f.Comment("##############################################################")
	f.Comment("This code is generated by `go run ./generate`. DON'T TOUCH IT!")
	f.Comment("##############################################################")
	f.Empty()
	for _, message := range messageMetadata.Messages {
		fields := make([]Code, 0)

		for _, field := range message.Fields {
			jsonName := field.JsonName

			if len(jsonName) == 0 {
				jsonName = strcase.ToSnake(field.Name)
			}

			fields = append(
				fields,
				Id(field.Name).Id(field.Type).Tag(map[string]string{"json": jsonName}),
			)
		}

		f.Type().Id(message.Name + "Message").Struct(fields...)
		f.Empty()
		f.Func().Params(Id("m").Id(message.Name + "Message")).Id("GetTypeId").
			Params().Id("MessageType").Block(Return(Lit(message.Type)))
	}
	f.Empty()
	registerStatements := make([]Code, 0)
	for _, messageType := range messageMetadata.Messages {
		registerStatements = append(
			registerStatements,
			Id("definition.Register").Call(Id(messageType.Name+"Message{}")),
		)
	}
	f.Func().Id("RegisterAllMessages").Params(Id("definition").Id("*Definition")).Block(registerStatements...)

	file, err := os.OpenFile("./net/tcp/protocol/messagetypes.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	_, _ = fmt.Fprintf(w, "%#v", f)
	_ = w.Flush()
}

func GenerateActions(messageMetadata MessageMetadata) {
	f := NewFile("actions")
	f.Comment("##############################################################")
	f.Comment("This code is generated by `go run ./generate`. DON'T TOUCH IT!")
	f.Comment("##############################################################")
	f.Empty()

	allContextsMap := make(map[string]struct{}, 0)

	for _, message := range messageMetadata.Messages {
		for _, ctx := range message.Contexts {
			allContextsMap[ctx] = struct{}{}
		}
	}

	allContexts := make([]string, 0)
	allContexts = append(allContexts, messageMetadata.ContextOnConnect)

	for ctx := range allContextsMap {
		if ctx != "ALL" && ctx != messageMetadata.ContextOnConnect {
			allContexts = append(allContexts, ctx)
		}
	}

	allContextStatements := make([]Code, 0)
	for i, ctx := range allContexts {
		allContextStatements = append(allContextStatements,
			Id(ctx).Op("=").Qual("kiv_ups_server/game/interfaces", "PlayerContext").Params(Lit(i)))
	}

	f.Const().Defs(allContextStatements...)


	for _, message := range messageMetadata.Messages {
		f.Type().Id(message.Name + "Action").Struct()

		f.Empty()

		var contexts []string

		if len(message.Contexts) == 1 && message.Contexts[0] == "ALL" {
			contexts = allContexts
		} else {
			contexts = message.Contexts
		}

		contextStatements := make([]Code, 0)

		for _, ctx := range contexts {
			contextStatements = append(contextStatements, Id(ctx))
		}

		f.Func().Params(Id("a").Id(message.Name+"Action")).Id("GetPlayerContexts").
			Params().Index().Qual("kiv_ups_server/game/interfaces", "PlayerContext").
			Block(Return(Index().Qual("kiv_ups_server/game/interfaces", "PlayerContext").Values(contextStatements...)))
		f.Empty()
		f.Func().Params(Id("a").Id(message.Name+"Action")).Id("GetMessage").
			Params().Qual("kiv_ups_server/net/tcp/protocol", "Message").
			Block(Return(Qual("kiv_ups_server/net/tcp/protocol", message.Name+"Message{}")))
	}
	f.Empty()
	registerStatements := make([]Code, 0)
	for _, messageType := range messageMetadata.Messages {
		registerStatements = append(
			registerStatements,
			Id("actionDefinition.Register").Call(Id(messageType.Name+"Action{}")),
		)
	}
	f.Func().Id("RegisterAllActions").Params(Id("actionDefinition").Id("*ActionDefinition")).Block(registerStatements...)

	file, err := os.OpenFile("./game/actions/defs.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	_, _ = fmt.Fprintf(w, "%#v", f)
	_ = w.Flush()
}

func GenerateGodotMessageTypes(messageMetadata MessageMetadata) {
	file, err := os.OpenFile("../client/networking/MessageTypes.gd", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	_, _ = fmt.Fprintln(w, "##################################################################")
	_, _ = fmt.Fprintln(w, "# This code is generated by `go run ./generate`. DON'T TOUCH IT! #")
	_, _ = fmt.Fprintln(w, "##################################################################")
	_, _ = fmt.Fprintf(w, "\n")

	_, _ = fmt.Fprintln(w, "extends Node\n")
	for _, message := range messageMetadata.Messages {
		_, _ = fmt.Fprintf(w, "const %s = %d\n",
			strings.ToUpper(strcase.ToSnake(message.Name)), message.Type)
	}
	_ = w.Flush()
}
