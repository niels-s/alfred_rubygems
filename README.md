Alfred Ruby Gems Workflow
=========================

This repo contains the go script to query rubygems search API. 

Before you build the script yourself make sure you installed the Go SDK. It should work with 1.0 ~ 1.2 because nothing fancy is used. 

To build it use the `go build` command. After building it you can trigger it in the commandline but you need to pass it a search query. Pass it the `-search` flag.

For example:

    ./main -search=rails
    
Feel free to make any pull request and extend for your personal use!
