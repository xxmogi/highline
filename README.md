# highline

A Powerline-style shell prompt generator written in Go.

`highline` renders colored, segmented prompts for bash and zsh — reading context directly from the filesystem and environment, with no dependency on external commands like `git` or `kubectl`.

```
 ~/dev/highline  main *  ⎈ prod-us-east
```

## Features

- **path** — current working directory, with home abbreviated to `~` and long paths shortened (`~/dev/.../project`)
- **git** — current branch name with dirty indicator (`*`), parsed directly from `.git/HEAD`
- **kube** — current Kubernetes context, read from `~/.kube/config` or `$KUBECONFIG`
- Segments that cannot be rendered (e.g. not in a git repo) are silently skipped
- Powerline separators with proper ANSI escape wrapping for bash and zsh
- Themeable colors and configurable segment order via a JSON config file

## Requirements

- Go 1.21+
- A terminal with a [Nerd Font](https://www.nerdfonts.com/) for Powerline glyphs

## Installation

```bash
git clone https://github.com/user/highline
cd highline
make install   # installs to $GOPATH/bin
```

Or build locally:

```bash
make build     # produces ./highline
```

## Shell Setup

### bash

Add to your `~/.bashrc`:

```bash
function _highline_ps1() {
    PS1="$(highline --shell bash) "
}

if [[ "$TERM" != "linux" ]]; then
    PROMPT_COMMAND="_highline_ps1${PROMPT_COMMAND:+; $PROMPT_COMMAND}"
fi
```

`PROMPT_COMMAND` runs `_highline_ps1` before each prompt is displayed, which sets `PS1` dynamically. The `"$TERM" != "linux"` guard skips the Powerline glyphs on the basic Linux console, where they are not supported.

### zsh

Add to your `~/.zshrc`:

```zsh
function _highline_precmd() {
    PS1="$(highline --shell zsh) "
}

precmd_functions+=(_highline_precmd)
```

`precmd_functions` is an array of hooks that zsh calls before each prompt. This is equivalent to bash's `PROMPT_COMMAND`.

## Usage

```
highline [--shell bash|zsh] [--config path] [segment ...]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--shell` | `bash` | Target shell (`bash` or `zsh`) |
| `--config` | `~/.config/highline/config.json` | Path to config file |

Segments can be passed as positional arguments in the desired order:

```bash
highline path git kube    # all three segments
highline git path         # git first, then path
highline path             # path only
```

## Configuration

`highline` reads `~/.config/highline/config.json` at startup (respects `$XDG_CONFIG_HOME`). All fields are optional.

```json
{
  "shell": "bash",
  "segments": ["path", "git", "kube"],
  "theme": {
    "path": { "fg": 7, "bg": 4 },
    "git":  { "fg": 0, "bg": 2 },
    "kube": { "fg": 7, "bg": 6 }
  },
  "path": {
    "max_depth": 4
  }
}
```

### Config fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `shell` | string | `"bash"` | Target shell |
| `segments` | []string | `["path","git","kube"]` | Segments to show and their order |
| `theme.<name>.fg` | int (0–15) | segment default | Foreground color (ANSI basic 16 colors) |
| `theme.<name>.bg` | int (0–15) | segment default | Background color (ANSI basic 16 colors) |
| `path.max_depth` | int | `4` | Path components before middle is collapsed |

### Priority

CLI flags take precedence over the config file:

- `--shell` overrides `config.shell`
- Segment arguments override `config.segments`

If the config file does not exist, `highline` runs with defaults. If it cannot be parsed, a warning is printed to stderr and defaults are used.

### Default colors

| Segment | Background | Foreground |
|---------|------------|------------|
| path    | Blue (4)   | White (7)  |
| git     | Yellow (3) | Black (0)  |
| kube    | Cyan (6)   | White (7)  |

## Development

```bash
make build    # build binary
make test     # run all tests
make lint     # run go vet
make clean    # remove build artifacts
make help     # list all targets
```

## Project structure

```
cmd/highline/       # entry point
internal/
  config/           # config file loading
  segment/          # path, git, kube segments
  renderer/         # Powerline prompt renderer
  color/            # ANSI color constants
specs/              # specification documents
```
