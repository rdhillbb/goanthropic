# Anthropic Client Public Functions Documentation

## Client Creation and Configuration

### NewClient
Creates and initializes a new AnthropicClient.
```go
func NewClient(apiKey string, opts ...ClientOption) *AnthropicClient
```
Parameters:
- `apiKey`: Your Anthropic API key
- `opts`: Optional configuration options

Example:
```go
client := NewClient("your-api-key",
    WithMaxConversationLength(10),
    WithDefaultParams(MessageParams{
        Model: "claude-3-5-sonnet-20241022",
        MaxTokens: 1000,
    }),
)
```

### Configuration Options

#### WithMaxConversationLength
Sets the maximum number of messages to retain in conversation history.
```go
func WithMaxConversationLength(length int) ClientOption
```

#### WithDefaultParams
Sets default parameters for all messages.
```go
func WithDefaultParams(params MessageParams) ClientOption
```

#### WithHTTPClient
Sets a custom HTTP client for API requests.
```go
func WithHTTPClient(client *http.Client) ClientOption
```

## Message Functions

### ChatMe
Handles a single message interaction while maintaining conversation history.
```go
func (c *AnthropicClient) ChatMe(
    ctx context.Context,
    message string,
    params *MessageParams,
) (*AnthropicResponse, error)
```
Parameters:
- `ctx`: Context for request cancellation
- `message`: User message text
- `params`: Message configuration parameters

Example:
```go
response, err := client.ChatMe(context.Background(), 
    "Tell me about quantum computing",
    &MessageParams{
        MaxTokens: 1000,
        Temperature: 0.7,
    },
)
```

### ChatWithTools
Implements tool interaction loop, allowing the assistant to use tools.
```go
func (c *AnthropicClient) ChatWithTools(
    ctx context.Context,
    message string,
    params *MessageParams,
    handlers map[string]func(context.Context, json.RawMessage) (string, error),
) (*AnthropicResponse, error)
```
Parameters:
- `ctx`: Context for request cancellation
- `message`: User message text
- `params`: Message configuration parameters
- `handlers`: Map of tool names to their handler functions

Example:
```go
handlers := map[string]func(context.Context, json.RawMessage) (string, error){
    "calculator": func(ctx context.Context, input json.RawMessage) (string, error) {
        // Calculator implementation
        return "4", nil
    },
}

response, err := client.ChatWithTools(context.Background(),
    "What is 2 + 2?",
    &MessageParams{
        Tools: []Tool{
            {
                Name: "calculator",
                Description: "Performs basic arithmetic",
                InputSchema: InputSchema{
                    Type: "object",
                    Properties: map[string]Property{
                        "expression": {
                            Type: "string",
                            Description: "Mathematical expression to evaluate",
                        },
                    },
                    Required: []string{"expression"},
                },
            },
        },
        ToolChoice: &ToolChoice{Type: "auto"},
    },
    handlers,
)
```

## Debug Logging Functions

### EnableDebug
Enables detailed debug logging with session tracking.
```go
func EnableDebug() error
```

### DisableDebug
Disables debug logging.
```go
func DisableDebug() error
```

### GetSessionID
Returns the current debug session identifier.
```go
func GetSessionID() string
```

### IsDebugEnabled
Returns the current debug logging state.
```go
func IsDebugEnabled() bool
```

## Message Logging Functions

### StartMessageLogging
Enables general message logging throughout the client.
```go
func StartMessageLogging() error
```

### StopMessageLogging
Disables message logging.
```go
func StopMessageLogging()
```

## Best Practices

1. **Error Handling**
   - Always check returned errors
   - Use context for request cancellation
   - Handle tool execution errors appropriately

2. **Resource Management**
   - Set appropriate conversation length limits
   - Configure suitable token limits
   - Use context timeouts for long-running operations

3. **Tool Implementation**
   - Keep tool handlers concise and focused
   - Validate tool inputs thoroughly
   - Handle tool errors gracefully
   - Document tool input schemas clearly

4. **Logging**
   - Enable debug logging during development
   - Use session IDs for tracking related requests
   - Disable logging in production if not needed
