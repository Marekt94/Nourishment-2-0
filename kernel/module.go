package kernel

import "nourishment_20/internal/logging"

const ADMIN_USER_NAME = "ADMIN"
const READER_USER_NAME = "READER"

type ModuleIntf interface {
	ExposeMethods()
	RegisterPermissions()
	GetName() string
}

type KernelIntf interface {
	RegisterModule(m ModuleIntf)
	Run()
	Init()
}

type Kernel struct {
	modules map[string]ModuleIntf
}

func NewKernel() KernelIntf {
	return &Kernel{make(map[string]ModuleIntf)}
}

func (k *Kernel) Init() {

}

func (k *Kernel) RegisterModule(m ModuleIntf) {
	_, exists := k.modules[m.GetName()]
	if !exists {
		k.modules[m.GetName()] = m
		logging.Global.Infof(`Module "%v" registered`, m.GetName())
	} else {
		logging.Global.Panicf(`Module "%v" duplicated`, m.GetName())
	}
}

func (k *Kernel) Run() {
	for _, m := range k.modules {
		m.RegisterPermissions()
	}
	for _, m := range k.modules {
		m.ExposeMethods()
	}
}
