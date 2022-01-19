package master

import (
	"path/filepath"
	"strings"

	fcom "github.com/meshplus/hyperbench-common/common"

	"github.com/meshplus/hyperbench/plugins/blockchain"
	"github.com/meshplus/hyperbench/vm"
	"github.com/meshplus/hyperbench/vm/base"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Master is the interface of mater node
type Master interface {
	// Prepare is used to prepare
	Prepare() error

	// GetContext generate the context, which will be passed to Worker
	GetContext() ([]byte, error)

	// Statistic query the remote statistic data from chain
	Statistic(from, to int64) (*fcom.RemoteStatistic, error)

	// LogStatus records blockheight and time
	LogStatus() (int64, error)
}

// LocalMaster is the implement of master in local
type LocalMaster struct {
	masterVM vm.VM
	params   []string
}

// Prepare is used to prepare
func (m *LocalMaster) Prepare() (err error) {
	// call user hook
	err = m.masterVM.BeforeDeploy()
	if err != nil {
		return errors.Wrap(err, "can not call user hook `BeforeDeploy`")
	}

	// prepare contract
	err = m.masterVM.DeployContract()
	if err != nil {
		return errors.Wrap(err, "can not deploy contract")
	}

	return nil
}

// GetContext generate the context, which will be passed to Worker
func (m *LocalMaster) GetContext() ([]byte, error) {
	err := m.masterVM.BeforeGet()
	if err != nil {
		return nil, err
	}
	return m.masterVM.GetContext()
}

// Statistic query the remote statistic data from chain
func (m *LocalMaster) Statistic(from, to int64) (*fcom.RemoteStatistic, error) {
	return m.masterVM.Statistic(from, to)
}

// LogStatus records blockheight and time
func (m *LocalMaster) LogStatus() (end int64, err error) {
	return m.masterVM.LogStatus()
}

// NewLocalMaster create LocalMaster.
func NewLocalMaster() (*LocalMaster, error) {
	blockchain.InitPlugin()

	params := viper.GetStringSlice(fcom.ClientContractArgsPath)
	scriptPath := viper.GetString(fcom.ClientScriptPath)
	vmType := strings.TrimPrefix(filepath.Ext(scriptPath), ".")
	masterVM, err := vm.NewVM(vmType, base.ConfigBase{
		Path: scriptPath,
		Ctx: fcom.VMContext{
			WorkerIdx: -1,
			VMIdx:     -1,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "can not create master")
	}

	return &LocalMaster{
		masterVM: masterVM,
		params:   params,
	}, nil
}
