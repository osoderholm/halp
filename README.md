# halp

Simple POC of CLI helper that uses ChatGPT to suggest commands.

## Build

Requires Go.

Install dependencies:
    
    go mod download

Compile:

    go build -o halp main.go

## Running

Set the `CHATGPT_TOKEN` environment variable before running.

Example:

    ./halp List all files in this directory that have the word go in the filename

You can move the compiled binary to `/usr/local/bin` or somewhere else and add to `PATH` to use from anywhere.
