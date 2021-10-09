// client create: UserMetaClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/usermeta/usermeta.proto
   gopackage : golang.lakeparime.info/apis/usermeta
   importname: ai_0
   varname   : client_UserMetaClient_0
   clientname: UserMetaClient
   servername: UserMetaServer
   gscvname  : usermeta.UserMeta
   lockname  : lock_UserMetaClient_0
   activename: active_UserMetaClient_0
*/

package usermeta

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_UserMetaClient_0 sync.Mutex
  client_UserMetaClient_0 UserMetaClient
)

func GetUserMetaClient() UserMetaClient { 
    if client_UserMetaClient_0 != nil {
        return client_UserMetaClient_0
    }

    lock_UserMetaClient_0.Lock() 
    if client_UserMetaClient_0 != nil {
       lock_UserMetaClient_0.Unlock()
       return client_UserMetaClient_0
    }

    client_UserMetaClient_0 = NewUserMetaClient(client.Connect("usermeta.UserMeta"))
    lock_UserMetaClient_0.Unlock()
    return client_UserMetaClient_0
}

