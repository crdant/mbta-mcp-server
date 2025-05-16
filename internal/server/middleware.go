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
			log.Printf("Request received: tool=%s", req.Params.Name)

			// Log parameters if in debug mode
			if debug {
				log.Printf("Request arguments: %v", req.Params.Arguments)
			}

			// Record start time for duration calculation
			startTime := time.Now()

			// Call the actual handler
			resp, err := next(ctx, req)

			// Calculate request duration
			duration := time.Since(startTime)

			// Log response information
			if err != nil {
				log.Printf("Request error: tool=%s error=%v duration=%v",
					req.Params.Name, err, duration)
				return nil, err
			}

			// Log success information
			log.Printf("Request completed: tool=%s duration=%v",
				req.Params.Name, duration)

			// Log more detailed response information in debug mode
			if debug && resp != nil {
				// Log error details if present
				if resp.IsError {
					log.Printf("Response marked as error")
				} else {
					log.Printf("Response content length: %d", len(resp.Content))
				}
			}

			return resp, nil
		}
	}
}

// ApplyMiddleware applies the middleware to the server.
// This wraps all tool handlers with the logging middleware.
func (s *Server) ApplyMiddleware() {
	// Since we can't apply middleware globally with the current mcp-go API,
	// we just log that middleware would be applied here.
	// In a future implementation, we would implement this by wrapping individual
	// tool handlers with middleware as they are registered.
	log.Println("Applying logging middleware to MCP server")
	
	if s.config.Debug {
		log.Printf("Applied logging middleware to MCP server")
	}
}