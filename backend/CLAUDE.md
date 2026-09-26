# Trading Backend
A mono repo which houses microservices and scripts in the cmd.

## Shared Go Guidelines
The code styleguide and git versioning rules are shared with other Go projects and live in goutil (https://github.com/kduong-dev/goutil), expected to be cloned next to trading-core:
@../../goutil/CLAUDE.md

## Backend Additions
- Please write integration tests as well in the `integration-tests` directory at the monorepo root.
- Use the goutil packages (`config`, `fatal`, `httpx`, `eventsource`, `iterator`, ...) rather than re-implementing them. If a helper isn't trading-specific, add it to goutil instead of the outer internal package.
