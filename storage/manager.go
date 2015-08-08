package storage

type managerStruct struct {
}

var manager *managerStruct

func getManager() *managerStruct {
	if manager == nil {
		manager = &managerStruct{}
	}
	return manager
}

/*
Manager is responsible for manipulating the counters and syncing to disk
*/
var Manager = getManager()
