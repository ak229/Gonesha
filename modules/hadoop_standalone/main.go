// Hadoop Analysis
package hadoop_standalone

import(
	"bytes"
	"golang.org/x/crypto/ssh"
	"strings"
	"../../utils"
	"../../models"
)

func RunTests(module string, host string, conn *ssh.Client) models.ModuleData {

	defer conn.Close()

	var output models.ModuleData
	if module == "hadoop_standalone" {

		output = HadoopStandalone(module, host, conn)
	}

	if module == "hadoop_cluster_datanode" {

		//output = HadoopClusterDataNode(host, conn)	

	}

	if module == "hadoop_cluster_namenode" {

		//output = HadoopClusterNameNode(host, conn)

	}

	if module == "hadoop_cluster_secondarynamenode" {

		//output = HadoopClusterSecondaryNameNode(host, conn)

	}

	return output	
}


func HadoopStandalone(module string, host string, conn *ssh.Client) models.ModuleData{

	hadoopHome := CheckHadoopHome(host, conn)
	dataNodeService := utils.CheckJpsService("DataNode", host, conn)
	nameNodeService := utils.CheckJpsService("NameNode", host, conn)
	secondaryNameNodeService := utils.CheckJpsService("SecondaryNameNode", host, conn)
	clusterIDConsistent := CheckClusterID(host, conn)


	var hostData models.ModuleData
	i := 0
	var response = make([]string, 10)
	if !hadoopHome {

		//utils.Log("ERROR", host, "HADOOP_HOME not set")
		response[i] = "Error : HADOOP_HOME not set"
		i++
	}

	if !dataNodeService {

		//utils.Log("ERROR", host, "Data node not running")
		response[i] = "Error : Data node not running"
		i++
	}

	if !nameNodeService {

		//utils.Log("ERROR", host, "Name node not running")
		response[i] = "Error : Name node not running"
		i++
	}

	if !secondaryNameNodeService {

		//utils.Log("ERROR", host, "Secondary name node not running")
		response[i] = "Error : Secondary name node not running"
		i++
	}

	if !clusterIDConsistent {

		//utils.Log("ERROR", host, "Cluster id mismatch")
		response[i] = "Error : Cluster id mismatch"
		i++
	}

	result := hadoopHome && dataNodeService && nameNodeService && secondaryNameNodeService && clusterIDConsistent

	hostData.Status = result
	hostData.Response = response
	hostData.Module = module

	return hostData
}


func HadoopClusterDataNode(host string, conn *ssh.Client) string{


	return "datanode"
}

func HadoopClusterNameNode(host string, conn *ssh.Client) string{


	return "namenode"
}


func HadoopClusterSecondaryNameNode(host string, conn *ssh.Client) string{


	return "secondarynamenode"
}

func CheckHadoopHome(host string, conn *ssh.Client) bool{

	var stdoutBuf bytes.Buffer
	homeVariable := utils.HomeVariables("HADOOP_HOME",host,utils.NewSession(conn),stdoutBuf)

	//utils.Log("INFO", host, "Checking for HADOOP_HOME")

	result := true	
	if len(strings.TrimSpace(homeVariable)) == 0{
		
		result = false	
	
	}
	
	return result

}

func CheckClusterID(host string, conn *ssh.Client) bool{

	var clusterIDString = strings.TrimSpace(ClusterID(host,utils.NewSession(conn)))

	//utils.Log("INFO", host, "Checking Cluster Compatibility")
	
	result := true
	if len(clusterIDString) == 0 {
		
		result = false	
	
	} else {
		var clusterID = make([] string,2)
		clusterID = strings.Split(clusterIDString,",")
	
		// cluster id mismatch
		if len(clusterID) == 2 {
			if clusterID[0] != clusterID[1] {

				result = false
			}
		}
	}

	return result	
}




func ClusterID( host string, session *ssh.Session) string{

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run("cat /home/bigdata/hadoop/{data,name}/current/VERSION | grep clusterID | awk '{split($$0,c,\"=\");print c[2]}'")
        
	var info string = stdoutBuf.String()
        info = strings.Trim(strings.Replace(info,"\n",",",-1),",")

	return info
}
