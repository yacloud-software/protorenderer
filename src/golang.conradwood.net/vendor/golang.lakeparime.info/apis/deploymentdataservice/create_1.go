// client create: DeploymentDataClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/deploymentdataservice/deploymentdataservice.proto
   gopackage : golang.lakeparime.info/apis/deploymentdataservice
   importname: ai_0
   varname   : client_DeploymentDataClient_0
   clientname: DeploymentDataClient
   servername: DeploymentDataServer
   gscvname  : deploymentdataservice.DeploymentData
   lockname  : lock_DeploymentDataClient_0
   activename: active_DeploymentDataClient_0
*/

package deploymentdataservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DeploymentDataClient_0 sync.Mutex
  client_DeploymentDataClient_0 DeploymentDataClient
)

func GetDeploymentDataClient() DeploymentDataClient { 
    if client_DeploymentDataClient_0 != nil {
        return client_DeploymentDataClient_0
    }

    lock_DeploymentDataClient_0.Lock() 
    if client_DeploymentDataClient_0 != nil {
       lock_DeploymentDataClient_0.Unlock()
       return client_DeploymentDataClient_0
    }

    client_DeploymentDataClient_0 = NewDeploymentDataClient(client.Connect("deploymentdataservice.DeploymentData"))
    lock_DeploymentDataClient_0.Unlock()
    return client_DeploymentDataClient_0
}

