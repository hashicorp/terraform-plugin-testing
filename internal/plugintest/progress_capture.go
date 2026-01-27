// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugintest

import (
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/actioncheck"
)

// ProgressCapture handles capturing and storing action progress messages during test execution.
type ProgressCapture struct {
	messages map[string][]actioncheck.ProgressMessage
	mu       sync.RWMutex
}

// NewProgressCapture creates a new ProgressCapture instance.
func NewProgressCapture() *ProgressCapture {
	return &ProgressCapture{
		messages: make(map[string][]actioncheck.ProgressMessage),
	}
}

// CaptureProgress records a progress message for the specified action.
func (pc *ProgressCapture) CaptureProgress(actionName, message string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.messages[actionName] = append(pc.messages[actionName], actioncheck.ProgressMessage{
		Message:   message,
		Timestamp: time.Now(),
	})
}

// GetMessages returns all captured progress messages for the specified action.
func (pc *ProgressCapture) GetMessages(actionName string) []actioncheck.ProgressMessage {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	messages := pc.messages[actionName]
	if messages == nil {
		return []actioncheck.ProgressMessage{}
	}

	// Return a copy to prevent race conditions
	result := make([]actioncheck.ProgressMessage, len(messages))
	copy(result, messages)
	return result
}

// GetAllActionNames returns all action names that have captured messages.
func (pc *ProgressCapture) GetAllActionNames() []string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	var actionNames []string
	for actionName := range pc.messages {
		actionNames = append(actionNames, actionName)
	}
	return actionNames
}

// Clear removes all captured messages.
func (pc *ProgressCapture) Clear() {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.messages = make(map[string][]actioncheck.ProgressMessage)
}
