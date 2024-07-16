// Copyright 2024 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package database_api

import (
	"net/netip"

	"github.com/nextmn/json-api/jsonapi"
)

type Uplink interface {
	GetUplinkAction(UplinkTeid uint32, SrgwIp netip.Addr, GnbIp netip.Addr) (jsonapi.Action, error)
	SetUplinkAction(UplinkTeid uint32, SrgwIp netip.Addr, GnbIp netip.Addr, UeIpAddress netip.Addr) (jsonapi.Action, error)
}
