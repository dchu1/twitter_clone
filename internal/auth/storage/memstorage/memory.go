package memstorage

import "sync"

var UsersCredRWmu = sync.RWMutex{}
var SessionManagerRWmu = sync.RWMutex{}
var UsersCred = make(map[string]string)
var SessionManager = make(map[string]uint64)
