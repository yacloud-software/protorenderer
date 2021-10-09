// client create: DevOnDutyClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/devonduty/devonduty.proto
   gopackage : golang.lakeparime.info/apis/devonduty
   importname: ai_0
   varname   : client_DevOnDutyClient_0
   clientname: DevOnDutyClient
   servername: DevOnDutyServer
   gscvname  : devonduty.DevOnDuty
   lockname  : lock_DevOnDutyClient_0
   activename: active_DevOnDutyClient_0
*/

package devonduty

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DevOnDutyClient_0 sync.Mutex
  client_DevOnDutyClient_0 DevOnDutyClient
)

func GetDevOnDutyClient() DevOnDutyClient { 
    if client_DevOnDutyClient_0 != nil {
        return client_DevOnDutyClient_0
    }

    lock_DevOnDutyClient_0.Lock() 
    if client_DevOnDutyClient_0 != nil {
       lock_DevOnDutyClient_0.Unlock()
       return client_DevOnDutyClient_0
    }

    client_DevOnDutyClient_0 = NewDevOnDutyClient(client.Connect("devonduty.DevOnDuty"))
    lock_DevOnDutyClient_0.Unlock()
    return client_DevOnDutyClient_0
}

