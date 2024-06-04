package manager

import (
	"github.com/sleepdeprecation/tfmgr/internal/config"
	"github.com/sleepdeprecation/tfmgr/internal/downloader"
)

type Manager struct {
	Config         *config.Config
	Downloader     *downloader.Downloader
	CurrentVersion string
}

// func New() *Manager {
// 	m := &Manager{
// 		Config:     config.Get(),
// 		Downloader: downloader.New(),
// 	}
//
// 	m.DetectVersion()
// 	return m
// }
//
// func (m *Manager) DetectVersion() {
// 	version := m.Config.DefaultVersion
// }
//
// func (m *Manager) Use(version string) {
// 	m.CurrentVersion = version
// }
