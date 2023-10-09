// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package tasks

import (
	"fmt"

	app_api "github.com/nextmn/srv6/internal/app/api"
	"github.com/nextmn/srv6/internal/config"
	"github.com/nextmn/srv6/internal/constants"
	"github.com/nextmn/srv6/internal/iproute2"
	"github.com/nextmn/srv6/internal/netfunc"
	netfunc_api "github.com/nextmn/srv6/internal/netfunc/api"
)

// TaskNextMNHeadend creates a new headend
type TaskNextMNHeadend struct {
	WithState
	headend    *config.Headend
	table      iproute2.Table
	registry   app_api.Registry
	iface_name string
	netfunc    netfunc_api.NetFunc
}

// Create a new TaskNextMNHeadend
func NewTaskNextMNHeadend(headend *config.Headend, table_name string, iface_name string, registry app_api.Registry) *TaskNextMNHeadend {
	return &TaskNextMNHeadend{
		WithState:  NewState(),
		headend:    headend,
		table:      iproute2.NewTable(table_name, constants.RT_PROTO_NEXTMN),
		iface_name: iface_name,
		registry:   registry,
		netfunc:    nil,
	}
}

// Init
func (t *TaskNextMNHeadend) RunInit() error {
	// Create and start headend
	tunIface, ok := t.registry.TunIface(t.iface_name)
	if !ok {
		return fmt.Errorf("Interface %s is not in registry", t.iface_name)
	}
	if ep, err := netfunc.NewHeadend(t.headend); err != nil {
		return err
	} else {
		t.netfunc = ep
	}
	t.netfunc.Start(tunIface)
	// Add route to headend
	if err := t.table.AddRoute4Tun(t.headend.To, t.iface_name); err != nil {
		return err
	}
	t.state = true
	return nil
}

// Exit
func (t *TaskNextMNHeadend) RunExit() error {
	// Remove route to endpoint
	if err := t.table.DelRoute4Tun(t.headend.To, t.iface_name); err != nil {
		return err
	}
	// Stop headend
	t.netfunc.Stop()
	t.state = false
	return nil
}
