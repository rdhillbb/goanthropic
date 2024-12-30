package main 

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/rdhillbb/goanthropic/types"
)

//
// GetWeather returns the weather tool definition using Anthropic types
func GetWeather() types.Tool {
    return types.Tool{
        Name: "get_weather",
        Description: "Get the current weather in a given location. Returns temperature, " +
            "conditions (sunny, cloudy, etc), and humidity. Always provide both Celsius " +
            "and Fahrenheit in your natural language response.",
        InputSchema: types.InputSchema{
            Type: "object",
            Properties: map[string]types.Property{
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

// GetStock returns the stock price tool definition using Anthropic types
func GetStock() types.Tool {
    return types.Tool{
        Name:        "get_stock_price",
        Description: "Get the current stock price for a given symbol",
        InputSchema: types.InputSchema{
            Type: "object",
            Properties: map[string]types.Property{
                "symbol": {
                    Type:        "string",
                    Description: "The stock symbol, e.g. AAPL",
                },
            },
            Required: []string{"symbol"},
        },
    }
}

// GetSearch returns the internet search tool definition using Anthropic types
func GetSearch() types.Tool {
    return types.Tool{
        Name:        "search",
        Description: "Search for information (placeholder - implementation needed)",
        InputSchema: types.InputSchema{
            Type: "object",
            Properties: map[string]types.Property{
                "query": {
                    Type:        "string",
                    Description: "The search query",
                },
            },
            Required: []string{"query"},
        },
    }
}

// Handler Implementations

// HandleWeather processes weather information requests
func HandleWeather(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Location string `json:"location"`
        Unit     string `json:"unit"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }

    // Placeholder implementation
    tempC := 22
    tempF := (tempC * 9 / 5) + 32
    weather := map[string]interface{}{
        "temperature_c": tempC,
        "temperature_f": tempF,
        "condition":    "sunny",
        "humidity":     65,
        "location":     params.Location,
    }
    
    jsonBytes, err := json.Marshal(weather)
    if err != nil {
        return "", err
    }
    return string(jsonBytes), nil
}

// HandleStock processes stock price requests
func HandleStock(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Symbol string `json:"symbol"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    
    // Placeholder implementation
    result := map[string]string{
        "symbol": params.Symbol,
        "price":  "150.00",
    }
    
    jsonBytes, err := json.Marshal(result)
    if err != nil {
        return "", err
    }
    return string(jsonBytes), nil
}

// HandleSearch processes search requests
func HandleSearch(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Query string `json:"query"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    
    // Placeholder implementation
    return fmt.Sprintf("Search results for: %s (implementation needed)", params.Query), nil
}

// GetDefaultTools returns the default set of available tools
func GetDefaultTools() []types.Tool {
    return []types.Tool{
        GetWeather(),
        GetStock(),
        GetSearch(),
    }
}

// GetDefaultHandlers returns the default set of tool handlers
func GetDefaultHandlers() map[string]func(context.Context, json.RawMessage) (string, error) {
    return map[string]func(context.Context, json.RawMessage) (string, error){
        "get_weather":     HandleWeather,
        "get_stock_price": HandleStock,
        "search":          HandleSearch,
    }
}
