package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    "github.com/joho/godotenv"
    
    "goanthropic"
    "goanthropic/tools"
)

const defaultModel = "claude-3-5-sonnet-20241022"

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get API key from environment variables
    apiKey := os.Getenv("ANTHROPIC_API_KEY")
    if apiKey == "" {
        log.Fatal("ANTHROPIC_API_KEY not found in environment variables")
    }

    // Only keep debug flag
    debug := flag.Bool("debug", false, "Enable debug logging")
    flag.Parse()

    if *debug {
        goanthropic.EnableDebug()
        defer goanthropic.DisableDebug()
    }

    client := goanthropic.NewClient(apiKey, 
        goanthropic.WithDefaultParams(anthropic.MessageParams{
            Model:      defaultModel,
            MaxTokens:  10000,
            Tools:      tools.GetDefaultTools(),
            ToolChoice: &goanthropic.ToolChoice{Type: anthropic.ToolChoiceAuto},
        }),
        goanthropic.WithMaxConversationLength(1000),
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
            &goanthropic.MessageParams{
                Model:      defaultModel,
                MaxTokens:  10000,
                Tools:      tools.GetDefaultTools(),
                ToolChoice: &goanthropic.ToolChoice{Type: goanthropic.ToolChoiceAuto},
            },
            handlers,
        )

        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        fmt.Println("\nAssistant:")
        for _, content := range response.Content {
            if content.Type == goanthropic.ContentTypeText {
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
