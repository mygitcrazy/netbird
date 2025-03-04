//go:build android
// +build android

package iface

import (
	"strings"

	"github.com/pion/transport/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"

	"github.com/netbirdio/netbird/iface/bind"
)

type tunDevice struct {
	address    WGAddress
	mtu        int
	tunAdapter TunAdapter
	iceBind    *bind.ICEBind

	fd      int
	name    string
	device  *device.Device
	wrapper *DeviceWrapper
}

func newTunDevice(address WGAddress, mtu int, tunAdapter TunAdapter, transportNet transport.Net) *tunDevice {
	return &tunDevice{
		address:    address,
		mtu:        mtu,
		tunAdapter: tunAdapter,
		iceBind:    bind.NewICEBind(transportNet),
	}
}

func (t *tunDevice) Create(mIFaceArgs MobileIFaceArguments) error {
	log.Info("create tun interface")
	var err error
	routesString := t.routesToString(mIFaceArgs.Routes)
	searchDomainsToString := t.searchDomainsToString(mIFaceArgs.SearchDomains)
	t.fd, err = t.tunAdapter.ConfigureInterface(t.address.String(), t.mtu, mIFaceArgs.Dns, searchDomainsToString, routesString)
	if err != nil {
		log.Errorf("failed to create Android interface: %s", err)
		return err
	}

	tunDevice, name, err := tun.CreateUnmonitoredTUNFromFD(t.fd)
	if err != nil {
		unix.Close(t.fd)
		return err
	}
	t.name = name
	t.wrapper = newDeviceWrapper(tunDevice)

	log.Debugf("attaching to interface %v", name)
	t.device = device.NewDevice(t.wrapper, t.iceBind, device.NewLogger(device.LogLevelSilent, "[wiretrustee] "))
	// without this property mobile devices can discover remote endpoints if the configured one was wrong.
	// this helps with support for the older NetBird clients that had a hardcoded direct mode
	// t.device.DisableSomeRoamingForBrokenMobileSemantics()

	err = t.device.Up()
	if err != nil {
		t.device.Close()
		return err
	}
	log.Debugf("device is ready to use: %s", name)
	return nil
}

func (t *tunDevice) Device() *device.Device {
	return t.device
}

func (t *tunDevice) DeviceName() string {
	return t.name
}

func (t *tunDevice) WgAddress() WGAddress {
	return t.address
}

func (t *tunDevice) UpdateAddr(addr WGAddress) error {
	// todo implement
	return nil
}

func (t *tunDevice) Close() (err error) {
	if t.device != nil {
		t.device.Close()
	}

	return
}

func (t *tunDevice) routesToString(routes []string) string {
	return strings.Join(routes, ";")
}

func (t *tunDevice) searchDomainsToString(searchDomains []string) string {
	return strings.Join(searchDomains, ";")
}
