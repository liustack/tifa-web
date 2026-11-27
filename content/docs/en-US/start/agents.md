---
title: Use Tifa with a coding agent
description: Working on a Tifa app with a coding agent
source: docs/start/agents.md
group: Start
groupOrder: 1
order: 5
---


## Start an app

```sh
go run github.com/liustack/tifa/cmd/tifa@latest new myapp
```

`tifa new` writes two things for coding agents into the app, both from the Tifa version the app uses:

- `AGENTS.md`, notes most coding agents read on their own: the conventions, the commands, and where the docs of this version are. Tifa manages the part between `<!-- tifa:begin -->` and `<!-- tifa:end -->`. Put your own notes outside it.
- The Tifa skill in `.claude/skills/tifa/`, with more detail and links to the docs.

## Which agents read the skill

| Agent | Reads the skill from |
| --- | --- |
| Claude Code | `.claude/skills/` |
| Cursor | `.claude/skills/` (and `.cursor/skills/`) |
| GitHub Copilot | `.claude/skills/` (and `.github/skills/`, `.agents/skills/`) |
| Codex | `.agents/skills/` |
| Gemini CLI | `.agents/skills/` (and `.gemini/skills/`) |

For Codex and Gemini CLI, copy the folder:

```sh
mkdir -p .agents/skills && cp -R .claude/skills/tifa .agents/skills/
```

Copy it again after you upgrade Tifa, since Tifa only writes `.claude/skills/tifa/`.

## Keep them current

`go tool tifa upgrade` writes `AGENTS.md` and the skill again when it moves the app to a newer Tifa. `go tool tifa agents` does it without an upgrade, also in an app that has neither yet. `go tool tifa agents --check` changes nothing and fails when they are out of date. See [`tifa agents`](/docs/build-and-deploy/upgrade#tifa-agents).

## Commands an agent uses

- `go tool tifa docs` prints the docs of the Tifa version the app uses, one line each. Read those, not the docs on GitHub, which may be of another version.
- `go tool tifa dev --background` starts the dev server and returns once it answers. `go tool tifa dev status`, `logs` and `stop` read and end it. See [Agents and scripts](/docs/build-and-deploy/dev-scripts#agents-and-scripts).
- `go tool tifa errors` prints the errors since the last build, those of the browser included. See [Tifa errors](/docs/build-and-deploy/dev-scripts#tifa-errors).
- `go tool tifa check` runs every check before the agent finishes.

The docs on [tifago.com](https://tifago.com) are also there as Markdown: add `.md` to the address of a page, or ask for it with `Accept: text/markdown`. [/llms.txt](https://tifago.com/llms.txt) lists them.
