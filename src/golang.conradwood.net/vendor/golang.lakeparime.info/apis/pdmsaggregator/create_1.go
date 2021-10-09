// client create: PDMSAggregatorClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/pdmsaggregator/pdmsaggregator.proto
   gopackage : golang.lakeparime.info/apis/pdmsaggregator
   importname: ai_0
   varname   : client_PDMSAggregatorClient_0
   clientname: PDMSAggregatorClient
   servername: PDMSAggregatorServer
   gscvname  : pdmsaggregator.PDMSAggregator
   lockname  : lock_PDMSAggregatorClient_0
   activename: active_PDMSAggregatorClient_0
*/

package pdmsaggregator

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PDMSAggregatorClient_0 sync.Mutex
  client_PDMSAggregatorClient_0 PDMSAggregatorClient
)

func GetPDMSAggregatorClient() PDMSAggregatorClient { 
    if client_PDMSAggregatorClient_0 != nil {
        return client_PDMSAggregatorClient_0
    }

    lock_PDMSAggregatorClient_0.Lock() 
    if client_PDMSAggregatorClient_0 != nil {
       lock_PDMSAggregatorClient_0.Unlock()
       return client_PDMSAggregatorClient_0
    }

    client_PDMSAggregatorClient_0 = NewPDMSAggregatorClient(client.Connect("pdmsaggregator.PDMSAggregator"))
    lock_PDMSAggregatorClient_0.Unlock()
    return client_PDMSAggregatorClient_0
}

