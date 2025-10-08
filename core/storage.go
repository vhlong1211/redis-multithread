package core

import "TCPServer/data_structure"

var dictStore *data_structure.Dict
var zsetStore map[string]*data_structure.ZSet
var setStore map[string]*data_structure.SimpleSet

func init() {
	dictStore = data_structure.CreateDict()
	zsetStore = make(map[string]*data_structure.ZSet)
	setStore = make(map[string]*data_structure.SimpleSet)
}
