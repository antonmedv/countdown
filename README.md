# Countdown

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

## Notify

![notify_preview](https://imgur.com/FsZpGwy.png)

## Key binding

- `p` or `P`: To pause the countdown.
- `c` or `C`: To resume the countdown.
- `Esc` or `Ctrl+C`: To stop the countdown without running next command.

## Install

```bash
go get github.com/antonmedv/countdown
```

## License

MIT

## Credit

<div>Icons made by <a href="https://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
