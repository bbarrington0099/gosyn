package main

import (
	"testing"
	"os"
	"reflect"
	"strings"
)

// Parse Command
func TestParseCommand(t *testing.T) {
	// Save original command-line arguments
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	tests := []struct {
		name        string
		args        []string
		wantCmd     command
		wantErr     bool
		errContains string
	}{
		{
			name:        "no arguments",
			args:        []string{"gosyn"},
			wantCmd:     command{},
			wantErr:     true,
			errContains: "no command provided",
		},
		{
			name:    "action only",
			args:    []string{"gosyn", "help"},
			wantCmd: command{action: "help", args: []string{""}},
			wantErr: false,
		},
		{
			name:    "action with one argument",
			args:    []string{"gosyn", "listSections", "variables"},
			wantCmd: command{action: "listSections", args: []string{"variables"}},
			wantErr: false,
		},
		{
			name:    "action with multiple arguments",
			args:    []string{"gosyn", "tax", "functions", "closures"},
			wantCmd: command{action: "tax", args: []string{"functions", "closures"}},
			wantErr: false,
		},
		{
			name:    "action with empty arguments",
			args:    []string{"gosyn", "update", "", "test"},
			wantCmd: command{action: "update", args: []string{"test"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test arguments
			os.Args = tt.args

			gotCmd, err := parseCommand()

			if (err != nil) != tt.wantErr {
				t.Fatalf("parseCommand() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("parseCommand() error message %q should contain %q", 
						err.Error(), tt.errContains)
				}
				return
			}

			if !reflect.DeepEqual(gotCmd, tt.wantCmd) {
				t.Errorf("parseCommand() = %+v, want %+v", gotCmd, tt.wantCmd)
			}
		})
	}
}

// Execute Command
func TestExecuteCommand(t *testing.T) {
    // Backup and restore original command-line arguments
    origArgs := os.Args
    defer func() { os.Args = origArgs }()

    // Backup original function and restore when done
    oldInit := initializeSectionsFn
    defer func() { initializeSectionsFn = oldInit }()

    // Mock sections data
    testSections := []section{
        {
            name: "Variables",
            subsections: []subsection{
                {name: "Declaration", content: "var x int"},
                {name: "Types", content: "int, string"},
            },
        },
        {
            name: "Functions",
            subsections: []subsection{
                {name: "Declaration", content: "func f() {}"},
            },
        },
    }

    // Set our mock implementation
    initializeSectionsFn = func() []section { return testSections }

    tests := []struct {
		name        string
		args        []string
		wantOutput  string
		wantErr     bool
		errContains string
	}{
		// Help commands
		{
			name:       "help command",
			args:       []string{"gosyn", "help"},
			wantOutput: listActions(),
			wantErr:    false,
		},
		{
			name:       "help alias",
			args:       []string{"gosyn", "h"},
			wantOutput: listActions(),
			wantErr:    false,
		},

		// List sections commands
		{
			name:       "listSections full",
			args:       []string{"gosyn", "listSections"},
			wantOutput: listSections(testSections),
			wantErr:    false,
		},
		{
			name:       "listSections alias",
			args:       []string{"gosyn", "lsec"},
			wantOutput: listSections(testSections),
			wantErr:    false,
		},

		// List subsections commands
		{
			name:        "listSubsections valid",
			args:        []string{"gosyn", "lsub", "Variables"},
			wantOutput: func() string {
				output, err := listSubsections(testSections, "Variables")
				if err != nil {
					t.Fatalf("listSubsections() error = %v", err)
				}
				return output
			}(),
			wantErr:     false,
		},
		{
			name:        "listSubsections missing arg",
			args:        []string{"gosyn", "lsub"},
			wantErr:     true,
			errContains: "no section name provided",
		},
		{
			name:        "listSubsections invalid section",
			args:        []string{"gosyn", "lsub", "Invalid"},
			wantErr:     true,
			errContains: "section \"Invalid\" not found",
		},

		// Tax commands
		{
			name:       "tax valid",
			args:       []string{"gosyn", "Variables", "Declaration"},
			wantOutput: func() string {
				output, err := tax(testSections, "Variables", "Declaration")
				if err != nil {
					t.Fatalf("tax() error = %v", err)
				}
				return output
			}(),
			wantErr:    false,
		},
		{
			name:        "tax missing subsection",
			args:        []string{"gosyn", "Variables"},
			wantErr:     true,
			errContains: "no subsection name provided",
		},
		{
			name:        "tax invalid section with subsection",
			args:        []string{"gosyn", "Invalid", "Sub"},
			wantErr:     true,
			errContains: "section \"Invalid\" not found",
		},

		// Edge cases
		{
			name:        "unknown section or command",
			args:        []string{"gosyn", "unknown"},
			wantErr:     true,
			errContains: "section \"unknown\" not found",
		},
		{
			name:        "empty command",
			args:        []string{"gosyn"},
			wantErr:     true,
			errContains: "no command provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			got, err := executeCommand()

			if (err != nil) != tt.wantErr {
				t.Fatalf("executeCommand() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("executeCommand() error = %v, want contains %q", 
						err.Error(), tt.errContains)
				}
				// Verify ANSI codes are present in errors
				if tt.wantErr && !strings.Contains(err.Error(), BoldRed) {
					t.Error("Error message should contain red formatting")
				}
				return
			}

			if got != tt.wantOutput {
				t.Errorf("executeCommand() = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}