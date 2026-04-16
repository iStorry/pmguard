# pmguard

> Use the right package manager. Always.

pmguard intercepts `npm`, `pnpm`, `yarn`, and `bun` commands and warns (or silently redirects) if you're using the wrong one for the current project — detected from lockfiles and the `packageManager` field in `package.json`.

No new commands to learn. Just add one line to your shell config.

---

## How it works

```
$ cd my-bun-project
$ pnpm install

⚠️  pmguard: This project uses bun, but you ran pnpm.
   Run: bun install
```

Or in redirect mode:

```
$ pnpm install
⚡ pmguard: redirecting pnpm → bun
```

---

## Install

**Mac (Homebrew):**

```bash
brew install istorry/pmguard/pmguard
```

**Linux / Mac without Homebrew:**

```bash
curl -sSL https://raw.githubusercontent.com/istorry/pmguard/main/scripts/install.sh | sh
```

**Windows:**

Download the latest binary from [GitHub Releases](https://github.com/istorry/pmguard/releases) and add it to your PATH.

**Activate in your shell:**

Add this to your `~/.zshrc` or `~/.bashrc`:

```bash
eval "$(pmguard install-hooks)"
```

Then restart your shell or run `source ~/.zshrc`. That's it.

---

## Configuration

```bash
# Default: warn and stop when wrong PM is used
pmguard config set-mode warn

# Silently redirect to the correct PM
pmguard config set-mode redirect

# View current config
pmguard config get
```

Config is stored at `~/.config/pmguard/config.yaml`.

---

## Detection priority

pmguard walks up from your current directory looking for:

1. A lockfile:
   - `bun.lockb` / `bun.lock` → bun
   - `pnpm-lock.yaml` → pnpm
   - `yarn.lock` → yarn
   - `package-lock.json` → npm
2. `packageManager` field in `package.json` (e.g. `"packageManager": "bun@1.0.0"`)

If nothing is found, the command passes through transparently — pmguard stays out of your way.

---

## Arg remapping (redirect mode)

When redirecting, pmguard translates flags between package managers where they differ:

| From | To | Original | Translated |
|------|----|----------|------------|
| pnpm | bun | `-D` | `-d` |
| npm | bun | `--save-dev` | `-d` |
| bun | pnpm | `-d` | `-D` |
| bun | npm | `-d` | `--save-dev` |

---

## Uninstall

Remove the `eval "$(pmguard install-hooks)"` line from your shell config, then:

```bash
# Homebrew
brew uninstall pmguard

# Manual install
rm /usr/local/bin/pmguard
```

---

## Contributing

PRs welcome! Key areas to improve:

- More flag translations in `internal/remap/remap.go`
- Fish shell hook support
- PowerShell hook support
- Per-project config override (`.pmguard.yaml`)

---

## License

MIT