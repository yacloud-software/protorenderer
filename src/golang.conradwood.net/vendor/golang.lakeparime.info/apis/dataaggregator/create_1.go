// client create: DataAggregatorClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/dataaggregator/dataaggregator.proto
   gopackage : golang.lakeparime.info/apis/dataaggregator
   importname: ai_0
   varname   : client_DataAggregatorClient_0
   clientname: DataAggregatorClient
   servername: DataAggregatorServer
   gscvname  : dataaggregator.DataAggregator
   lockname  : lock_DataAggregatorClient_0
   activename: active_DataAggregatorClient_0
*/

package dataaggregator

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DataAggregatorClient_0 sync.Mutex
  client_DataAggregatorClient_0 DataAggregatorClient
)

func GetDataAggregatorClient() DataAggregatorClient { 
    if client_DataAggregatorClient_0 != nil {
        return client_DataAggregatorClient_0
    }

    lock_DataAggregatorClient_0.Lock() 
    if client_DataAggregatorClient_0 != nil {
       lock_DataAggregatorClient_0.Unlock()
       return client_DataAggregatorClient_0
    }

    client_DataAggregatorClient_0 = NewDataAggregatorClient(client.Connect("dataaggregator.DataAggregator"))
    lock_DataAggregatorClient_0.Unlock()
    return client_DataAggregatorClient_0
}

