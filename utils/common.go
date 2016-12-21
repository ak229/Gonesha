package utils
// Common Utilities
import (
    "bytes"
    "golang.org/x/crypto/ssh"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "../models"
    "strings"
    "github.com/fatih/color"
    "github.com/olekukonko/tablewriter"
    "time"
    "net"
    "sync"
)
const CONFIG_ENV string = "GONESH_HOME"

var red *color.Color
var green *color.Color
var white *color.Color

func init() {

        red = color.New(color.FgRed, color.Bold)
        green = color.New(color.FgGreen, color.Bold)
        white = color.New(color.FgWhite, color.Bold)


}

// get the environment variable
func ShellSpace() string{
	// get the environment variable
	result := os.Getenv(CONFIG_ENV)
	return result
}

// check if the environment variable is set
func CheckShellSpace() {
	// get the workspace directory
	result := ShellSpace()

	// check if the workspace directory is set 
	if len(strings.TrimSpace(result)) == 0 {
		fmt.Println("The Environment Variable " + CONFIG_ENV + " is not set. This Shell needs it buddy...")
		os.Exit(1)	
	}
}

// get the application config
func GetGanpatiConfig() models.GanpatiConfig{
	// read from the input file
        raw, err := ioutil.ReadFile(ShellSpace() + "/config/ganpati.json")
        if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
        }

	// json decoding of the file contents
        var c models.GanpatiConfig
        json.Unmarshal(raw, &c)
        return c
}


// get host information
func GetHostInfo() []models.HostInfo {
	// read from the input file
        raw, err := ioutil.ReadFile(ShellSpace() + "/config/hostinfo.json")
        if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
        }
	
	// json decoding of the file contents
        var c []models.HostInfo
        json.Unmarshal(raw, &c)
        return c
}

// remote login
func RemoteCredentials(user string, password string) *ssh.ClientConfig {
	// initialize login parameters
        config := &ssh.ClientConfig{
                User: user,
                Auth: []ssh.AuthMethod{
                ssh.Password(password),
                },
        }
        return config

}

// start of log entry ( may not be needed )
func Start(host string) {

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = ">>>START"
        if _, err = f.WriteString(info+"\n"); err != nil {
               panic(err)
        }

}

// end of log entry ( may not be needed )
func End(host string) {

        f, err := os.OpenFile("/root/data/"+host+".dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
        var info string = ">>>END"
        if _, err = f.WriteString(info+"\n"); err != nil {
               panic(err)
        }

}

// get all the services ( may not be needed )
func Services( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) string{

        session.Stdout = &stdoutBuf
        session.Run("jps | grep -v Jps | awk '{print \"_\"$2}'")

        var info string = stdoutBuf.String()
        info = strings.Trim(strings.Replace(info,"\n",",",-1),",")
	
	return info

}

// get all the home variables ( may not be needed )
func HomeVariables( variableName string, host string, session *ssh.Session, stdoutBuf bytes.Buffer ) string {
        session.Stdout = &stdoutBuf
        session.Run("echo $"+variableName)

        var info string = stdoutBuf.String()
	return info
}

// get the  path ( may not be needed )
func Path( host string, session *ssh.Session, stdoutBuf bytes.Buffer ) string{
        session.Stdout = &stdoutBuf
        session.Run("echo $PATH")

        var info string = stdoutBuf.String()

	return info
}


// intro
func Intro(){
        ganpatiConfig := GetGanpatiConfig()
        color.Green("----Welcome----")
        color.Green(ganpatiConfig.Name + " : " + ganpatiConfig.Version)
}

// new session
func NewSession(conn *ssh.Client) *ssh.Session{
	session,err := conn.NewSession()
	if err != nil {

		panic(err)
	}
	return session
}

// logging
func Log(level string,host string, message string) {
	var log = "["+host+"]["+level+"] : "+message
	if level == "ERROR" {
		color.Red(log)
	} else if level == "WARNING" {
		color.Yellow(log)
	} else if level == "INFO" {
		color.Cyan(log)
	}
}

// display output in tabular format
func RenderTable( headers []string, data [][]string ) {
	// initialize table
	table := tablewriter.NewWriter(os.Stdout)
	
	// set the headers
	table.SetHeader(headers)
	
	// iterate over the data to be displayed
	for _, v := range data {
    		table.Append(v)
	}

	// render the output
	table.Render()
}

// execute the command on the required host
func Execute( ip string, user string, password string, command string ) string{
	// establish the connection
	conn, err := ssh.Dial("tcp", ip+":22", RemoteCredentials(user, password))
	if err != nil {
        	fmt.Println(err)
	}

	// create a session
	var stdoutBuf bytes.Buffer
	session := NewSession(conn)
	session.Stdout = &stdoutBuf

	// run the command
 	session.Run(command)

	// get the output of the command
 	var info string = stdoutBuf.String()
	return info
}

// check if required services are running ( may not be needed )
func CheckJpsService(name string, host string, conn *ssh.Client) bool{

	var stdoutBuf bytes.Buffer
        var services = strings.TrimSpace(Services(host,NewSession(conn),stdoutBuf))


        result := true
        if len(services) == 0{

		result = false

        } else {
		
		if !strings.Contains("_"+services,name) {

			result = false
		
		}
	}

	return result

}

// color coding of host information
func HostHealth(hostInfo models.HostInfo, hostData models.HostData) {

	red := color.New(color.FgRed, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	white := color.New(color.FgWhite, color.Bold)
	if len(hostData.Ip) == 0 {
		red.Printf(hostInfo.Ip + "\t")
		for _,m := range hostInfo.Module {

			red.Printf(m + "\t")
		}

	} else {

		white.Printf(hostInfo.Ip+ "\t")
		for _,m := range hostData.ModuleDataList {

			if m.Status {
		
				green.Printf(m.Module + "\t")	
		
			} else {


				red.Printf(m.Module + "\t")
			}
			
		}
	}

	fmt.Println("\n\n")
}

// get ip address from reference ( tag or ip address )
func GetHostEntry( alias string) models.HostInfo {
	// get the cluster info
	var result models.HostInfo
	hostInfo := GetHostInfo()

	// iterate over the cluster
	for _,h := range hostInfo {
		// check for a match
        	if h.Ip == alias || h.Tag == alias {
			result = h
         		break
		}
	}
	return result
}

// get the config files
func ConfigFiles() []models.Module {
	// read from the input file
	raw, err := ioutil.ReadFile(ShellSpace() + "/config/modules.json")
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	// json decoding of the file contents
	var c []models.Module	
	json.Unmarshal(raw, &c)
	return c
}

// get host health
func GetReachableHosts(hostinfo []models.HostInfo) []models.HostInfo{

	resultHosts := []models.HostInfo{}

	var wg sync.WaitGroup
 	wg.Add(len(hostinfo))
	
	for _,h := range hostinfo {

		 go func(h models.HostInfo) {

			// check if the machine is up
         		conn_init, err := net.DialTimeout("tcp",h.Ip+":22",time.Duration(5) * time.Second)

			_  = conn_init

         		if err != nil {
				red.Printf("Host "+ h.Ip + " is not reachable ")
                 		fmt.Println()
				wg.Done()
                 		return
			} else {

				// check if login is possible 
				conn, err := ssh.Dial("tcp", h.Ip+":22", RemoteCredentials(h.User, h.Password))
 				_ = conn
				if err != nil {
					red.Printf("Incorrect credentials for Host "+ h.Ip)
         				fmt.Println()
					wg.Done()
         				return
 				}
	
				
				resultHosts = append(resultHosts, h)
			}
         		wg.Done()
		}(h)


	}

	wg.Wait()

	return resultHosts
}
