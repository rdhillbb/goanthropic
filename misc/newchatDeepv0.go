package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "os"
    "strings"
    
    "anthropicfunc/anthropic"
    "anthropicfunc/tools"
)

const defaultModel = "claude-3-5-sonnet-20241022"

func main() {
    apiKey := flag.String("api-key", "", "Anthropic API key")
    debug := flag.Bool("debug", false, "Enable debug logging")
    flag.Parse()

    if *apiKey == "" {
        fmt.Println("Error: API key is required")
        flag.Usage()
        os.Exit(1)
    }

    if *debug {
        anthropic.EnableDebug()
        defer anthropic.DisableDebug()
    }

    client := anthropic.NewClient(*apiKey, 
        anthropic.WithDefaultParams(anthropic.MessageParams{
            Model:      defaultModel,
            MaxTokens:  1000,
            Tools:      tools.GetDefaultTools(),
            ToolChoice: &anthropic.ToolChoice{Type: anthropic.ToolChoiceAuto},
        }),
        anthropic.WithMaxConversationLength(10),
    )

    handlers := tools.GetDefaultHandlers()
    scanner := bufio.NewScanner(os.Stdin)
    ctx := context.Background()

    fmt.Println("Chat initialized with tools. Type 'exit' to quit.")
    fmt.Println("Available tools:")
    for _, tool := range tools.GetDefaultTools() {
        fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
    }
    fmt.Println("\nEnter your message:")

    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }

        input := strings.TrimSpace(scanner.Text())
        if input == "exit" {
            break
        }

        if input == "" {
            continue
        }

        response, err := client.ChatWithTools(
            ctx,
            input,
            &anthropic.MessageParams{
                Model:      defaultModel,
                MaxTokens:  1000,
                Tools:      tools.GetDefaultTools(),
                ToolChoice: &anthropic.ToolChoice{Type: anthropic.ToolChoiceAuto},
            },
            handlers,
        )

        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        fmt.Println("\nAssistant:")
        for _, content := range response.Content {
            if content.Type == anthropic.ContentTypeText {
                fmt.Println(content.Text)
            }
        }
        fmt.Println()
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading input: %v\n", err)
        os.Exit(1)
    }
}
