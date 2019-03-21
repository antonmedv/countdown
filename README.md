# Terminal Countdown Timer

<p align="center"><img src="https://user-images.githubusercontent.com/141232/54696023-9ed03e00-4b5d-11e9-9c7b-d6f67691e70c.gif" width="450" alt="Screen shot"></p>

## Usage

Specify duration in go format `1h2m3s`.

```bash
countdown 25s
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
