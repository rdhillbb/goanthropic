# Function Tools Documentation

## Overview
The Function Tools package provides a set of predefined tools and handlers for use with the Anthropic client. These tools enable capabilities like weather information retrieval, stock price lookups, and internet searches.

## Available Tools

### 1. Weather Tool
```go
func GetWeather() anthropic.Tool
```
Provides current weather information for a specified location.

**Input Schema:**
- `location` (required): City, country, or region (e.g., "San Francisco, CA")
- `unit` (optional): Temperature unit ("celsius" or "fahrenheit")

**Example Response:**
```json
{
    "temperature_c": 22,
    "temperature_f": 71,
    "condition": "sunny",
    "humidity": 65,
    "location": "San Francisco, CA"
}
```

### 2. Stock Price Tool
```go
func GetStock() anthropic.Tool
```
Retrieves current stock price information.

**Input Schema:**
- `symbol` (required): Stock symbol (e.g., "AAPL")

**Example Response:**
```json
{
    "symbol": "AAPL",
    "price": "150.00"
}
```

### 3. Internet Search Tool
```go
func GetSearch() anthropic.Tool
```
Performs basic internet searches using the Tavily search API.

**Input Schema:**
- `query` (required): Search query or question

### 4. Deep Search Tool
```go
func GetDeepSearch() anthropic.Tool
```
Performs comprehensive internet searches for detailed analysis.

**Input Schema:**
- `query` (required): Search query for detailed analysis

## Utility Functions

### GetDefaultTools
```go
func GetDefaultTools() []anthropic.Tool
```
Returns a slice containing all available tools. Use this when you want to enable all functionality.

### GetDefaultHandlers
```go
func GetDefaultHandlers() map[string]func(context.Context, json.RawMessage) (string, error)
```
Returns a map of all tool handlers. Use this in conjunction with `ChatWithTools`.

## Usage Example

```go
func main() {
    client := anthropic.NewClient("your-api-key")
    
    // Get all available tools and handlers
    tools := tools.GetDefaultTools()
    handlers := tools.GetDefaultHandlers()
    
    // Use tools in a chat session
    response, err := client.ChatWithTools(
        context.Background(),
        "What's the weather in London?",
        &anthropic.MessageParams{
            Tools: tools,
            ToolChoice: &anthropic.ToolChoice{
                Type: "auto",
            },
        },
        handlers,
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## Handler Implementation Details

Each tool has a corresponding handler function that processes the tool requests:

### HandleWeather
Processes weather information requests. Currently returns simulated data.

### HandleStock
Processes stock price requests. Currently returns simulated data.

### HandleSearch
Processes internet search requests using the Tavily search API.

### HandleDeepSearch
Processes comprehensive search requests using the Tavily search API's deep search capability.

## Best Practices

1. **Error Handling**
   - Always check for errors from handler functions
   - Validate input parameters before processing
   - Handle JSON unmarshaling errors appropriately

2. **Tool Selection**
   - Only include tools that are necessary for your use case
   - Consider creating custom tools for specific needs
   - Use the auto tool choice unless specific tool selection is required

3. **Search Usage**
   - Use regular search for basic information needs
   - Reserve deep search for complex queries requiring comprehensive analysis
   - Consider rate limits and costs associated with search API usage

4. **Data Validation**
   - Validate location formats for weather requests
   - Ensure stock symbols are properly formatted
   - Sanitize search queries before processing

## Note on Implementation
The current weather and stock handlers return simulated data. In a production environment, these should be connected to actual weather and stock market data APIs.

## Customization
You can create custom tools by following the same pattern:
1. Define the tool using the `anthropic.Tool` structure
2. Create a handler function that processes the tool's input
3. Add the tool and handler to your tools and handlers collections

Example of a custom tool:
```go
func GetCustomTool() anthropic.Tool {
    return anthropic.Tool{
        Name: "custom_tool",
        Description: "Description of your custom tool",
        InputSchema: anthropic.InputSchema{
            Type: "object",
            Properties: map[string]anthropic.Property{
                "parameter": {
                    Type: "string",
                    Description: "Parameter description",
                },
            },
            Required: []string{"parameter"},
        },
    }
}
```
