// client create: CSVFileServerClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/csvfileserver/csvfileserver.proto
   gopackage : golang.lakeparime.info/apis/csvfileserver
   importname: ai_0
   varname   : client_CSVFileServerClient_0
   clientname: CSVFileServerClient
   servername: CSVFileServerServer
   gscvname  : csvfileserver.CSVFileServer
   lockname  : lock_CSVFileServerClient_0
   activename: active_CSVFileServerClient_0
*/

package csvfileserver

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CSVFileServerClient_0 sync.Mutex
  client_CSVFileServerClient_0 CSVFileServerClient
)

func GetCSVFileServerClient() CSVFileServerClient { 
    if client_CSVFileServerClient_0 != nil {
        return client_CSVFileServerClient_0
    }

    lock_CSVFileServerClient_0.Lock() 
    if client_CSVFileServerClient_0 != nil {
       lock_CSVFileServerClient_0.Unlock()
       return client_CSVFileServerClient_0
    }

    client_CSVFileServerClient_0 = NewCSVFileServerClient(client.Connect("csvfileserver.CSVFileServer"))
    lock_CSVFileServerClient_0.Unlock()
    return client_CSVFileServerClient_0
}

