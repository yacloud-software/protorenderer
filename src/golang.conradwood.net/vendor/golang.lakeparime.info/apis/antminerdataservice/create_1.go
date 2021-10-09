// client create: AntminerDataClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/antminerdataservice/antminerdataservice.proto
   gopackage : golang.lakeparime.info/apis/antminerdataservice
   importname: ai_0
   varname   : client_AntminerDataClient_0
   clientname: AntminerDataClient
   servername: AntminerDataServer
   gscvname  : antminerdataservice.AntminerData
   lockname  : lock_AntminerDataClient_0
   activename: active_AntminerDataClient_0
*/

package antminerdataservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AntminerDataClient_0 sync.Mutex
  client_AntminerDataClient_0 AntminerDataClient
)

func GetAntminerDataClient() AntminerDataClient { 
    if client_AntminerDataClient_0 != nil {
        return client_AntminerDataClient_0
    }

    lock_AntminerDataClient_0.Lock() 
    if client_AntminerDataClient_0 != nil {
       lock_AntminerDataClient_0.Unlock()
       return client_AntminerDataClient_0
    }

    client_AntminerDataClient_0 = NewAntminerDataClient(client.Connect("antminerdataservice.AntminerData"))
    lock_AntminerDataClient_0.Unlock()
    return client_AntminerDataClient_0
}

