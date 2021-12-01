// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	blitzshare "github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/blitzshare"
	mock "github.com/stretchr/testify/mock"
)

// BlitzshareApi is an autogenerated mock type for the BlitzshareApi type
type BlitzshareApi struct {
	mock.Mock
}

// GetPeerAddr provides a mock function with given fields: oneTimePass
func (_m *BlitzshareApi) GetPeerAddr(oneTimePass *string) *blitzshare.PeerAddress {
	ret := _m.Called(oneTimePass)

	var r0 *blitzshare.PeerAddress
	if rf, ok := ret.Get(0).(func(*string) *blitzshare.PeerAddress); ok {
		r0 = rf(oneTimePass)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*blitzshare.PeerAddress)
		}
	}

	return r0
}

// RegisterAsPeer provides a mock function with given fields: multiAddr, oneTimePass
func (_m *BlitzshareApi) RegisterAsPeer(multiAddr string, oneTimePass string) bool {
	ret := _m.Called(multiAddr, oneTimePass)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(multiAddr, oneTimePass)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
