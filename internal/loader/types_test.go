package loader

import "testing"

func TestExtractStackName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid single word", "10-nginx", "nginx"},
		{"valid multi word", "20-test-stack", "test-stack"},
		{"valid with underscores", "30-my_stack_name", "my_stack_name"},
		{"invalid no dash", "invalidname", ""},
		{"invalid empty after dash", "40-", ""},
		{"valid many dashes", "50-stack-with-many-dashes", "stack-with-many-dashes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractStackName(tt.input)
			if result != tt.expected {
				t.Errorf("extractStackName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected StackStatus
	}{
		{"running", "running", StackStatusRunning},
		{"running with count", "running(2)", StackStatusRunning},
		{"exited maps to stopped", "exited", StackStatusStopped},
		{"stopped", "stopped", StackStatusStopped},
		{"unknown maps to down", "unknown", StackStatusDown},
		{"empty maps to down", "", StackStatusDown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeStatus(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeStatus(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestActionIsValid(t *testing.T) {
	validActions := []Action{ActionStart, ActionStop, ActionRestart, ActionReload, ActionDown, ActionList}
	for _, action := range validActions {
		if !action.IsValid() {
			t.Errorf("Expected %q to be valid", action)
		}
	}

	invalidActions := []Action{"invalid", "unknown", ""}
	for _, action := range invalidActions {
		if action.IsValid() {
			t.Errorf("Expected %q to be invalid", action)
		}
	}
}
