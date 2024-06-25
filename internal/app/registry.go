// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package app

import (
	"fmt"

	"github.com/nextmn/srv6/internal/ctrl"
	"github.com/nextmn/srv6/internal/iproute2"
)

type Registry struct {
	ifaces             map[string]*iproute2.TunIface
	controllerRegistry *ctrl.ControllerRegistry
}

func NewRegistry() *Registry {
	return &Registry{
		ifaces:             make(map[string]*iproute2.TunIface),
		controllerRegistry: nil,
	}
}

func (r *Registry) TunIface(name string) (*iproute2.TunIface, bool) {
	iface, exists := r.ifaces[name]
	return iface, exists
}

func (r *Registry) RegisterTunIface(iface *iproute2.TunIface) error {
	if _, exists := r.ifaces[iface.Name()]; exists {
		return fmt.Errorf("Iface %s is already registered.", iface.Name())
	}
	r.ifaces[iface.Name()] = iface
	return nil
}

func (r *Registry) DeleteTunIface(name string) {
	delete(r.ifaces, name)
}

func (r *Registry) RegisterControllerRegistry(cr *ctrl.ControllerRegistry) {
	r.controllerRegistry = cr
}

func (r *Registry) ControllerRegistry() (*ctrl.ControllerRegistry, bool) {
	if r.controllerRegistry == nil {
		return nil, false
	} else {
		return r.controllerRegistry, true
	}
}
func (r *Registry) DeleteControllerRegistry() {
	r.controllerRegistry = nil
}
