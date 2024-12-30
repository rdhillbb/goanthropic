package goanthropic

import "encoding/json"

// Core API configuration constants
const (
    defaultAPIEndpoint = "https://api.anthropic.com/v1/messages"
    defaultModel      = "claude-3-5-sonnet-20241022"
)

// Role and content type constants
const (
    RoleSystem    = "system"
    RoleUser      = "user"
    RoleAssistant = "assistant"
    
    ContentTypeText       = "text"
    ContentTypeToolUse    = "tool_use"
    ContentTypeToolResult = "tool_result"
    ContentTypeThinking   = "thinking"  
    
    StopReasonToolUse      = "tool_use"
    StopReasonEndTurn      = "end_turn"
    StopReasonMaxTokens    = "max_tokens"
    StopReasonStopSequence = "stop_sequence"  
    
    ToolChoiceAuto = "auto"
    ToolChoiceNone = "none"
    ToolChoiceTool = "tool"
)

// ClientOption defines functions that can modify client configuration
type ClientOption func(*AnthropicClient)

// Message represents a single message in the conversation
type Message struct {
    Role    string           `json:"role"`    
    Content []MessageContent `json:"content"` 
}

// MessageContent represents different types of content within a message
type MessageContent struct {
    Type       string          `json:"type"`               
    Text       string          `json:"text,omitempty"`     
    ID         string          `json:"id,omitempty"`       
    Name       string          `json:"name,omitempty"`     
    Input      json.RawMessage `json:"input,omitempty"`    
    ToolUseID  string          `json:"tool_use_id,omitempty"`  
    Content    string          `json:"content,omitempty"`      
    IsError    bool            `json:"is_error,omitempty"`     
}

// ToolUse represents a tool call from the assistant
type ToolUse struct {
    ID    string          `json:"id"`
    Name  string          `json:"name"`
    Input json.RawMessage `json:"input"`
}

// MessageParams contains all possible parameters for a message request
type MessageParams struct {
    Model       string                 `json:"model"`
    MaxTokens   int                    `json:"max_tokens"`
    Temperature float64                `json:"temperature,omitempty"`
    TopP        float64                `json:"top_p,omitempty"`
    TopK        int                    `json:"top_k,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    System      string                 `json:"system,omitempty"`
    Tools       []Tool                 `json:"tools,omitempty"`
    ToolChoice  *ToolChoice            `json:"tool_choice,omitempty"`
}

// Request represents the complete structure sent to the Anthropic API
type Request struct {
    Model       string      `json:"model"`
    Messages    []Message   `json:"messages"`
    MaxTokens   int         `json:"max_tokens"`
    Temperature float64     `json:"temperature,omitempty"`
    TopP        float64     `json:"top_p,omitempty"`
    TopK        int         `json:"top_k,omitempty"`
    System      string      `json:"system,omitempty"`
    Tools       []Tool      `json:"tools,omitempty"`
    ToolChoice  *ToolChoice `json:"tool_choice,omitempty"`
}

// Tool-related types
type Tool struct {
    Name         string      `json:"name"`
    Description  string      `json:"description"`
    InputSchema  InputSchema `json:"input_schema"`
}

type InputSchema struct {
    Type       string              `json:"type"`
    Properties map[string]Property `json:"properties"`
    Required   []string           `json:"required"`
}

type Property struct {
    Type        string   `json:"type"`
    Description string   `json:"description"`
    Enum        []string `json:"enum,omitempty"`
}

type ToolChoice struct {
    Type string `json:"type"`
    Name string `json:"name,omitempty"`
}

// Response types
type AnthropicResponse struct {
    ID          string          `json:"id"`
    Type        string          `json:"type"`
    Role        string          `json:"role"`
    Content     []MessageContent `json:"content"`
    Model       string          `json:"model"`
    StopReason  string          `json:"stop_reason"`
    Usage       Usage           `json:"usage"`
}

type Usage struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}
