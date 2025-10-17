# better-fg rules/markdown splicing tool mono-repo

This is a more userfriendly and interactive version of `fg`. If the user runs it with only one background job it should work as fg normally would but if it has more than one background job the user should be presented with a list to select from. It should also be possible to fuzzy search for the job in the list. If the user presses enter straight away it should act as `fg` would when it has multiple jobs available but the user didn't specify a specific one. Additionally we want to both support `fg %<number>` but also `fg <number>`.

## Tech
- Language: GO
- Frameworks: charm.sh (bubble tea, huh, lip gloss, bubbles, log)
- git
- Nix (for building and managing dev-shell, uses direnv)
- `direnv exec` for making sure we always execute our commands in the correct dev-shell.

## Best practices
- Always make commits after making changes
- Write simple commit messages prefixed using feat/chore/bug
- Before starting work, define the task and add it to TASKS.md
- Always add documentation for all cli arguments, config flags, etc.
- This tool will support many other cli tools with their different names and oddities for these types of files, we should include this in our code architecture such that there can be support for multiple "output formats"
- Always add unit tests for all functionality you add, aim for at least full branch coverage
- Write integration tests that test the full functionality, make sure XDG_CONFIG_HOME is overridden for the test execution such that the test-fragments can be part of the repository
- Always run linters with `--fix` or similar flags to automatically use safe fixes

## Git
- Always start new features in a feature branch `feat/<some-name>`
- Make many small commits

## NEVER DO
- NEVER EVER ADD CODE ATTRIBUTIONS IN COMMIT DESCRIPTIONS THAT REFERENCES THE CLI TOOL!

