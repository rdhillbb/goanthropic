package goanthropic

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "github.com/rdhillbb/goanthropic/types"
    "github.com/rdhillbb/logging"
)

const (
    defaultAPIEndpoint = "https://api.anthropic.com/v1/messages"
    defaultModel      = "claude-3-5-sonnet-20241022"
)

type ClientOption func(*AnthropicClient)

// AnthropicClient handles all communication with the Anthropic API
type AnthropicClient struct {
    apiKey          string
    defaultParams   types.MessageParams
    httpClient      *http.Client
    conversation    []types.Message
    maxConvLength   int
    systemPrompt    string
}

// NewClient creates a new AnthropicClient
func NewClient(apiKey string, opts ...ClientOption) *AnthropicClient {
    logMessage("Creating new AnthropicClient")
    client := &AnthropicClient{
        apiKey:     apiKey,
        httpClient: &http.Client{},
    }
    
    for _, opt := range opts {
        opt(client)
    }
    
    logJSON("Client configuration", map[string]interface{}{
        "maxConvLength": client.maxConvLength,
        "hasDefaults":   len(client.defaultParams.Tools) > 0 || 
                        client.defaultParams.MaxTokens > 0 ||
                        client.defaultParams.Model != "",
    })
    return client
}

// ChatWithTools handles chat interactions with tool support
func (c *AnthropicClient) ChatWithTools(ctx context.Context, message string, params *types.MessageParams, handlers []types.ToolHandler) (*types.AnthropicResponse, error) {
    // Use default params if none provided
    finalParams := c.defaultParams
    if params != nil {
        // Merge any non-zero params from the provided params
        if params.Model != "" {
            finalParams.Model = params.Model
        }
        if params.MaxTokens != 0 {
            finalParams.MaxTokens = params.MaxTokens
        }
        if params.Temperature != 0 {
            finalParams.Temperature = params.Temperature
        }
        if params.TopP != 0 {
            finalParams.TopP = params.TopP
        }
        if params.TopK != 0 {
            finalParams.TopK = params.TopK
        }
        if params.Tools != nil {
            finalParams.Tools = params.Tools
        }
        if params.ToolChoice != nil {
            finalParams.ToolChoice = params.ToolChoice
        }
    }

    // Validate the merged parameters
    if err := validateToolParams(&finalParams); err != nil {
        return nil, fmt.Errorf("invalid parameters: %w", err)
    }

    content := []types.MessageContent{{
        Type: types.ContentTypeText,
        Text: message,
    }}

    c.addMessageToConversation(types.RoleUser, content)
    c.trimConversationHistory()

    reqBody := types.Request{
        Model:       finalParams.Model,
        System:      c.systemPrompt,
        Messages:    c.conversation,
        MaxTokens:   finalParams.MaxTokens,
        Temperature: finalParams.Temperature,
        TopP:        finalParams.TopP,
        TopK:        finalParams.TopK,
        Tools:       finalParams.Tools,
        ToolChoice:  finalParams.ToolChoice,
    }

    response, err := c.sendRequest(ctx, reqBody)
    if err != nil {
        return nil, err
    }

    if len(response.Content) > 0 {
        c.addMessageToConversation(types.RoleAssistant, response.Content)
        c.trimConversationHistory()
    }

    return response, nil
}

// ChatMe handles basic chat interactions without tools
func (c *AnthropicClient) ChatMe(ctx context.Context, message string, params *types.MessageParams) (*types.AnthropicResponse, error) {
    finalParams := c.defaultParams
    if params != nil {
        if params.Model != "" {
            finalParams.Model = params.Model
        }
        if params.MaxTokens != 0 {
            finalParams.MaxTokens = params.MaxTokens
        }
        if params.Temperature != 0 {
            finalParams.Temperature = params.Temperature
        }
        if params.TopP != 0 {
            finalParams.TopP = params.TopP
        }
        if params.TopK != 0 {
            finalParams.TopK = params.TopK
        }
    }

    content := []types.MessageContent{{
        Type: types.ContentTypeText,
        Text: message,
    }}

    c.addMessageToConversation(types.RoleUser, content)
    c.trimConversationHistory()

    reqBody := types.Request{
        Model:       finalParams.Model,
        System:      c.systemPrompt,
        Messages:    c.conversation,
        MaxTokens:   finalParams.MaxTokens,
        Temperature: finalParams.Temperature,
        TopP:        finalParams.TopP,
        TopK:        finalParams.TopK,
    }

    response, err := c.sendRequest(ctx, reqBody)
    if err != nil {
        return nil, err
    }

    if len(response.Content) > 0 {
        c.addMessageToConversation(types.RoleAssistant, response.Content)
        c.trimConversationHistory()
    }

    return response, nil
}

// sendRequest handles the HTTP communication with the Anthropic API
func (c *AnthropicClient) sendRequest(ctx context.Context, reqBody types.Request) (*types.AnthropicResponse, error) {
    logMessage("Preparing API request")
    logJSON("Request payload", reqBody)

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        logMessage("Error marshaling request: %v", err)
        return nil, fmt.Errorf("error marshaling request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", defaultAPIEndpoint, bytes.NewBuffer(jsonData))
    if err != nil {
        logMessage("Error creating HTTP request: %v", err)
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("anthropic-version", "2023-06-01")
    req.Header.Set("x-api-key", c.apiKey)

    logMessage("Sending request to Anthropic API")
    resp, err := c.httpClient.Do(req)
    if err != nil {
        logMessage("API request failed: %v", err)
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logMessage("Error reading response body: %v", err)
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        logMessage("Received error response (status %d)", resp.StatusCode)
        var errorResp struct {
            Error struct {
                Type    string `json:"type"`
                Message string `json:"message"`
            } `json:"error"`
        }
        if err := json.Unmarshal(body, &errorResp); err != nil {
            logMessage("Failed to parse error response: %v", err)
            return nil, fmt.Errorf("error response status %d: %s", resp.StatusCode, body)
        }
        logMessage("API error: %s - %s", errorResp.Error.Type, errorResp.Error.Message)
        return nil, fmt.Errorf("API error: %s - %s", errorResp.Error.Type, errorResp.Error.Message)
    }

    var anthropicResp types.AnthropicResponse
    if err := json.Unmarshal(body, &anthropicResp); err != nil {
        logMessage("Error parsing response JSON: %v", err)
        return nil, fmt.Errorf("error parsing response: %w", err)
    }

    logJSON("API response", anthropicResp)
    return &anthropicResp, nil
}

// Conversation management methods
func (c *AnthropicClient) addMessageToConversation(role string, content []types.MessageContent) {
    logMessage("Adding message to conversation (role: %s)", role)
    c.conversation = append(c.conversation, types.Message{
        Role:    role,
        Content: content,
    })
}

func (c *AnthropicClient) trimConversationHistory() {
    if c.maxConvLength > 0 && len(c.conversation) > c.maxConvLength {
        logMessage("Trimming conversation to max length: %d", c.maxConvLength)
        c.conversation = c.conversation[len(c.conversation)-c.maxConvLength:]
    }
}

// Client options
func WithMaxConversationLength(length int) ClientOption {
    return func(c *AnthropicClient) {
        if length > 0 {
            c.maxConvLength = length
        }
    }
}

func WithDefaultParams(params types.MessageParams) ClientOption {
    return func(c *AnthropicClient) {
        c.defaultParams = params
    }
}

func WithHTTPClient(client *http.Client) ClientOption {
    return func(c *AnthropicClient) {
        if client != nil {
            c.httpClient = client
        }
    }
}

// Parameter validation
func validateToolParams(params *types.MessageParams) error {
    if params == nil {
        return fmt.Errorf("message parameters cannot be nil")
    }
    if params.Tools == nil {
        return fmt.Errorf("tools cannot be nil")
    }
    if params.ToolChoice == nil {
        return fmt.Errorf("tool choice cannot be nil")
    }
    return nil
}

// Logging helpers
func logMessage(format string, args ...interface{}) {
    if logging.IsLoggingEnabled() {
        msg := fmt.Sprintf(format, args...)
        logging.WriteLogs(fmt.Sprintf("[Anthropic] %s", msg))
    }
}

func logJSON(prefix string, data interface{}) {
    if logging.IsLoggingEnabled() {
        jsonBytes, err := json.MarshalIndent(data, "", "  ")
        if err != nil {
            logMessage("%s: failed to marshal JSON: %v", prefix, err)
            return
        }
        logging.WriteLogs(fmt.Sprintf("[Anthropic] %s: %s", prefix, string(jsonBytes)))
    }
}

// Package-level logging control
/*
func EnableDebug() error {
    return logging.EnableLogging()
}

func DisableDebug() {
    logging.DisableLogging()
}
*/
