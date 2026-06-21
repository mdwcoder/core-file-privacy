# Core File Privacy

Encrypt, hide and restore private files from your terminal.

**Core File Privacy** (`cfp`) is a cross-platform CLI tool for encrypting, hiding, verifying, decrypting, and changing passwords of files or folders using a user-provided password or a secure auto-generated key.

## Features

- **AES-256-GCM** encryption with **Argon2id** key derivation
- Encrypt files and directories into `.cfp` containers
- Generate secure keys in a convenient format
- Hide encrypted files using OS-specific mechanisms
- Verify container integrity without decryption
- Change passwords/keys without re-encrypting from scratch
- Multi-platform: Windows, Linux, macOS
- Automatic installation and PATH configuration

## Installation

### Quick Install

Download the latest release for your platform from [GitHub Releases](https://github.com/yourusername/core-file-privacy/releases).

#### Linux/macOS

```bash
# Download and extract
tar -xzf cfp_linux_amd64.tar.gz

# Install
./cfp install
```

#### Windows

```powershell
# Download and extract the .zip file
# Then run:
.\cfp.exe install
```

### Manual Installation

Simply run the binary from any location. If it's not in the expected installation directory, it will offer to install itself:

```bash
./cfp
```

The installer will:
1. Copy the binary to `~/core-file-privacy/` (Linux/macOS) or `%USERPROFILE%\core-file-privacy\` (Windows)
2. Add the directory to your PATH
3. Configure your shell (bash, zsh, fish)

## Quick Start

```bash
# Encrypt a file with a generated key
cfp hide .env --generate-key --hidden

# Decrypt a file
cfp show .env.cfp

# Verify a container
cfp verify .env.cfp

# Change the password/key
cfp rekey .env.cfp --generate-key
```

## Usage

### Encrypt a File

```bash
# With password
cfp hide secret.txt

# With generated key
cfp hide secret.txt --generate-key

# Hide the encrypted file
cfp hide secret.txt --generate-key --hidden

# Keep the original file
cfp hide secret.txt --keep

# Delete the original after encryption
cfp hide secret.txt --delete-original
```

### Encrypt a Directory

```bash
# Archive and encrypt a directory
cfp hide private-folder --archive --generate-key
```

### Decrypt a Container

```bash
# Decrypt to original location
cfp show secret.txt.cfp

# Decrypt to specific location
cfp show secret.txt.cfp --output restored.txt

# Force overwrite existing files
cfp show secret.txt.cfp --force

# Delete the encrypted file after decryption
cfp show secret.txt.cfp --delete-encrypted
```

### Verify Container Integrity

```bash
cfp verify secret.txt.cfp
```

This will ask for the password/key and verify it's correct without writing the decrypted file.

### View Container Information

```bash
# Show public information (no password required)
cfp info secret.txt.cfp

# Show private metadata (requires password)
cfp info secret.txt.cfp --unlock
```

### Change Password/Key

```bash
# Change to a new password
cfp rekey secret.txt.cfp

# Generate a new key
cfp rekey secret.txt.cfp --generate-key

# Change security profile
cfp rekey secret.txt.cfp --profile paranoid
```

### System Commands

```bash
# Install cfp
cfp install

# Install with options
cfp install --force          # Force reinstall
cfp install --no-path        # Don't modify PATH
cfp install --target /custom/path  # Custom location

# Uninstall cfp
cfp uninstall

# Check system status
cfp doctor

# Show version
cfp version
```

## Security Profiles

Three security profiles are available, controlling the Argon2id parameters:

| Profile  | Memory | Iterations | Parallelism | Use Case |
|----------|--------|------------|-------------|----------|
| `fast`   | 32 MB  | 2          | 2           | Quick operations, lower security |
| `default`| 64 MB  | 3          | 4           | Balanced security and performance |
| `paranoid`| 256 MB| 4          | 4           | Maximum security, slower |

Default profile is `default`.

## Generated Key Format

When using `--generate-key`, a secure key is generated in this format:

```
cfp_V7KQ-9X2M-P4DA-R8JN-T6HZ-YW3C-LB5S-EQ1F
```

- Uses `crypto/rand` for true randomness
- 8 groups of 4 alphanumeric characters
- Can be used as password for `show`, `verify`, and `rekey`
- **Save it in a password manager** - CFP cannot recover it if lost

## File Format

The `.cfp` container format (CFP1):

```
[magic: CFP1]
[header_length: uint32]
[header_json]
[ciphertext]
```

The header contains:
- Format version
- Cipher and KDF information
- Security profile
- KDF parameters (salt, iterations, memory, parallelism)
- Nonce
- Creation timestamp

The encrypted payload contains:
- Metadata (original filename, type, permissions, size)
- Original file data (or tar archive for directories)

## Security Considerations

### Encryption

- **AES-256-GCM**: Authenticated encryption ensures both confidentiality and integrity
- **Argon2id**: Memory-hard key derivation function, resistant to GPU/ASIC attacks
- **crypto/rand**: Cryptographically secure random number generation
- **No password recovery**: If you lose your password/key, your data is gone forever

### Important Warnings

**No Password Recovery**: CFP cannot recover your password or generated key. If you lose it, your encrypted data is permanently inaccessible. Always backup your keys in a password manager.

**No Guaranteed Secure Deletion**: When using `--delete-original`, CFP deletes files using standard OS mechanisms. This does **not** guarantee secure deletion on:
- SSDs (wear leveling, over-provisioning)
- Journaling filesystems
- Systems with snapshots or backups
- Cloud-synced folders (Dropbox, OneDrive, etc.)

**Hiding is Not Security**: The `--hidden` option makes files less visible but does not provide real security. The encryption is what protects your data.

**No Cloud Sync**: CFP does not sync encrypted files to the cloud. Manage your `.cfp` files yourself.

## Compatibility

### Operating Systems

- **Linux**: x86_64, ARM64
- **macOS**: x86_64, ARM64 (Apple Silicon)
- **Windows**: x86_64, ARM64

### Shells

The installer supports:
- bash
- zsh
- fish (Linux/macOS)
- PowerShell (Windows)
- CMD (Windows)

## Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/core-file-privacy.git
cd core-file-privacy

# Build
go build -o cfp ./cmd/cfp

# Run tests
go test ./...

# Install locally
go install ./cmd/cfp
```

## Contributing

Contributions are welcome! Please open an issue or pull request on GitHub.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/core-file-privacy/issues)
- **Security**: See [SECURITY.md](SECURITY.md)
