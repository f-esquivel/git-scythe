package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/franklinesquivel/git-scythe/internal/types"
)

type mockGitter struct {
	getMergedBranchesFunc func() ([]string, error)
	deleteBranchFunc      func(branchName string, force bool) error
}

func (m *mockGitter) GetMergedBranches() ([]string, error) {
	if m.getMergedBranchesFunc != nil {
		return m.getMergedBranchesFunc()
	}
	return []string{}, nil
}

func (m *mockGitter) DeleteBranch(branchName string, force bool) error {
	if m.deleteBranchFunc != nil {
		return m.deleteBranchFunc(branchName, force)
	}
	return nil
}

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		git            types.Gitter
		input          string
		expectedOutput string
		expectedErr    bool
	}{
		{
			name: "successful deletion",
			git: &mockGitter{
				getMergedBranchesFunc: func() ([]string, error) {
					return []string{"branch1", "branch2"}, nil
				},
				deleteBranchFunc: func(branchName string, force bool) error {
					return nil
				},
			},
			input:          "Y\n",
			expectedOutput: "The following merged branches will be deleted:\n  > branch1\n  > branch2\n\nDo you want to continue? (Y/n) \nDeleting 2 branches...\nDeleted 2 branches.\n",
			expectedErr:    false,
		},
		{
			name: "user cancels",
			git: &mockGitter{
				getMergedBranchesFunc: func() ([]string, error) {
					return []string{"branch1"}, nil
				},
			},
			input:          "n\n",
			expectedOutput: "The following merged branches will be deleted:\n  > branch1\n\nDo you want to continue? (Y/n) Operation cancelled.\n",
			expectedErr:    false,
		},
		{
			name: "no branches to delete",
			git: &mockGitter{
				getMergedBranchesFunc: func() ([]string, error) {
					return []string{}, nil
				},
			},
			input:          "",
			expectedOutput: "No merged branches to display.\n",
			expectedErr:    false,
		},
		{
			name: "error fetching branches",
			git: &mockGitter{
				getMergedBranchesFunc: func() ([]string, error) {
					return nil, errors.New("git fetch error")
				},
			},
			input:          "",
			expectedOutput: "",
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}

			err := run(tt.git, reader, writer)

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if writer.String() != tt.expectedOutput {
				t.Errorf("expected output:\n%q\ngot:\n%q", tt.expectedOutput, writer.String())
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedErr    bool
	}{
		{
			name:           "version flag",
			args:           []string{"git-scythe", "--version"},
			expectedOutput: "dev\n",
			expectedErr:    false,
		},
		{
			name:           "v flag",
			args:           []string{"git-scythe", "-v"},
			expectedOutput: "dev\n",
			expectedErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}

			err := execute(tt.args, writer)

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if writer.String() != tt.expectedOutput {
				t.Errorf("expected output:\n%q\ngot:\n%q", tt.expectedOutput, writer.String())
			}
		})
	}
}