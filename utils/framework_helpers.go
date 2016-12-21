package utils

import (
	"../models"
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
)

// get test information
func GetTest( name string ) models.Test {
	// read from the input file
	raw, err := ioutil.ReadFile(ShellSpace() + "/config/tests.json")
	if err != nil {
 		fmt.Println(err.Error())
 		os.Exit(1)
 	}

	// json decoding of the file contents	
	var test []models.Test
 	json.Unmarshal(raw, &test)

	// iterate over each object
	// terminate if a match is found
	var result models.Test
	for _, t := range test {
		
		if t.Name == name {
			result = t
			break
		}
	}
	
	return result
}

// get module information
func GetModule( name string ) models.Module {
	// read from the input file
	raw, err := ioutil.ReadFile(ShellSpace() + "/config/modules.json")
 	if err != nil {
        	fmt.Println(err.Error())
        	os.Exit(1)
 	}

	// json decoding of the file contents
	var modules []models.Module
	json.Unmarshal(raw, &modules)

	// iterate over each object
	// terminate if a match is found
	var result models.Module
 	for _, m := range modules {
        	if m.Name == name {
                	result = m
			break
         	}
	}

	return result

}

// run test
func RunTest( params map[string] interface {}, hostinfo models.HostInfo) string {
	// data required for the test
 	GetMethod,_ := params["GetMethod"].(string)
 	command := RunJS(params,"tests",GetMethod)
	output := Execute(&hostinfo.Ip, &hostinfo.User, &hostinfo.Password, &command)

	// run the check on the data obtained from the host	
	params["data"] = output
	CheckMethod,_ := params["CheckMethod"].(string)

	return RunJS(params,"tests", CheckMethod)

}

// run evaluation
func RunModuleEvaluation(params map[string] interface {}) string {
	
	return RunJS(params,"role", "result")
}


// render table
func RenderRecord(Ip string, tag string, result []models.Result) {

	white.Printf(Ip + "\t" + tag + "\t")
	
	// iterate over each module
	for _,r := range result {

		if r.Status {
	
			green.Printf(r.Name)		
	
		} else {

			red.Printf(r.Name)
		}	

	}
	fmt.Println()
}
