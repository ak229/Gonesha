package utils

import(
	"strings"
	"../models"
	"github.com/fatih/structs"
	"encoding/json"
)

// the entry point of the framework
func Framework(hostinfo models.HostInfo, clusterInfoJSON string) []models.Result {
	// json decode the host info
	hostinfoJSON,_ :=  json.Marshal(hostinfo)
	moduleResult := "["
	
	// iterate over each module	
	for _,m := range hostinfo.Module {
		// get module information
		modules := GetModule(m)	

		// result of all tests for each module
		testResult := "["

		// iterate over each test to be run for the module
		for _,t := range modules.TestsToRun {
			// get test information
			test := GetTest(t)

			// convert struct to key-value pair
			mt := structs.Map(test)

			// add additional cluster and host information to map
 			mt["HostInfo"] = string(hostinfoJSON)
 			mt["ClusterInfo"] = clusterInfoJSON

			// run the test with the required information
			testResult += ( RunTest(mt,hostinfo) + "," )
		}

		testResult = strings.Trim(testResult,",")
		testResult += "]"

		eval := make(map[string] interface {})
		eval["HostInfo"] = string(hostinfoJSON)
		eval["ClusterInfo"] = clusterInfoJSON
		eval["TestResult"] = testResult
		eval["File"] = modules.EvaluationFile
		eval["Name"] = modules.Name
		moduleResult += ( RunModuleEvaluation(eval) + "," )
	}
	
	moduleResult = strings.Trim(moduleResult,",")
	moduleResult += "]"	

	var result []models.Result
	json.Unmarshal([]byte(moduleResult), &result)

	return result
}
