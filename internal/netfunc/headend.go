// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package netfunc

import (
	"fmt"
	"net/netip"

	"github.com/nextmn/srv6/internal/config"
	netfunc_api "github.com/nextmn/srv6/internal/netfunc/api"
)

func NewHeadend(he *config.Headend, ttl uint8, hopLimit uint8, debug bool) (netfunc_api.NetFunc, error) {
	p, err := netip.ParsePrefix(he.To)
	if err != nil {
		return nil, err
	}
	switch he.Behavior {
	case config.H_M_GTP4_D:
		policy := he.Policy
		return NewNetFunc(NewHeadendGTP4(p, policy, ttl, hopLimit), debug), nil
	default:
		return nil, fmt.Errorf("Unsupported headend behavior (%s) with this provider (%s)", he.Behavior, he.Provider)
	}
}
