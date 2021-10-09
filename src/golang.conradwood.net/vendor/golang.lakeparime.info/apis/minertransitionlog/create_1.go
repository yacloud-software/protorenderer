// client create: MinerTransitionLogClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/minertransitionlog/minertransitionlog.proto
   gopackage : golang.lakeparime.info/apis/minertransitionlog
   importname: ai_0
   varname   : client_MinerTransitionLogClient_0
   clientname: MinerTransitionLogClient
   servername: MinerTransitionLogServer
   gscvname  : minertransitionlog.MinerTransitionLog
   lockname  : lock_MinerTransitionLogClient_0
   activename: active_MinerTransitionLogClient_0
*/

package minertransitionlog

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_MinerTransitionLogClient_0 sync.Mutex
  client_MinerTransitionLogClient_0 MinerTransitionLogClient
)

func GetMinerTransitionLogClient() MinerTransitionLogClient { 
    if client_MinerTransitionLogClient_0 != nil {
        return client_MinerTransitionLogClient_0
    }

    lock_MinerTransitionLogClient_0.Lock() 
    if client_MinerTransitionLogClient_0 != nil {
       lock_MinerTransitionLogClient_0.Unlock()
       return client_MinerTransitionLogClient_0
    }

    client_MinerTransitionLogClient_0 = NewMinerTransitionLogClient(client.Connect("minertransitionlog.MinerTransitionLog"))
    lock_MinerTransitionLogClient_0.Unlock()
    return client_MinerTransitionLogClient_0
}

