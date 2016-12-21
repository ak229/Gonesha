package shellwrapper

import (
	"strings"
	"golang.org/x/crypto/ssh"
	"sync"
	"../modules/hadoop_standalone"
	"../utils"
	"net"
	"time"
	"../models"
	"github.com/fatih/color"
	"fmt"
	"../help"
	"encoding/json"
)

// get information regarding the cluster
func Cluster (args ...string) (string,error) {
	// get the host info
	hostInfo := utils.GetHostInfo()
	cluster := make([][] string,200)

	// iterate over each host
	i := 0
	for _,h := range hostInfo {

        	cluster[i] = [] string{h.Ip, h.Tag, strings.Join(h.Module,",")}

	i++
	}

	// render the output in the form of a table
	utils.RenderTable([]string{"Host","Tag","Module"}, cluster)

	return "",nil
}


// run a command on a remote machine
func Run (args ...string) (string,error) {
	// check for help
	if help.Needed(args) {
        	fmt.Println(utils.RUN)
        	return "", nil
	}

	// get the required parameters
	reference := args[0]
	command := args[1]
	var info = ""

	// get the ip address if host tag is passed
	h := utils.GetHostEntry(reference)

	// run the command on the host
	info = utils.Execute(h.Ip,h.User,h.Password,command)

	return info,nil
}

 
// get the status of the cluster ( this may not be needed eventually )
func Status (args ...string) (string,error) {
	// check for help	
	if help.Needed(args) {
        	fmt.Println(utils.STATUS)
        	return "", nil
	}

	white := color.New(color.FgWhite,color.Bold)
	white.Println("\t\t\t\t*****Cluster Status*****")
	hostInfo := utils.GetHostInfo()
 	// thread initialization
 	var wg sync.WaitGroup
 	wg.Add(len(hostInfo))
	i := 0
	hostData := make([] models.HostData,len(hostInfo))
	for _,h := range hostInfo {

		go func(i int,module []string,ip string,user string, password string) {
                 
			conn_init, err := net.DialTimeout("tcp",ip+":22",time.Duration(5) * time.Second)
			if err != nil {
				wg.Done()
				return	
			}

			if conn_init == conn_init {

			}
			
			conn, err := ssh.Dial("tcp", ip+":22", utils.RemoteCredentials(user, password))
			if err != nil {
                 		wg.Done()
				return
			}
                 	
			defer wg.Done()
                 	
			moduleData := []models.ModuleData{} 
 
			for _,m := range module {
				
				moduleData = append(moduleData, resolve(m, ip, conn))
                 	
			}
		
			hostData[i] = models.HostData{Ip: ip, ModuleDataList : moduleData}

         	}(i, h.Module, h.Ip, h.User, h.Password)
	i++
	}	

	wg.Wait()
	
	i = 0
	for _,h := range hostInfo {

		utils.HostHealth(h, hostData[i])
	i++	
	}

	return "",nil


}

// gets the configuration information corresponding to each module
func Config(args ...string) (string,error) {
	result := "config accepts exactly 2 arguments"	
	if len(args) > 0 && len(args) < 3 {
		if len(args) == 1 {
			if args[0] == "files" {
				result = ""
				moduleConfig := []models.Module{} 
				moduleConfig = utils.ConfigFiles()
				headers := []string{}	
				data := [][]string{[]string{}}
				for _,mc := range moduleConfig {
					fmt.Println("\n\t\tModule : "+ mc.Name)
					headers = []string{"Tag","Path"}
					data = [][]string{[]string{}}
					for _,f := range mc.FileInfo {
						data = append(data,[]string{ f.Tag, f.Path })	
					}
					utils.RenderTable(headers,data)
					fmt.Println()
				}
			
			} else {

			}
		}

	}
	return result,nil
}



// Mapper ( may not be needed later )
func resolve(module string,host string, conn *ssh.Client) models.ModuleData{

	var output models.ModuleData
        if module == "hadoop_standalone" {
                output = hadoop_standalone.RunTests(module, host, conn) 
        } else if module == "hadoop_cluster" {

        }

	return output
}



func FetchInfo(args ...string) (string,error) {
	// get cluster info
	hostInfo := utils.GetHostInfo()
 
	// json decode cluster info
	hostInfoJSON,_ :=  json.Marshal(hostInfo)
	
 	
	var result []models.Result

	// get subset of hosts that are reachable
	resultHosts := utils.GetReachableHosts(hostInfo)

	// thread initialization
	var wg sync.WaitGroup
	wg.Add(len(resultHosts))

	// iterate over the hosts
 	for _,h := range resultHosts {

		go func(h models.HostInfo) {
			
			// call the framework
			result = utils.Framework(h, string(hostInfoJSON))	
			utils.RenderRecord(h.Ip, h.Tag, result)
			defer wg.Done()
		
		}(h)

	}
	wg.Wait()
		return "",nil

}

