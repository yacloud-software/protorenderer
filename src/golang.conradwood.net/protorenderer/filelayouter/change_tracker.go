package filelayouter

import (
	pr "golang.conradwood.net/apis/protorenderer"
	pc "golang.conradwood.net/protorenderer/protocache"
)

/*
* for each "version" we list which files changed
 */
type trackedVersion struct {
	version      int
	previous     *trackedVersion
	fileVersions []*TrackedChange
}

func (tv *trackedVersion) Clone() *trackedVersion {
	pl := &trackedVersion{previous: tv}
	if tv == nil {
		return pl
	}
	pl.version = 0
	return pl

}

// accumulated changes since a given version
func (tv *trackedVersion) ChangesSince(version int) []*TrackedChange {
	if tv.version == version {
		// there is no different between me and myself
		return nil
	}
	var changes []*TrackedChange
	lv := tv.getVersion(version)
	if lv == nil {
		return tv.fileVersions
	}
	for _, tc := range tv.fileVersions {
		if lv.IsDifferent(tc) {
			changes = append(changes, tc)
		}
	}
	return changes

}

// get all previous versions higher and including a number
func (tv *trackedVersion) getVersion(version int) *trackedVersion {
	if tv == nil {
		return nil
	}
	if tv.version == version {
		return tv
	}
	if tv.previous != nil {
		return tv.previous.getVersion(version)
	}
	return nil
}

// return changes since last version
func (tv *trackedVersion) Changes() []*TrackedChange {
	if tv.previous == nil {
		// no previous - everything changed!
		return tv.fileVersions
	}
	var res []*TrackedChange
	for _, tc := range tv.fileVersions {
		if tv.previous.IsDifferent(tc) {
			res = append(res, tc)
		}
	}
	return res
}

// returns true if the tracked change is different than what this tracked version "remembers"
func (tv *trackedVersion) IsDifferent(tc *TrackedChange) bool {
	for _, vtc := range tv.fileVersions {
		if vtc.cf.ProtoFile() != tc.cf.ProtoFile() {
			continue
		}
		if vtc.cf.Version() == tc.cf.Version() {
			return false
		}
		return true
	}
	return true
}

// submit the current version of file. tracked version will mark it as changed *if* it has changed
func (tv *trackedVersion) CurrentFile(cf *pc.CachedFile) *TrackedChange {
	tc := &TrackedChange{cf: cf}
	tv.fileVersions = append(tv.fileVersions, tc)
	return tc
}

type TrackedChange struct {
	cf          *pc.CachedFile
	filename    string
	prefix      string // e.g. "golang.conradwood.net/apis"
	relativeDir string // e.g. "golang.conradwood.net/apis/common"
}

func (tc *TrackedChange) Protofile() *pr.ProtoFile {
	return tc.cf.ProtoFile()
}
func (tc *TrackedChange) Filename() string {
	return tc.filename
}
func (tc *TrackedChange) RelativeDir() string {
	return tc.relativeDir
}























































































