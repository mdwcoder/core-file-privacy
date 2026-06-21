# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue in Core File Privacy, please follow these steps:

### Private Disclosure

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please report security issues via email to: **[security@example.com]** (replace with actual security contact)

Please include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution Target**: Within 30 days

### What We Consider Security Issues

- Encryption implementation flaws
- Key derivation weaknesses
- Path traversal vulnerabilities
- Buffer overflows or memory safety issues
- Authentication bypass
- Information leakage (passwords, keys in logs)
- Insecure random number generation
- Format specification vulnerabilities

### What We Don't Consider Security Issues

- Features not yet implemented
- Performance issues (unless they cause DoS)
- Usability concerns
- Missing features

## Security Design

### Encryption

- **Algorithm**: AES-256-GCM (authenticated encryption)
- **Key Derivation**: Argon2id (memory-hard KDF)
- **Random Generation**: crypto/rand (cryptographically secure)

### Threat Model

Core File Privacy is designed to protect against:
- Unauthorized access to file contents
- Tampering with encrypted files
- Brute-force attacks on passwords

### Limitations

Core File Privacy does **not** protect against:
- Malware on the same system
- Memory dumps while the program is running
- Loss of password/key (no recovery possible)
- Forensic analysis of storage media (no secure deletion guarantee)

## Security Best Practices

### For Users

1. **Use Strong Passwords**: If not using generated keys, use long, random passwords
2. **Backup Keys**: Store generated keys in a password manager
3. **Verify Before Trust**: Use `cfp verify` to check container integrity
4. **Update Regularly**: Keep CFP updated to the latest version
5. **Test Recovery**: Periodically test that you can decrypt your files

### For Developers

1. **No Password Logging**: Never log passwords or keys
2. **Secure Memory**: Clear sensitive data from memory when possible
3. **Input Validation**: Validate all inputs, especially file paths
4. **Use Standard Libraries**: Prefer well-audited crypto libraries
5. **Code Review**: All crypto-related changes require review

## Cryptographic Parameters

### Security Profiles

| Profile  | Memory | Iterations | Parallelism | Salt Size | Nonce Size |
|----------|--------|------------|-------------|-----------|------------|
| fast     | 32 MB  | 2          | 2           | 16 bytes  | 12 bytes   |
| default  | 64 MB  | 3          | 4           | 16 bytes  | 12 bytes   |
| paranoid | 256 MB | 4          | 4           | 16 bytes  | 12 bytes   |

### Key Size

- **Derived Key**: 256 bits (32 bytes)
- **Salt**: 128 bits (16 bytes)
- **Nonce**: 96 bits (12 bytes)

## Audit Status

This software has not yet been independently audited. Use at your own risk for critical data.

## Known Limitations

1. **No Forward Secrecy**: If a password is compromised, all past and future files encrypted with that password are at risk
2. **No Key Rotation**: Changing passwords requires full re-encryption
3. **Single Password**: Each container is protected by a single password/key
4. **No Multi-Party**: Containers cannot be shared among multiple users without sharing the password

## Future Security Enhancements

Potential future improvements:
- Memory clearing for sensitive data
- Hardware security module (HSM) support
- Multi-party threshold encryption
- Forward secrecy mechanisms
- Independent security audit
