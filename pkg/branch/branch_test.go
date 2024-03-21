package branch_test

import (
	"os/exec"
	"testing"

	"github.com/elhub/gh-dxp/pkg/branch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Command(name string, args ...string) (string, error) {
	argsCalled := m.Called(name, args)
	return argsCalled.String(0), argsCalled.Error(1)
}

type FakeExitError struct {
	code int
}

func (f FakeExitError) Error() string {
	return "Fake error"
}

func (f FakeExitError) ExitCode() int {
	return f.code
}

// This tests both branch.CheckoutBranch and branch.BranchExists
func TestCheckoutBranch(t *testing.T) {
	t.Run("should checkout to existing branch", func(t *testing.T) {
		mockExec := new(MockExecutor)

		// Set up expectation
		mockExec.On("Command", "git", []string{"checkout", "existing-branch"}).Return("", nil)
		mockExec.On("Command", "git", []string{"show-ref", "--verify", "--quiet", "refs/heads/existing-branch"}).Return("", nil)

		// Call the method under test
		err := branch.CheckoutBranch(mockExec, "existing-branch")

		// Assert that the expectations were met
		mockExec.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("should checkout to new branch if it does not exist", func(t *testing.T) {
		mockExec := new(MockExecutor)

		// Start a process that will exit with code 1
		// Mocking a os.ProcessState is not possible, so we need to mock the ExitError
		cmd := exec.Command("go", "version", "nonexistent")
		cmdErr := cmd.Run()

		// Set up expectation
		mockExec.On("Command", "git", []string{"checkout", "-b", "nonexistent-branch"}).Return("", nil)
		mockExec.On("Command", "git", []string{"show-ref", "--verify", "--quiet", "refs/heads/nonexistent-branch"}).Return("", cmdErr)

		// Call the method under test
		err := branch.CheckoutBranch(mockExec, "nonexistent-branch")

		// Assert that the expectations were met
		mockExec.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("should throw error if checkout fails to existing branch", func(t *testing.T) {
		mockExec := new(MockExecutor)

		// Start a process that will exit with code 1
		// Mocking a os.ProcessState is not possible, so we need to mock the ExitError
		cmd := exec.Command("go", "version", "nonexistent")
		cmdErr := cmd.Run()

		// Set up expectation
		mockExec.On("Command", "git", []string{"checkout", "-b", "nonexistent-branch"}).Return("", FakeExitError{code: 1})
		mockExec.On("Command", "git", []string{"show-ref", "--verify", "--quiet", "refs/heads/nonexistent-branch"}).Return("", cmdErr)

		// Call the method under test
		err := branch.CheckoutBranch(mockExec, "nonexistent-branch")

		// Assert that the expectations were met
		mockExec.AssertExpectations(t)
		assert.Error(t, err)
	})

	t.Run("should throw error if checkout fails to new branch", func(t *testing.T) {
		mockExec := new(MockExecutor)

		// Start a process that will exit with code 1
		// Mocking a os.ProcessState is not possible, so we need to mock the ExitError
		cmd := exec.Command("go", "version", "nonexistent")
		cmdErr := cmd.Run()

		// Set up expectation
		mockExec.On("Command", "git", []string{"checkout", "-b", "nonexistent-branch"}).Return("", FakeExitError{code: 1})
		mockExec.On("Command", "git", []string{"show-ref", "--verify", "--quiet", "refs/heads/nonexistent-branch"}).Return("", cmdErr)

		// Call the method under test
		err := branch.CheckoutBranch(mockExec, "nonexistent-branch")

		// Assert that the expectations were met
		mockExec.AssertExpectations(t)
		assert.Error(t, err)
	})

	t.Run("should throw error if branch exists check fails", func(t *testing.T) {
		mockExec := new(MockExecutor)

		// Set up expectation
		mockExec.On("Command", "git", []string{"show-ref", "--verify", "--quiet", "refs/heads/failing-branch"}).Return("", FakeExitError{code: 1})

		// Call the method under test
		err := branch.CheckoutBranch(mockExec, "failing-branch")

		// Assert that the expectations were met
		mockExec.AssertExpectations(t)
		assert.Error(t, err)
	})

}
