// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package app_api

import (
	"github.com/nextmn/srv6/internal/ctrl"
	ctrl_api "github.com/nextmn/srv6/internal/ctrl/api"
	"github.com/nextmn/srv6/internal/database"
	"github.com/nextmn/srv6/internal/iproute2"
)

type Registry interface {
	// ifaces
	TunIface(name string) (*iproute2.TunIface, bool)
	RegisterTunIface(iface *iproute2.TunIface) error
	DeleteTunIface(name string)
	RegisterControllerRegistry(*ctrl.ControllerRegistry)
	ControllerRegistry() (*ctrl.ControllerRegistry, bool)
	DeleteControllerRegistry()
	RegisterDB(*database.Database)
	DB() (*database.Database, bool)
	DeleteDB()
	RegisterRulesRegistry(rr ctrl_api.RulesRegistry)
	RulesRegistry() (ctrl_api.RulesRegistry, bool)
	DeleteRulesRegistry()
}
