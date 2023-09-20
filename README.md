# Countdown

<p align="center">
  <br>
  <img src="demo.gif" width="600" alt="Countdown Demo">
  <br>
</p>

## Install

```sh
brew install countdown
```

```sh
go install github.com/antonmedv/countdown@latest
```

Or download prebuilt binary from [releases](https://github.com/antonmedv/countdown/releases).

## Usage

Specify duration in Go format `1h2m3s` or a target time: `02:15pm`, `14:15`.

```sh
countdown 25s
countdown 11:32
```

Add a command with `&&` to run after the countdown.

```sh
countdown 1m30s && say "Hello, world"
```

Count from up from the zero.

```sh
countdown -up 30s
```

Announce (via macOS `say` command) last 10 seconds.

```sh
countdown -say 10s
```

## Key binding

- `Space`: Pause/Resume the countdown.
- `Esc` or `Ctrl+C`: Stop the countdown without running the next command.

## Become a sponsor

Every line of code in my repositories 📖 signifies my unwavering commitment to open source 💡. Your support 🤝 ensures
these projects keep thriving, innovating, and benefiting all 💼. If my work has ever resonated 🎵 or helped you, kindly
consider showing love ❤️ by sponsoring. [**🚀 Sponsor Me Today! 🚀**](https://github.com/sponsors/antonmedv)

## License

[MIT](LICENSE)
