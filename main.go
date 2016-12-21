package main

import (
        "github.com/abiosoft/ishell"
        "./shellwrapper"
        "./utils"
)

func main() {

	// run pre-requisite check
        utils.CheckShellSpace()

	// create a new shell instance
        shell := ishell.New()

	// set intial parameters
        shell.SetPrompt(utils.PROMPT + utils.TERMINATOR)
        shell.Println(utils.NAME)
        shell.Println(utils.VERSION)

	// intialize commands
	shell.Register("cluster",shellwrapper.Cluster)
        shell.Register("run",shellwrapper.Run)
        shell.Register("status",shellwrapper.Status)
        shell.Register("status",shellwrapper.Status)
        shell.Register("config", shellwrapper.Config)
        shell.Register("fetchinfo",shellwrapper.FetchInfo)
        
	// start the shell session
	shell.Start()
}
