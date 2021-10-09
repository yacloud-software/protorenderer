// client create: OrkneyScriptSupportClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/orkneybetascripts/orkneybetascripts.proto
   gopackage : golang.lakeparime.info/apis/orkneybetascripts
   importname: ai_0
   varname   : client_OrkneyScriptSupportClient_0
   clientname: OrkneyScriptSupportClient
   servername: OrkneyScriptSupportServer
   gscvname  : orkneybetascripts.OrkneyScriptSupport
   lockname  : lock_OrkneyScriptSupportClient_0
   activename: active_OrkneyScriptSupportClient_0
*/

package orkneybetascripts

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_OrkneyScriptSupportClient_0 sync.Mutex
  client_OrkneyScriptSupportClient_0 OrkneyScriptSupportClient
)

func GetOrkneyScriptSupportClient() OrkneyScriptSupportClient { 
    if client_OrkneyScriptSupportClient_0 != nil {
        return client_OrkneyScriptSupportClient_0
    }

    lock_OrkneyScriptSupportClient_0.Lock() 
    if client_OrkneyScriptSupportClient_0 != nil {
       lock_OrkneyScriptSupportClient_0.Unlock()
       return client_OrkneyScriptSupportClient_0
    }

    client_OrkneyScriptSupportClient_0 = NewOrkneyScriptSupportClient(client.Connect("orkneybetascripts.OrkneyScriptSupport"))
    lock_OrkneyScriptSupportClient_0.Unlock()
    return client_OrkneyScriptSupportClient_0
}

