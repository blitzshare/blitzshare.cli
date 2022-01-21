// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	bufio "bufio"

	config "bootstrap.cli/app/config"
	host "github.com/libp2p/go-libp2p-core/host"

	mock "github.com/stretchr/testify/mock"

	network "github.com/libp2p/go-libp2p-core/network"
)

// P2p is an autogenerated mock type for the P2p type
type P2p struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *P2p) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConnectToBootsrapNode provides a mock function with given fields: conf
func (_m *P2p) ConnectToBootsrapNode(conf *config.AppConfig) *host.Host {
	ret := _m.Called(conf)

	var r0 *host.Host
	if rf, ok := ret.Get(0).(func(*config.AppConfig) *host.Host); ok {
		r0 = rf(conf)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*host.Host)
		}
	}

	return r0
}

// ConnectToPeer provides a mock function with given fields: conf, address, otp
func (_m *P2p) ConnectToPeer(conf *config.AppConfig, address *string, otp *string) *bufio.ReadWriter {
	ret := _m.Called(conf, address, otp)

	var r0 *bufio.ReadWriter
	if rf, ok := ret.Get(0).(func(*config.AppConfig, *string, *string) *bufio.ReadWriter); ok {
		r0 = rf(conf, address, otp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bufio.ReadWriter)
		}
	}

	return r0
}

// StartPeer provides a mock function with given fields: conf, otp, handler
func (_m *P2p) StartPeer(conf *config.AppConfig, otp *string, handler func(network.Stream)) string {
	ret := _m.Called(conf, otp, handler)

	var r0 string
	if rf, ok := ret.Get(0).(func(*config.AppConfig, *string, func(network.Stream)) string); ok {
		r0 = rf(conf, otp, handler)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
