package tools

import (
    "context"
    "encoding/json"
    "errors"
     "anthropicfunc/anthropic"
)

var (
    ErrNilHandler     = errors.New("handler cannot be nil")
    ErrMissingHandler = errors.New("missing handler for tool")
)

type Tool struct {
    Name         string
    Description  string
    InputSchema  InputSchema
}

type InputSchema struct {
    Type       string
    Properties map[string]Property
    Required   []string
}

type Property struct {
    Type        string
    Description string
    Enum        []string
}

type ToolManager struct {
    tools    []anthropic.Tool
    handlers map[string]func(context.Context, json.RawMessage) (string, error)
}

func NewToolManager() *ToolManager {
    return &ToolManager{
        handlers: make(map[string]func(context.Context, json.RawMessage) (string, error)),
    }
}

func (tm *ToolManager) AddTool(tool anthropic.Tool, handler func(context.Context, json.RawMessage) (string, error)) error {
    if handler == nil {
        return ErrNilHandler
    }
    
    tm.tools = append(tm.tools, tool)
    tm.handlers[tool.Name] = handler
    return nil
}

func (tm *ToolManager) AddTools(tools []anthropic.Tool, handlers map[string]func(context.Context, json.RawMessage) (string, error)) error {
    for _, tool := range tools {
        handler, exists := handlers[tool.Name]
        if !exists {
            return ErrMissingHandler
        }
        if err := tm.AddTool(tool, handler); err != nil {
            return err
        }
    }
    return nil
}

func (tm *ToolManager) GetTools() []anthropic.Tool {
    return tm.tools
}

func (tm *ToolManager) GetHandlers() map[string]func(context.Context, json.RawMessage) (string, error) {
    return tm.handlers
}

func (tm *ToolManager) ValidateTool(toolName string) bool {
    for _, tool := range tm.tools {
        if tool.Name == toolName {
            _, hasHandler := tm.handlers[toolName]
            return hasHandler
        }
    }
    return false
}

func (t Tool) ToAnthropicTool() anthropic.Tool {
    return anthropic.Tool{
        Name:        t.Name,
        Description: t.Description,
        InputSchema: anthropic.InputSchema{
            Type:       t.InputSchema.Type,
            Properties: toAnthropicProperties(t.InputSchema.Properties),
            Required:   t.InputSchema.Required,
        },
    }
}

func toAnthropicProperties(props map[string]Property) map[string]anthropic.Property {
    anthropicProps := make(map[string]anthropic.Property)
    for key, prop := range props {
        anthropicProps[key] = anthropic.Property{
            Type:        prop.Type,
            Description: prop.Description,
            Enum:        prop.Enum,
        }
    }
    return anthropicProps
}

func (tm *ToolManager) GetAnthropicTools() []anthropic.Tool {
    anthropicTools := make([]anthropic.Tool, len(tm.tools))
    for i, tool := range tm.tools {
        anthropicTools[i] = tool.ToAnthropicTool()
    }
    return anthropicTools
}
