// Package tools provides functionality for managing and executing various chat tools
package tools 

import (
    "context"
    "encoding/json"
    "errors"
)

// =====================================
// Type Definitions and Interfaces
// =====================================

// Tool represents a function that can be called by the AI

type Tool struct {
    Name         string
    Description  string
    InputSchema  InputSchema
}

// InputSchema defines the expected input format for a tool
type InputSchema struct {
    Type       string
    Properties map[string]Property
    Required   []string
}

// Property defines a single property in the input schema
type Property struct {
    Type        string
    Description string
    Enum        []string
}

// ToolManager handles the setup and configuration of all available tools
type ToolManager struct {
    tools    []Tool
    handlers map[string]func(context.Context, json.RawMessage) (string, error)
}


// =====================================
// Constructor and Initialization
// =====================================

// NewToolManager creates and initializes a new ToolManager
func NewToolManager() (*ToolManager, error) {
    tm := &ToolManager{
        handlers: make(map[string]func(context.Context, json.RawMessage) (string, error)),
    }
    
    if err := tm.setupTools(); err != nil {
        return nil, err
    }
    
    if err := tm.setupHandlers(); err != nil {
        return nil, err
    }
    
    return tm, nil
}

// =====================================
// Core Methods
// =====================================

// GetTools returns the list of available tools
func (tm *ToolManager) GetTools() []Tool {
    return tm.tools
}

// GetHandlers returns the map of tool handlers
func (tm *ToolManager) GetHandlers() map[string]func(context.Context, json.RawMessage) (string, error) {
    return tm.handlers
}

// ValidateTool checks if a tool exists and has a handler
func (tm *ToolManager) ValidateTool(toolName string) bool {
    for _, tool := range tm.tools {
        if tool.Name == toolName {
            _, hasHandler := tm.handlers[toolName]
            return hasHandler
        }
    }
    return false
}

// =====================================
// Tool Setup and Configuration
// =====================================

// setupTools initializes all available tools
func (tm *ToolManager) setupTools() error {
    tm.tools = []Tool{
        getWeatherTool(),
        getStockTool(),
        getSearchTool(),
        getDeepSearchTool(),
    }
    return nil
}

// Tool Definitions
func getWeatherTool() Tool {
    return Tool{
        Name: "get_weather",
        Description: "Get the current weather in a given location. Returns temperature, " +
            "conditions (sunny, cloudy, etc), and humidity. Always provide both Celsius " +
            "and Fahrenheit in your natural language response.",
        InputSchema: InputSchema{
            Type: "object",
            Properties: map[string]Property{
                "location": {
                    Type:        "string",
                    Description: "The location name (city, country, or region), e.g. 'San Francisco, CA' or 'Cambodia'",
                },
                "unit": {
                    Type:        "string",
                    Description: "Temperature unit (celsius or fahrenheit)",
                    Enum:        []string{"celsius", "fahrenheit"},
                },
            },
            Required: []string{"location"},
        },
    }
}

func getStockTool() Tool {
    return Tool{
        Name:        "get_stock_price",
        Description: "Get the current stock price for a given symbol",
        InputSchema: InputSchema{
            Type: "object",
            Properties: map[string]Property{
                "symbol": {
                    Type:        "string",
                    Description: "The stock symbol, e.g. AAPL",
                },
            },
            Required: []string{"symbol"},
        },
    }
}

func getSearchTool() Tool {
    return Tool{
        Name:        "SearchInternet",
        Description: "Search the internet for information when user requests it or when information is needed",
        InputSchema: InputSchema{
            Type: "object",
            Properties: map[string]Property{
                "query": {
                    Type:        "string",
                    Description: "The search query or question",
                },
            },
            Required: []string{"query"},
        },
    }
}

func getDeepSearchTool() Tool {
    return Tool{
        Name:        "DeepSearch",
        Description: "Perform a comprehensive search when deep analysis is requested",
        InputSchema: InputSchema{
            Type: "object",
            Properties: map[string]Property{
                "query": {
                    Type:        "string",
                    Description: "The search query or question for detailed analysis",
                },
            },
            Required: []string{"query"},
        },
    }
}

// =====================================
// Handler Setup and Implementation
// =====================================

// setupHandlers initializes all tool handlers
func (tm *ToolManager) setupHandlers() error {
    tm.handlers["get_weather"] = tm.handleWeather
    tm.handlers["get_stock_price"] = tm.handleStockPrice
    tm.handlers["SearchInternet"] = tm.handleSearch
    tm.handlers["DeepSearch"] = tm.handleDeepSearch
    return nil
}

// Handler Implementations
func (tm *ToolManager) handleWeather(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Location string `json:"location"`
        Unit     string `json:"unit"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    if params.Location == "" {
        return "", errors.New("location is required")
    }

    // Mock weather data - replace with actual API call
    tempC := 22
    tempF := (tempC * 9 / 5) + 32

    weather := map[string]interface{}{
        "temperature_c": tempC,
        "temperature_f": tempF,
        "condition":     "sunny",
        "humidity":      65,
        "location":      params.Location,
    }

    return json.Marshal(weather)
}

func (tm *ToolManager) handleStockPrice(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Symbol string `json:"symbol"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    if params.Symbol == "" {
        return "", errors.New("stock symbol is required")
    }

    // Mock stock data - replace with actual API call
    result := map[string]interface{}{
        "symbol":   params.Symbol,
        "price":    150.00,
        "currency": "USD",
    }

    return json.Marshal(result)
}

func (tm *ToolManager) handleSearch(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Query string `json:"query"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    if params.Query == "" {
        return "", errors.New("search query is required")
    }

    // Mock search result - replace with actual search implementation
    result := map[string]interface{}{
        "query": params.Query,
        "results": []string{
            "Sample search result 1",
            "Sample search result 2",
        },
    }

    return json.Marshal(result)
}

func (tm *ToolManager) handleDeepSearch(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Query string `json:"query"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    if params.Query == "" {
        return "", errors.New("deep search query is required")
    }

    // Mock deep search result - replace with actual implementation
    result := map[string]interface{}{
        "query": params.Query,
        "detailed_results": []string{
            "Comprehensive analysis result 1",
            "Comprehensive analysis result 2",
            "Related topics and insights",
        },
    }

    return json.Marshal(result)
}
