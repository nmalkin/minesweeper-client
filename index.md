---
layout: default
title: Minesweeper Client
---

Minesweeper Client
==================

A client to help you [play Minesweeper](http://minesweeper.nm.io)

## How to use it

The runner is meant to save you time by implementing the interaction with
Minesweeper's HTTP API. Instead, you just have to read from standard input
and print to standard output (which is easier in most languages).

Your program should start a game by reading from standard input. It will get a
line of text that is actually three numbers separated by commas: the dimensions
of the board (width, then height) and the number of mines on it.

At this point, your AI should start opening cells. To open a cell, print a new
line with its comma-separated x,y coordinate. For example, to open the
top-left cell of the board, print `0,0`. (Note that, in some languages, you
may need to add a new line or flush your buffer for the runner to receive
your output).

After making your guess, read from standard input again. You'll get one of the
following inputs:

- `lost` if you tried to open a cell with a mine. Sorry about that!
- a number, specifying how many mines neighbor the opened cell
- `win` if you opened the last mine-free cell. Congratulations!
- `Bad Request` if something went wrong. Make sure that your outputs are
  formatted exactly as expected, and don't hesitate to ask us for help.

## How to get it

[Visit the downloads page](downloads.html) and get the binary for your
architecture.

#### I don't trust your binaries!

Clone [this repository](https://github.com/nmalkin/minesweeper-client), make sure you have
[Go installed](http://golang.org/doc/install), then run `go build` from inside
the repository.

## How to run it

Run the helper program, followed by the name of your AI's executable as you
want it executed. (Your command can have many parts.)

    ./minesweeper ./my-ai --win=always

Actually, that won't work exactly: the client executable has one required flag: 
your name, to distinguish your submission from everybody else's. 
You use it like this:

    ./minesweeper -name="Vercingetorix" ./my-ai --win=always

There are also some optional flags:

- `games` to specify how many games to play back-to-back
- `version` the version number of your AI (for scoring)
