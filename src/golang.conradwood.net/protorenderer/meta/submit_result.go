package meta

import (
	"fmt"
	"time"
	//	"golang.conradwood.net/apis/create"
	pr "golang.conradwood.net/apis/protorenderer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cache"
	"golang.conradwood.net/protorenderer/db"
	"sync"
)

var (
	pkgidgen  *idGen
	svcidgen  *idGen
	rpcidgen  *idGen
	msgidgen  *idGen
	persistid *db.DBPersistID
)

func init() {
	pkgidgen = newidgen("P")
	svcidgen = newidgen("S")
	rpcidgen = newidgen("R")
	msgidgen = newidgen("M")
}
func newidgen(t string) *idGen {
	s := fmt.Sprintf("id_%s_cache", t)
	res := &idGen{t: t,
		pkgids: cache.NewResolvingCache(s, time.Duration(999999)*time.Hour, 9999999),
	}
	return res
}

type idGen struct {
	t       string
	pkgids  cache.CachingResolver
	pkglock sync.Mutex
}

func (id *idGen) Retrieve(key string) (string, error) {
	v, err := id.pkgids.Retrieve(key, func(k string) (interface{}, error) {
		dbkey := fmt.Sprintf("%s_%s", id.t, key)
		id.pkglock.Lock()
		defer id.pkglock.Unlock()
		ctx := authremote.Context()
		pids, err := persistid.ByKey(ctx, dbkey)
		if err != nil {
			return nil, err
		}
		if len(pids) > 0 {
			return fmt.Sprintf("%d", pids[0].ID), nil
		}
		newid := &pr.PersistID{Key: dbkey}
		id, err := persistid.Save(ctx, newid)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%d", id), nil
	})
	if err != nil {
		return "", err
	}
	s := v.(string)
	return s, nil
}

// we 'submit' (aka remember result), but we also make sure
// that each object has a unique ID
func (m *MetaCompiler) submitResult(result *Result) error {
	if persistid == nil {
		persistid = db.DefaultDBPersistID()
	}
	for _, pkg := range result.Packages {
		// all packages need an id
		key := fmt.Sprintf("%s/%s/%s", pkg.FQDN, pkg.Name, pkg.Filename)
		v, err := pkgidgen.Retrieve(key)
		if err != nil {
			return err
		}
		pkg.Proto.ID = v

		// all messages need an ID
		for _, m := range pkg.Messages {
			if m.ID != "" {
				continue
			}
			key := fmt.Sprintf("%s_%s", pkg.Proto.ID, m.Name)
			mid, err := msgidgen.Retrieve(key)
			if err != nil {
				return err
			}
			m.ID = mid
		}

		// all services and RPCs need an ID
		for _, s := range pkg.Services {
			key := fmt.Sprintf("%s_%s", pkg.Proto.ID, s.Name)
			sid, err := svcidgen.Retrieve(key)
			if err != nil {
				return err
			}
			s.ID = sid
			for _, r := range s.RPCs {
				if r.ID != "" {
					continue
				}
				key := fmt.Sprintf("%s_%s_%s", pkg.Proto.ID, s.ID, r.Name)
				rid, err := rpcidgen.Retrieve(key)
				if err != nil {
					return err
				}
				r.ID = rid
			}
		}

	}

	// set the packageid in the meta set of each protofile
	fmt.Printf("Meta: Meta compiler has compiled %d packages\n", len(result.Packages))
	for _, pkg := range result.Packages {
		if pkg.Proto.ID == "" {
			fmt.Printf("Meta: BUG package %s has no ID!\n", pkg.FQDN)
			continue
		}
		for _, pf := range pkg.Protofiles {
			if pf.Meta == nil {
				pf.Meta = &pr.MetaProtoFile{}
			}
			pf.Meta.PackageID = pkg.Proto.ID
			pf.Meta.Package = pkg.Proto
			//			fmt.Printf("Meta: setting packageid to \"%s\" for protofile %s\n", pkg.Proto.ID, pf.Filename)
		}

	}
	return nil
}

func (m *MetaCompiler) GetMostRecentResult() *Result {
	return m.result
}



































































































