[![Build Status](https://drone.monkiato.com/api/badges/monkiato/word-challenge/status.svg?ref=refs/heads/master)](https://drone.monkiato.com/monkiato/word-challenge)
[![codecov](https://codecov.io/gh/monkiato/word-challenge/branch/master/graph/badge.svg?token)](https://codecov.io/gh/monkiato/word-challenge)


# Word Challenge Game

This is a server in Go for a real-time multiplayer game.

The game consists in a challenge where users have to write words as fast
as they can in order to obtain score based on the speed.

The first one to write the word correctly get the scores and a new word
will appear for everyone.

## Build Docker Image

`docker build . -t word-challenge`

default port: 8080

## TODO

 - Remove client webpage and run a different server for this
 - Configurable port
 - change challenge type to renew words every X time
 
 