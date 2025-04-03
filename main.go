package main

import (
	"fmt"

	"github.com/invopop/jsonschema"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

// Tool arguments are just structs, annotated with jsonschema tags
// More at https://mcpgolang.com/tools#schema-generation
type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}

type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai, google, claude, etc)"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

type CalculateArguments struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	jsonschema.Version = "https://json-schema.org/draft-07/schema"
	// Removed setting jsonschema.Version as Draft202012Version is not declared in the package

	err := server.RegisterTool(
		"hello",
		"Say hello to a person",
		func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s!", arguments.Submitter))), nil
		})
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool(
		"calculate",
		"Calculate two numbers",
		func(arguments CalculateArguments) (*mcp_golang.ToolResponse, error) {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("A + B is %d", arguments.A+arguments.B))), nil
		})
	if err != nil {
		panic(err)
	}

	err = server.RegisterPrompt("promt_test", "This is a test prompt", func(arguments Content) (*mcp_golang.PromptResponse, error) {
		return mcp_golang.NewPromptResponse("description", mcp_golang.NewPromptMessage(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s!", arguments.Title)), mcp_golang.RoleUser)), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterResource("test://resource", "resource_test", "This is a test resource", "application/json", func() (*mcp_golang.ResourceResponse, error) {
		return mcp_golang.NewResourceResponse(mcp_golang.NewTextEmbeddedResource("test://resource", "This is a test resource", "application/json")), nil
	})

	err = server.Serve()

	fmt.Println("Server started, waiting for requests...")

	if err != nil {
		panic(err)
	}

	<-done
}
