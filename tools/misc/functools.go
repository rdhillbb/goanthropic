package tools

import (
    "context"
    "encoding/json"
    "anthropicfunc/tavilysearch"
)

// Tool Definitions

// GetWeather returns the weather tool definition
func GetWeather() Tool {
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

// GetStock returns the stock price tool definition
func GetStock() Tool {
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

// GetSearch returns the internet search tool definition
func GetSearch() Tool {
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

// GetDeepSearch returns the deep search tool definition
func GetDeepSearch() Tool {
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

// HandleSearch processes internet search requests
func HandleSearch(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Query string `json:"query"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    
    results, err := tavilysearch.SearchInternet(params.Query)
    if err != nil {
        return "", err
    }
    return results, nil
}

// HandleDeepSearch processes comprehensive search requests
func HandleDeepSearch(ctx context.Context, args json.RawMessage) (string, error) {
    var params struct {
        Query string `json:"query"`
    }
    if err := json.Unmarshal(args, &params); err != nil {
        return "", err
    }
    
    results, err := tavilysearch.DeepSearch(params.Query)
    if err != nil {
        return "", err
    }
    return results, nil
}

// Tool and Handler Collections

// GetDefaultTools returns the default set of available tools
func GetDefaultTools() []Tool {
    return []Tool{
        GetWeather(),
        GetStock(),
        GetSearch(),
        GetDeepSearch(),
    }
}

// GetDefaultHandlers returns the default set of tool handlers
func GetDefaultHandlers() map[string]func(context.Context, json.RawMessage) (string, error) {
    return map[string]func(context.Context, json.RawMessage) (string, error){
        "get_weather":     HandleWeather,
        "get_stock_price": HandleStock,
        "SearchInternet":  HandleSearch,
        "DeepSearch":      HandleDeepSearch,
    }
}
