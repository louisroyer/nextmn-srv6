// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package tasks

import (
	app_api "github.com/nextmn/srv6/internal/app/api"
	iproute2 "github.com/nextmn/srv6/internal/iproute2"
	iproute2_api "github.com/nextmn/srv6/internal/iproute2/api"
)

// TaskIface
type TaskIface struct {
	WithState
	iface    iproute2_api.Iface
	registry app_api.Registry
}

// Create a new Task for DummyIface
func NewTaskDummyIface(name string) *TaskIface {
	return &TaskIface{
		iface:    iproute2.NewDummyIface(name),
		registry: nil,
	}
}

// Create a new Task for TunIface
func NewTaskTunIface(name string, registry app_api.Registry) *TaskIface {
	return &TaskIface{
		iface:    iproute2.NewTunIface(name),
		registry: registry,
	}
}

// Create and set up the Iface
func (t *TaskIface) RunInit() error {
	if err := t.iface.CreateAndUp(); err != nil {
		return err
	}
	if t.registry != nil {
		if err := t.registry.RegisterIface(t.iface); err != nil {
			return err
		}
	}
	t.state = true
	return nil
}

// Delete the Iface
func (t *TaskIface) RunExit() error {
	if err := t.iface.Delete(); err != nil {
		return err
	}
	if t.registry != nil {
		t.registry.DeleteIface(t.iface.Name())
	}
	t.state = false
	return nil
}
