# Core File Privacy - Agent Instructions

## Build Commands

```bash
# Install dependencies
go mod tidy

# Build binary
go build -buildvcs=false -o cfp ./cmd/cfp

# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run integration demo
go run testdata/integration_demo.go

# Run complete test script
./testdata/complete_test.sh
```

## Project Structure

- `cmd/cfp/main.go` - CLI entry point with Cobra commands
- `internal/app/` - Command implementations (hide, show, verify, info, rekey, install, doctor)
- `internal/crypto/` - Cryptographic operations (AES-256-GCM, Argon2id, key generation)
- `internal/format/` - CFP1 container format (header, metadata, encode/decode)
- `internal/files/` - File operations (atomic write, hidden files, paths)
- `internal/archive/` - Directory archiving (tar create/extract)
- `internal/install/` - Installation logic (detect OS/shell, PATH management)
- `internal/prompt/` - User interaction (password input, confirmations, progress)
- `internal/version/` - Version information

## Key Features

1. **Encryption**: AES-256-GCM with Argon2id key derivation
2. **Security Profiles**: fast, default, paranoid (different Argon2id parameters)
3. **Key Generation**: Secure keys in format `cfp_XXXX-XXXX-...`
4. **File Hiding**: Dotfile prefix on Unix, hidden attribute on Windows
5. **Directory Support**: Archive directories with tar before encryption
6. **Auto-Installation**: Detects if not installed and offers to install
7. **PATH Management**: Adds install directory to shell config (bash, zsh, fish, PowerShell)

## Testing

All tests pass:
- Unit tests for crypto, format, files, archive, install
- Integration tests for complete workflows
- Manual testing with test scripts

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `golang.org/x/crypto/argon2` - Key derivation
- `golang.org/x/term` - Secure password input

## Release

Use GoReleaser for cross-platform builds:
```bash
goreleaser release --clean
```

This builds for:
- linux/amd64, linux/arm64
- windows/amd64, windows/arm64
- darwin/amd64, darwin/arm64
