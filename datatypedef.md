# Anthropic Client Data Types Documentation

## Core Types

### AnthropicClient
The primary client for interacting with the Anthropic API. Initialize using `NewClient()`.

### Message
Represents a single message in a conversation.
```go
type Message struct {
    Role    string           // "system", "user", or "assistant"
    Content []MessageContent // Array of content elements
}
```

### MessageContent
Defines different types of content within a message.
```go
type MessageContent struct {
    Type       string          // Content type identifier
    Text       string          // Plain text content
    ID         string          // Content identifier
    Name       string          // Tool name (for tool calls)
    Input      json.RawMessage // Tool input parameters
    ToolUseID  string          // Reference ID for tool results
    Content    string          // Tool result content
    IsError    bool            // Error indicator for tool results
}
```

### MessageParams
Configuration parameters for message requests.
```go
type MessageParams struct {
    Model       string                 // Model identifier (e.g., "claude-3-5-sonnet-20241022")
    MaxTokens   int                    // Maximum tokens in response
    Temperature float64                // Response randomness (0.0-1.0)
    TopP        float64                // Nucleus sampling parameter
    TopK        int                    // Top-k sampling parameter
    Metadata    map[string]interface{} // Optional request metadata
    System      string                 // System-level instructions
    Tools       []Tool                 // Available tools
    ToolChoice  *ToolChoice            // Tool selection preferences
}
```

## Tool-Related Types

### Tool
Defines a callable tool configuration.
```go
type Tool struct {
    Name         string      // Tool identifier
    Description  string      // Tool purpose description
    InputSchema  InputSchema // Expected input format
}
```

### InputSchema
Defines the structure of tool inputs.
```go
type InputSchema struct {
    Type       string              // Schema type (usually "object")
    Properties map[string]Property // Input field definitions
    Required   []string           // Required field names
}
```

### Property
Defines a single property in an input schema.
```go
type Property struct {
    Type        string   // Data type
    Description string   // Property description
    Enum        []string // Optional enumerated values
}
```

### ToolChoice
Configures tool selection preferences.
```go
type ToolChoice struct {
    Type string // "auto", "none", or "tool"
    Name string // Required when Type is "tool"
}
```

## Response Types

### AnthropicResponse
Contains the API response data.
```go
type AnthropicResponse struct {
    ID          string           // Response identifier
    Type        string           // Response type
    Role        string           // Responder role
    Content     []MessageContent // Response content array
    Model       string           // Model used
    StopReason  string           // Completion reason
    Usage       Usage           // Token usage statistics
}
```

### Usage
Tracks token usage in requests and responses.
```go
type Usage struct {
    InputTokens  int // Number of input tokens
    OutputTokens int // Number of output tokens
}
```

## Constants

### Role Constants
```go
const (
    RoleSystem    = "system"
    RoleUser      = "user"
    RoleAssistant = "assistant"
)
```

### Content Type Constants
```go
const (
    ContentTypeText       = "text"
    ContentTypeToolUse    = "tool_use"
    ContentTypeToolResult = "tool_result"
    ContentTypeThinking   = "thinking"
)
```

### Stop Reason Constants
```go
const (
    StopReasonToolUse      = "tool_use"
    StopReasonEndTurn      = "end_turn"
    StopReasonMaxTokens    = "max_tokens"
    StopReasonStopSequence = "stop_sequence"
)
```

### Tool Choice Constants
```go
const (
    ToolChoiceAuto = "auto"
    ToolChoiceNone = "none"
    ToolChoiceTool = "tool"
)
```
