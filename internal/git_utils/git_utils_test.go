package git_utils

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/franklinesquivel/git-scythe/internal/types"
)

// mockCmd is a mock implementation of the types.Cmd interface for testing.
type mockCmd struct {
	stdoutPipeFunc func() (io.ReadCloser, error)
	startFunc      func() error
	waitFunc       func() error
	outputFunc     func() ([]byte, error)
	setStdinFunc   func(io.Reader)
}

func (m *mockCmd) StdoutPipe() (io.ReadCloser, error) {
	if m.stdoutPipeFunc != nil {
		return m.stdoutPipeFunc()
	}
	return nil, nil
}

func (m *mockCmd) Start() error {
	if m.startFunc != nil {
		return m.startFunc()
	}
	return nil
}

func (m *mockCmd) Wait() error {
	if m.waitFunc != nil {
		return m.waitFunc()
	}
	return nil
}

func (m *mockCmd) Output() ([]byte, error) {
	if m.outputFunc != nil {
		return m.outputFunc()
	}
	return nil, nil
}

func (m *mockCmd) SetStdin(stdin io.Reader) {
	if m.setStdinFunc != nil {
		m.setStdinFunc(stdin)
	}
}

// mockExecutor is a mock implementation of the types.CommandExecutor interface.
type mockExecutor struct {
	commandFunc func(name string, arg ...string) types.Cmd
}

func (m *mockExecutor) Command(name string, arg ...string) types.Cmd {
	if m.commandFunc != nil {
		return m.commandFunc(name, arg...)
	}
	return &mockCmd{}
}

func TestGetDefaultBranch(t *testing.T) {
	tests := []struct {
		name        string
		executor    types.CommandExecutor
		expected    string
		expectedErr bool
	}{
		{
			name: "Success",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "remote" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("origin\n"), nil
							},
						}
					}
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return []byte("origin/main\n"), nil
						},
					}
				},
			},
			expected:    "main",
			expectedErr: false,
		},
		{
			name: "Error",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("git error")
						},
					}
				},
			},
			expected:    "",
			expectedErr: true,
		},
		{
			name: "Git not in path",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New(`exec: "git": executable file not found in $PATH`)
						},
					}
				},
			},
			expected:    "",
			expectedErr: true,
		},
		{
			name: "Empty repository",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("fatal: ambiguous argument 'origin/HEAD': unknown revision or path not in the working tree")
						},
					}
				},
			},
			expected:    "",
			expectedErr: true,
		},
		{
			name: "non-origin remote",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "remote" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("upstream\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "rev-parse" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("upstream/main\n"), nil
							},
						}
					}
					return &mockCmd{}
				},
			},
			expected:    "main",
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := New(tt.executor)
			branch, err := git.GetDefaultBranch()

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if branch != tt.expected {
				t.Errorf("expected branch: %s, got: %s", tt.expected, branch)
			}
		})
	}
}

func TestGetMergedBranches(t *testing.T) {
	tests := []struct {
		name           string
		executor       types.CommandExecutor
		expectedResult []string
		expectedError  bool
	}{
		{
			name: "successful execution",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "remote" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("origin\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "rev-parse" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("main\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "branch" {
						return &mockCmd{
							stdoutPipeFunc: func() (io.ReadCloser, error) {
								return io.NopCloser(nil), nil
							},
						}
					}
					if name == "grep" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("  branch1\n  branch2\n"), nil
							},
						}
					}
					return &mockCmd{}
				},
			},
			expectedResult: []string{"branch1", "branch2"},
			expectedError:  false,
		},
		{
			name: "error in getting default branch",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("git error")
						},
					}
				},
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "no merged branches",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "remote" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("origin\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "rev-parse" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("main\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "branch" {
						return &mockCmd{
							stdoutPipeFunc: func() (io.ReadCloser, error) {
								return io.NopCloser(nil), nil
							},
						}
					}
					if name == "grep" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								// grep will return an empty output and a non-zero exit code
								// when no lines match, which is expected here.
								return []byte(""), nil
							},
						}
					}
					return &mockCmd{}
				},
			},
			expectedResult: []string{},
			expectedError:  false,
		},
		{
			name: "unusual branch names",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "remote" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("origin\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "rev-parse" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("main\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "branch" {
						return &mockCmd{
							stdoutPipeFunc: func() (io.ReadCloser, error) {
								return io.NopCloser(nil), nil
							},
						}
					}
					if name == "grep" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("  feature/new-login\n  fix/a-bug\n"), nil
							},
						}
					}
					return &mockCmd{}
				},
			},
			expectedResult: []string{"feature/new-login", "fix/a-bug"},
			expectedError:  false,
		},
		{
			name: "detached HEAD state",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "remote" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("origin\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "rev-parse" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("main\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "branch" {
						return &mockCmd{
							stdoutPipeFunc: func() (io.ReadCloser, error) {
								return io.NopCloser(nil), nil
							},
						}
					}
					if name == "grep" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								// Simulate grep filtering out the detached head line
								return []byte("  a-branch\n  another-branch\n"), nil
							},
						}
					}
					return &mockCmd{}
				},
			},
			expectedResult: []string{"a-branch", "another-branch"},
			expectedError:  false,
		},
		{
			name: "branch names with spaces",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "remote" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("origin\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "rev-parse" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("main\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "branch" {
						return &mockCmd{
							stdoutPipeFunc: func() (io.ReadCloser, error) {
								return io.NopCloser(nil), nil
							},
						}
					}
					if name == "grep" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("  a feature branch\n  another branch with spaces\n"), nil
							},
						}
					}
					return &mockCmd{}
				},
			},
			expectedResult: []string{"a feature branch", "another branch with spaces"},
			expectedError:  false,
		},
		{
			name: "branch name conflicts with flag",
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "remote" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("origin\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "rev-parse" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("main\n"), nil
							},
						}
					}
					if name == "git" && arg[0] == "branch" {
						return &mockCmd{
							stdoutPipeFunc: func() (io.ReadCloser, error) {
								return io.NopCloser(nil), nil
							},
						}
					}
					if name == "grep" {
						// This will likely fail without the `--` separator
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte("  --version\n"), nil
							},
						}
					}
					return &mockCmd{}
				},
			},
			expectedResult: []string{"--version"},
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := New(tt.executor)
			result, err := git.GetMergedBranches()

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("Expected result: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}

func TestDeleteBranch(t *testing.T) {
	tests := []struct {
		name        string
		branchName  string
		force       bool
		executor    types.CommandExecutor
		expectedErr bool
	}{
		{
			name:       "delete branch successfully",
			branchName: "feature/test-branch",
			force:      false,
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "branch" && arg[1] == "-d" && arg[2] == "feature/test-branch" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte(""), nil
							},
						}
					}
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("unexpected command")
						},
					}
				},
			},
			expectedErr: false,
		},
		{
			name:       "force delete branch successfully",
			branchName: "feature/test-branch",
			force:      true,
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "branch" && arg[1] == "-D" && arg[2] == "feature/test-branch" {
						return &mockCmd{

							outputFunc: func() ([]byte, error) {
								return []byte(""), nil
							},
						}
					}
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("unexpected command")
						},
					}
				},
			},
			expectedErr: false,
		},
		{
			name:       "delete branch with error",
			branchName: "feature/test-branch",
			force:      false,
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("git error")
						},
					}
				},
			},
			expectedErr: true,
		},
		{
			name:       "delete unmerged branch without force",
			branchName: "unmerged-branch",
			force:      false,
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("error: The branch 'unmerged-branch' is not fully merged.\nIf you are sure you want to delete it, run 'git branch -D unmerged-branch'.")
						},
					}
				},
			},
			expectedErr: true,
		},
		{
			name:       "delete non-existent branch",
			branchName: "non-existent-branch",
			force:      false,
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("error: branch 'non-existent-branch' not found.")
						},
					}
				},
			},
			expectedErr: true,
		},
		{
			name:       "delete current branch",
			branchName: "main",
			force:      false,
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("error: Cannot delete branch 'main' checked out at...")
						},
					}
				},
			},
			expectedErr: true,
		},
		{
			name:       "force delete unmerged branch",
			branchName: "unmerged-branch",
			force:      true,
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					if name == "git" && arg[0] == "branch" && arg[1] == "-D" && arg[2] == "unmerged-branch" {
						return &mockCmd{
							outputFunc: func() ([]byte, error) {
								return []byte(""), nil
							},
						}
					}
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("unexpected command")
						},
					}
				},
			},
			expectedErr: false,
		},
		{
			name:       "delete branch with empty name",
			branchName: "",
			force:      false,
			executor: &mockExecutor{
				commandFunc: func(name string, arg ...string) types.Cmd {
					return &mockCmd{
						outputFunc: func() ([]byte, error) {
							return nil, errors.New("fatal: '' is not a valid branch name.")
						},
					}
				},
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := New(tt.executor)
			err := git.DeleteBranch(tt.branchName, tt.force)

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
		})
	}
}
