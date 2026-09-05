# Contributing

Thank you for your interest in MasjidPi!

We welcome:

- Bug reports
- Feature requests
- Documentation improvements
- Code contributions
- Testing on Raspberry Pi hardware

## Development Principles

- Keep dependencies minimal.
- Prefer the Go standard library.
- Build for reliability.
- Keep the project beginner-friendly.
- Write clear commit messages.
- Document new features.

Every pull request should leave the project in a working state.

## Before Opening a Pull Request

Run the checks that apply to your change:

```bash
make test

cd backend
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
```

Frontend tests are standalone Node.js scripts:

```bash
for test_file in tests/*test.js; do
    node "$test_file" || exit 1
done
```

Also run `git diff --check`, update the relevant documentation, and avoid
committing credentials, device state or unredacted upstream payloads containing
personal information.
