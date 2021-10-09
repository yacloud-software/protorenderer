// client create: AntminerapiClient
/* geninfo:
   filename  : golang.lakeparime.com/apis/antminerapi/antminerapi.proto
   gopackage : golang.lakeparime.com/apis/antminerapi
   importname: ai_0
   varname   : client_AntminerapiClient_0
   clientname: AntminerapiClient
   servername: AntminerapiServer
   gscvname  : antminerapi.Antminerapi
   lockname  : lock_AntminerapiClient_0
   activename: active_AntminerapiClient_0
*/

package antminerapi

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AntminerapiClient_0 sync.Mutex
  client_AntminerapiClient_0 AntminerapiClient
)

func GetAntminerapiClient() AntminerapiClient { 
    if client_AntminerapiClient_0 != nil {
        return client_AntminerapiClient_0
    }

    lock_AntminerapiClient_0.Lock() 
    if client_AntminerapiClient_0 != nil {
       lock_AntminerapiClient_0.Unlock()
       return client_AntminerapiClient_0
    }

    client_AntminerapiClient_0 = NewAntminerapiClient(client.Connect("antminerapi.Antminerapi"))
    lock_AntminerapiClient_0.Unlock()
    return client_AntminerapiClient_0
}

