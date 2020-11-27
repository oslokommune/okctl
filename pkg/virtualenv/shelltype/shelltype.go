// Package shelltype contains constants to identify a shell (bash, zsh, etc)
package shelltype

// ShellType enumerates shells we recognize
type ShellType string

const (
	// Bash is a constant that identifies the Bash shell
	Bash ShellType = "bash"

	// Zsh is a constant that identifies the Zsh shell
	Zsh ShellType = "zsh"

	// Unknown is a constant that is identifies the case when an unknown shell is used
	Unknown ShellType = "unknown"
)
