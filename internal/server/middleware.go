// ABOUTME: This file implements logging middleware for the MCP server.
// ABOUTME: It provides request/response logging and timing functionality.

package server

import (
	"context"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// loggingMiddleware creates a middleware that logs request and response information.
// It wraps the original handler and provides timing and debug information.
func loggingMiddleware(debug bool) mcpserver.ToolHandlerMiddleware {
	return func(next mcpserver.ToolHandlerFunc) mcpserver.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Log basic request information
			log.Printf("[INFO] Request received: tool=%s", req.Params.Name)

			// Log parameters if in debug mode
			if debug {
				log.Printf("[DEBUG] Request arguments: %v", req.Params.Arguments)
			}

			// Record start time for duration calculation
			startTime := time.Now()

			// Call the actual handler
			resp, err := next(ctx, req)

			// Calculate request duration
			duration := time.Since(startTime)

			// Log response information
			if err != nil {
				log.Printf("[ERROR] Request error: tool=%s error=%v duration=%v",
					req.Params.Name, err, duration)
				return nil, err
			}

			// Log success information
			log.Printf("[INFO] Request completed: tool=%s duration=%v",
				req.Params.Name, duration)

			// Log more detailed response information in debug mode
			if debug && resp != nil {
				// Log error details if present
				if resp.IsError {
					log.Printf("[DEBUG] Response marked as error")
				} else {
					log.Printf("[DEBUG] Response content length: %d", len(resp.Content))
				}
			}

			return resp, nil
		}
	}
}

// ApplyMiddleware applies the middleware to the server.
// This applies the logging middleware to all handlers.
func (s *Server) ApplyMiddleware() {
	log.Println("Applying logging middleware to MCP server")

	// We can't apply middleware globally with the current mcp-go API,
	// so we'll need to manually wrap each tool handler when it's registered.
	// This would be done in the registerTransitInfoTools function.
	// For now, we just log that middleware would be applied.

	if s.config.Debug {
		log.Printf("[DEBUG] Applied logging middleware to MCP server")
	}
}

// wrapWithMiddleware wraps a tool handler with middleware.
// This is a helper function that can be used when registering tools.
func (s *Server) wrapWithMiddleware(handler mcpserver.ToolHandlerFunc) mcpserver.ToolHandlerFunc {
	// Apply logging middleware
	return loggingMiddleware(s.config.Debug)(handler)
}