// Package config provides runtime configuration for the envdiff tool.
//
// It defines the Config struct that controls how envdiff compares .env files,
// formats its output, and handles secret masking. Configuration can be
// populated from CLI flags via FromArgs or built programmatically.
//
// Supported output formats:
//
//	- text  (default): human-readable line-by-line diff
//	- table          : aligned table view with truncated values
//	- json           : machine-readable JSON, optionally indented
//
// Example usage:
//
//	cfg, err := config.FromArgs(os.Args[1:])
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if err := cfg.Validate(); err != nil {
//	    log.Fatal(err)
//	}
package config
