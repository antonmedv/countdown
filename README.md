# Terminal Countdown Timer

## Usage

Specify duration in go format `1h2m3s`.

```bash
countdown 5s
```

Add command with `&&` to run after countdown.

```bash
countdown 1m30s && say "Hello, world"
```

Press `Esc` or `Ctrl+C` to stop countdown without running next command.

## Install

```bash
go get github.com/antonmedv/countdown
``` 

## License

MIT
