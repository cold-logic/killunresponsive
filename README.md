# Kill Unresponsive

Terminal command to terminate unresponsive processes.

## How it works

This command generates a spin dump report using Apple's provided "spindump" command. The report generated contains samples of user and kernel stacks for every process in the system. This command parses that report to obtain a list of unresponsive processes and then terminates them.

## How to compile

Simply run the `go build` command from the terminal while browsing this directory.

## How to install

I copied the compiled app into my /usr/local/sbin directory but you're welcome to install it anywhere you like.