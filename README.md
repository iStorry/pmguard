# pmguard

You've done it. We've all done it.

```
$ pnpm install
...
Lockfile is up to date, resolution step is skipped
...
```

Then you look up and realize — *this is a bun project.*

Now you've got a `pnpm-lock.yaml` sitting in a repo that's never seen pnpm in its life. Your teammates are confused. Your CI is broken. Your lockfile is cooked.

**pmguard stops you before that happens.**

---

## What it does

It wraps `npm`, `pnpm`, `yarn`, and `bun` in your shell. Same commands you always type — but if you're in the wrong project, it either warns you or silently runs the right one instead.

```
$ cd my-bun-project
$ pnpm install

⚠️  pmguard: This project uses bun, but you ran pnpm.
   Run: bun install
```

Or if you're a "just fix it for me" type:

```
$ pnpm install
⚡ pmguard: redirecting pnpm → bun
```

Zero new commands. Zero new habits. Just add one line to your shell and forget about it.

---

## Install

**Mac:**
```bash
brew install iStorry/pmguard/pmguard
```

**Linux / Mac without Homebrew:**
```bash
curl -sSL https://raw.githubusercontent.com/iStorry/pmguard/main/scripts/install.sh | sh
```

**Windows:**

Grab the latest binary from [Releases](https://github.com/iStorry/pmguard/releases) and add it to your PATH.

**Then activate it:**

```bash
# Add to your ~/.zshrc or ~/.bashrc
eval "$(pmguard install-hooks)"
```

Restart your shell and you're done. Seriously, that's it.

---

## Modes

```bash
pmguard config set-mode warn      # yell at me but don't do anything (default)
pmguard config set-mode redirect  # just fix it silently
pmguard config get                # what mode am I in?
```

Config lives at `~/.config/pmguard/config.yaml`.

---

## How it detects your package manager

It walks up your directory tree looking for a lockfile:

| Lockfile | Package manager |
|----------|----------------|
| `bun.lockb` / `bun.lock` | bun |
| `pnpm-lock.yaml` | pnpm |
| `yarn.lock` | yarn |
| `package-lock.json` | npm |

No lockfile? It also checks the `packageManager` field in `package.json`. Nothing found at all? It gets out of your way and passes the command through untouched.

---

## Redirect mode flag translation

When redirecting, pmguard translates flags so the command actually works:

| From | To | Original | Becomes |
|------|----|----------|---------|
| pnpm | bun | `-D` | `-d` |
| npm | bun | `--save-dev` | `-d` |
| bun | pnpm | `-d` | `-D` |
| bun | npm | `-d` | `--save-dev` |

---

## Uninstall

```bash
# Homebrew
brew uninstall pmguard

# Manual
rm /usr/local/bin/pmguard
```

And remove the `eval "$(pmguard install-hooks)"` line from your shell config.

---

## Contributing

This is a small focused tool and PRs are very welcome. Good places to start:

- More flag translations in `internal/remap/remap.go`
- Fish shell support
- PowerShell support
- Per-project config via `.pmguard.yaml`

---

## License

MIT
