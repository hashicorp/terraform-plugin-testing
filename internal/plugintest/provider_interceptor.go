// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugintest

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// NewProviderInterceptor creates a new provider interceptor for tfprotov5
func NewProviderInterceptor(provider tfprotov5.ProviderServer, capture *ProgressCapture) tfprotov5.ProviderServer {
	// Check if provider supports actions
	if actionProvider, ok := provider.(tfprotov5.ProviderServerWithActions); ok {
		return &ActionProviderInterceptorV5{
			ProviderServerWithActions: actionProvider,
			capture:                   capture,
		}
	}

	// Return regular provider if no actions
	return &ProviderInterceptorV5{
		ProviderServer: provider,
		capture:        capture,
	}
}

// NewProviderInterceptorV6 creates a new provider interceptor for tfprotov6
func NewProviderInterceptorV6(provider tfprotov6.ProviderServer, capture *ProgressCapture) tfprotov6.ProviderServer {
	// Check if provider supports actions
	if actionProvider, ok := provider.(tfprotov6.ProviderServerWithActions); ok {
		return &ActionProviderInterceptorV6{
			ProviderServerWithActions: actionProvider,
			capture:                   capture,
		}
	}

	// Return regular provider if no actions
	return &ProviderInterceptorV6{
		ProviderServer: provider,
		capture:        capture,
	}
}

// ProviderInterceptorV5 wraps a tfprotov5.ProviderServer
type ProviderInterceptorV5 struct {
	tfprotov5.ProviderServer
	capture *ProgressCapture
}

// ProviderInterceptorV6 wraps a tfprotov6.ProviderServer
type ProviderInterceptorV6 struct {
	tfprotov6.ProviderServer
	capture *ProgressCapture
}

// ActionProviderInterceptorV5 wraps a ProviderServerWithActions for tfprotov5
type ActionProviderInterceptorV5 struct {
	tfprotov5.ProviderServerWithActions
	capture *ProgressCapture
}

// ActionProviderInterceptorV6 wraps a ProviderServerWithActions for tfprotov6
type ActionProviderInterceptorV6 struct {
	tfprotov6.ProviderServerWithActions
	capture *ProgressCapture
}

// InvokeAction intercepts action invocations to capture progress messages (tfprotov5)
func (p *ActionProviderInterceptorV5) InvokeAction(ctx context.Context, req *tfprotov5.InvokeActionRequest) (*tfprotov5.InvokeActionServerStream, error) {
	stream, err := p.ProviderServerWithActions.InvokeAction(ctx, req)
	if err != nil {
		return stream, err
	}

	// Wrap the events iterator to capture progress messages
	originalEvents := stream.Events
	stream.Events = func(yield func(tfprotov5.InvokeActionEvent) bool) {
		for event := range originalEvents {
			// Capture progress messages
			if progress, ok := event.Type.(tfprotov5.ProgressInvokeActionEventType); ok {
				p.capture.CaptureProgress(req.ActionType, progress.Message)
			}

			// Continue with original event
			if !yield(event) {
				break
			}
		}
	}

	return stream, nil
}

// InvokeAction intercepts action invocations to capture progress messages (tfprotov6)
func (p *ActionProviderInterceptorV6) InvokeAction(ctx context.Context, req *tfprotov6.InvokeActionRequest) (*tfprotov6.InvokeActionServerStream, error) {
	stream, err := p.ProviderServerWithActions.InvokeAction(ctx, req)
	if err != nil {
		return stream, err
	}

	// Wrap the events iterator to capture progress messages
	originalEvents := stream.Events
	stream.Events = func(yield func(tfprotov6.InvokeActionEvent) bool) {
		for event := range originalEvents {
			// Capture progress messages
			if progress, ok := event.Type.(tfprotov6.ProgressInvokeActionEventType); ok {
				p.capture.CaptureProgress(req.ActionType, progress.Message)
			}

			// Continue with original event
			if !yield(event) {
				break
			}
		}
	}

	return stream, nil
}
