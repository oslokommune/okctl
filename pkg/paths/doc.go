// Package paths knows how to construct paths to different locations used by okctl.
//
// Full paths are constructed by combining the GetAbsoluteRepositoryRootDirectory() with a relative path constructor function.
//
// - All relative path construction functions MUST be relative to the users IAC repository root.
package paths
