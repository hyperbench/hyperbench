package vmpool

import (
	fcom "github.com/hyperbench/hyperbench-common/common"

	"path"
	"strings"

	"github.com/hyperbench/hyperbench/vm"
	"github.com/hyperbench/hyperbench/vm/base"
	"github.com/spf13/viper"
)

// Pool is the reusable vm.VM pool, Pop and Push is concurrent-safe
type Pool interface {
	// Pop gets a vm.VM from Pool concurrent-safely
	// if it's implement in no-block way, it may return nil
	Pop() vm.VM

	// Push sets a vm.WorkerVM to Pool concurrent-safely.
	Push(vm.VM)

	// Walk try to apply each vm.VM with wf until wf return false or all vm.WorkerVM in Pool have been applied.
	Walk(wf func(v vm.VM) (terminal bool))

	// Close close all vm.WorkerVM in Pool.
	Close()
}

// PoolImpl implement Pool.
type PoolImpl struct {
	ch chan vm.VM
	//len int64
	cap int64
}

// NewPoolImpl create PoolImpl.
func NewPoolImpl(workerID int64, cap int64) (*PoolImpl, error) {
	p := &PoolImpl{
		cap: cap,
		ch:  make(chan vm.VM, cap),
	}

	scriptPath := viper.GetString(fcom.ClientScriptPath)
	t := strings.TrimPrefix(path.Ext(scriptPath), ".")
	configBase := base.ConfigBase{
		Path: scriptPath,
		Ctx: fcom.VMContext{
			WorkerIdx: workerID,
			VMIdx:     0,
		},
	}
	configBase.Ctx.WorkerIdx = workerID
	var i int64
	fcom.GetLogger("pool").Notice(workerID, cap, scriptPath, t)
	for i = 0; i < cap; i++ {
		nvm, err := vm.NewVM(t, configBase)
		if err != nil {
			return nil, err
		}
		configBase.Ctx.VMIdx++
		// generate each vm with given index
		p.Push(nvm)
	}

	return p, nil
}

// Close close all vm.WorkerVM in Pool.
func (p *PoolImpl) Close() {
	p.Walk(func(v vm.VM) bool {
		v.Close()
		return false
	})
}

// Pop gets a vm.VM from Pool concurrent-safely
// if it's implement in no-block way, it may return nil.
func (p *PoolImpl) Pop() (worker vm.VM) {
	select {
	case worker = <-p.ch:
		return
	default:
		return
	}
}

// Push sets a vm.WorkerVM to Pool concurrent-safely.
func (p *PoolImpl) Push(worker vm.VM) {
	select {
	case p.ch <- worker:
	default:
	}
	return
}

// Walk try to apply each vm.VM with wf until wf return false or all vm.WorkerVM in Pool have been applied
func (p *PoolImpl) Walk(wf func(v vm.VM) bool) {
	l := len(p.ch)
	for i := 0; i < l; i++ {
		worker := <-p.ch
		exit := wf(worker)
		p.ch <- worker
		if exit {
			return
		}
	}
}
