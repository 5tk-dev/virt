package virt

import (
	"net"
	"os"
	"sync"
)

var (
	vmListeners sync.Map
)

func CreateListener(qmpSockPath string) (net.Listener, error) {
	if _, ok := vmListeners.Load(qmpSockPath); ok {
		return nil, os.ErrExist
	}
	if err := os.RemoveAll(qmpSockPath); err != nil {
		return nil, err
	}
	listener, err := net.Listen("unix", qmpSockPath)
	if err != nil {
		return nil, err
	}
	if err := os.Chown(qmpSockPath, 1000, 1000); err != nil {
		return nil, err
	}
	vmListeners.Store(qmpSockPath, listener)
	return listener, nil
}

func DeleteListener(qmpSockPath string) error {
	if l, ok := vmListeners.LoadAndDelete(qmpSockPath); ok {
		if listener, ok := l.(net.Listener); ok {
			if err := listener.Close(); err != nil {
				return err
			}
		}
	}
	if err := os.RemoveAll(qmpSockPath); err != nil {
		return err
	}
	return nil
}
