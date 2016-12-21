package utils

import (
   	"github.com/robertkrimen/otto"
	"io/ioutil"
	"fmt"
)

// run javascript file
func RunJS(params map[string]interface {}, scriptType string, method string) string{	
	// get the path	
	path, _ := params["File"].(string)

	// read from the input file
	raw, err := ioutil.ReadFile(ShellSpace() + "/workspace/" + scriptType + "/" + path)
	if err != nil {
        	fmt.Println("File : " + ShellSpace() + "/workspace/" + scriptType + "/" + path + " not found ")
		fmt.Println("Leaving...")
		return "failed"
 	}

	// create new instance of javascript engine
	vm := otto.New()
 	
	// set the parameters to be passed to the engine
 	for k,v := range params {
		vs,_ := v.(string)
		vm.Set(k,vs)
 	}

	// run the contents of the file
	vm.Run(raw)

	// get the required result
 	rawResult,err := vm.Get(method)
	if err != nil {
		fmt.Println(" The give javascript does not have a result ")
		panic(err)
	}
	result,err := rawResult.ToString()
	if err != nil {
		
		panic(err)

	}
	return result
}
