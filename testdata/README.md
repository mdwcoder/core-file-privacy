# Test Data

This directory contains test files and integration demos for Core File Privacy.

## Files

- `integration_demo.go`: Complete integration test demonstrating all core functionality
  - Key generation
  - File encryption/decryption
  - Directory archiving
  - File hiding
  - Container info reading

## Running the Integration Demo

```bash
go run testdata/integration_demo.go
```

This will create test files in `/tmp/cfp-integration-test/` and verify all core functionality.

## Note

The `.cfp` files in this directory are generated during testing and can be safely deleted.
