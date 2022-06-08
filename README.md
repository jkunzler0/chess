# P2P Chess w/ Score Database

WIP

TODO:
- Stop other peers from connecting to ongoing game
- Stop peers from crashing on game completion
- Update/streamline game tests
- Add a test for p2p?
- Database core functionality
- Validate p2p game results for the database

Misc. Questions to Anwser:
- Where should I put defer close statements?
- Should P2pGame params be a single struct?
- Is the global *bufio.Reader in game.go appropriate?
- 