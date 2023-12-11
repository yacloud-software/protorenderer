// client create: ProtoRendererServiceClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/protorenderer/protorenderer.proto
   gopackage : golang.conradwood.net/apis/protorenderer
   importname: ai_0
   clientfunc: GetProtoRendererService
   serverfunc: NewProtoRendererService
   lookupfunc: ProtoRendererServiceLookupID
   varname   : client_ProtoRendererServiceClient_0
   clientname: ProtoRendererServiceClient
   servername: ProtoRendererServiceServer
   gsvcname  : protorenderer.ProtoRendererService
   lockname  : lock_ProtoRendererServiceClient_0
   activename: active_ProtoRendererServiceClient_0
*/

package protorenderer

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ProtoRendererServiceClient_0 sync.Mutex
  client_ProtoRendererServiceClient_0 ProtoRendererServiceClient
)

func GetProtoRendererClient() ProtoRendererServiceClient { 
    if client_ProtoRendererServiceClient_0 != nil {
        return client_ProtoRendererServiceClient_0
    }

    lock_ProtoRendererServiceClient_0.Lock() 
    if client_ProtoRendererServiceClient_0 != nil {
       lock_ProtoRendererServiceClient_0.Unlock()
       return client_ProtoRendererServiceClient_0
    }

    client_ProtoRendererServiceClient_0 = NewProtoRendererServiceClient(client.Connect(ProtoRendererServiceLookupID()))
    lock_ProtoRendererServiceClient_0.Unlock()
    return client_ProtoRendererServiceClient_0
}

func GetProtoRendererServiceClient() ProtoRendererServiceClient { 
    if client_ProtoRendererServiceClient_0 != nil {
        return client_ProtoRendererServiceClient_0
    }

    lock_ProtoRendererServiceClient_0.Lock() 
    if client_ProtoRendererServiceClient_0 != nil {
       lock_ProtoRendererServiceClient_0.Unlock()
       return client_ProtoRendererServiceClient_0
    }

    client_ProtoRendererServiceClient_0 = NewProtoRendererServiceClient(client.Connect(ProtoRendererServiceLookupID()))
    lock_ProtoRendererServiceClient_0.Unlock()
    return client_ProtoRendererServiceClient_0
}

func ProtoRendererServiceLookupID() string { return "protorenderer.ProtoRendererService" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("protorenderer.ProtoRendererService")
}

































