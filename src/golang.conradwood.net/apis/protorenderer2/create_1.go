// client create: ProtoRenderer2Client
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/protorenderer2/protorenderer2.proto
   gopackage : golang.conradwood.net/apis/protorenderer2
   importname: ai_0
   clientfunc: GetProtoRenderer2
   serverfunc: NewProtoRenderer2
   lookupfunc: ProtoRenderer2LookupID
   varname   : client_ProtoRenderer2Client_0
   clientname: ProtoRenderer2Client
   servername: ProtoRenderer2Server
   gsvcname  : protorenderer2.ProtoRenderer2
   lockname  : lock_ProtoRenderer2Client_0
   activename: active_ProtoRenderer2Client_0
*/

package protorenderer2

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ProtoRenderer2Client_0 sync.Mutex
  client_ProtoRenderer2Client_0 ProtoRenderer2Client
)

func GetProtoRenderer2Client() ProtoRenderer2Client { 
    if client_ProtoRenderer2Client_0 != nil {
        return client_ProtoRenderer2Client_0
    }

    lock_ProtoRenderer2Client_0.Lock() 
    if client_ProtoRenderer2Client_0 != nil {
       lock_ProtoRenderer2Client_0.Unlock()
       return client_ProtoRenderer2Client_0
    }

    client_ProtoRenderer2Client_0 = NewProtoRenderer2Client(client.Connect(ProtoRenderer2LookupID()))
    lock_ProtoRenderer2Client_0.Unlock()
    return client_ProtoRenderer2Client_0
}

func ProtoRenderer2LookupID() string { return "protorenderer2.ProtoRenderer2" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("protorenderer2.ProtoRenderer2")
   AddService("protorenderer2.ProtoRenderer2")
}
