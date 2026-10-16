# Feature requests

Future features, not near-term work. (Basic build steps go in [TODO.md](TODO.md).)

## Multiple daemon instances on one machine

For now there is exactly one daemon per machine. Later, handle every case
where someone starts more than one:
- detecting a daemon that's already running and refusing a second instance
- concurrent `localcloud ... start` races, where two CLIs both spawn a daemon
- the single hardcoded socket path `/tmp/daemon.sock` colliding

## One machine as both master and worker

Let a single machine run as a master and a worker at the same time, instead of
one role per machine.
