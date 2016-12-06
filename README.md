#Installation

Install this littul utility with:

`go get -u github.com/RobMeades/udp-echo-server`

There is command-line help.  An example command line that leaves the echo server running on port 1000 might be:

`nohup udp-echo-server -p 1000 > udp.log &`