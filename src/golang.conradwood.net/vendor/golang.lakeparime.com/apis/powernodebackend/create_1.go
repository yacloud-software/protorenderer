// client create: PowerNodeBackEndClient
/* geninfo:
   filename  : golang.lakeparime.com/apis/powernodebackend/powernodebackend.proto
   gopackage : golang.lakeparime.com/apis/powernodebackend
   importname: ai_0
   varname   : client_PowerNodeBackEndClient_0
   clientname: PowerNodeBackEndClient
   servername: PowerNodeBackEndServer
   gscvname  : powernodebackend.PowerNodeBackEnd
   lockname  : lock_PowerNodeBackEndClient_0
   activename: active_PowerNodeBackEndClient_0
*/

package powernodebackend

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PowerNodeBackEndClient_0 sync.Mutex
  client_PowerNodeBackEndClient_0 PowerNodeBackEndClient
)

func GetPowerNodeBackEndClient() PowerNodeBackEndClient { 
    if client_PowerNodeBackEndClient_0 != nil {
        return client_PowerNodeBackEndClient_0
    }

    lock_PowerNodeBackEndClient_0.Lock() 
    if client_PowerNodeBackEndClient_0 != nil {
       lock_PowerNodeBackEndClient_0.Unlock()
       return client_PowerNodeBackEndClient_0
    }

    client_PowerNodeBackEndClient_0 = NewPowerNodeBackEndClient(client.Connect("powernodebackend.PowerNodeBackEnd"))
    lock_PowerNodeBackEndClient_0.Unlock()
    return client_PowerNodeBackEndClient_0
}

