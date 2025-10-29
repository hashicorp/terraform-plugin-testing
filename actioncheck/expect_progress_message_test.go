// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actioncheck_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/actioncheck"
)

func TestExpectProgressMessageContains_Success(t *testing.T) {
	check := actioncheck.ExpectProgressMessageContains("test_action", "expected content")
	
	req := actioncheck.CheckActionRequest{
		ActionName: "test_action",
		Messages: []actioncheck.ProgressMessage{
			{Message: "some expected content here", Timestamp: time.Now()},
		},
	}
	
	resp := &actioncheck.CheckActionResponse{}
	check.CheckAction(context.Background(), req, resp)
	
	if resp.Error != nil {
		t.Errorf("Expected no error, got: %v", resp.Error)
	}
}

func TestExpectProgressMessageContains_Failure(t *testing.T) {
	check := actioncheck.ExpectProgressMessageContains("test_action", "missing content")
	
	req := actioncheck.CheckActionRequest{
		ActionName: "test_action",
		Messages: []actioncheck.ProgressMessage{
			{Message: "some other content", Timestamp: time.Now()},
		},
	}
	
	resp := &actioncheck.CheckActionResponse{}
	check.CheckAction(context.Background(), req, resp)
	
	if resp.Error == nil {
		t.Error("Expected error, got none")
	}
}
