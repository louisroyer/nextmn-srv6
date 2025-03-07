// Copyright 2023 Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Provider int

const (
	ProviderLinux Provider = iota
	ProviderNextMN
	ProviderNextMNWithController
)

func (p Provider) String() string {
	switch p {
	case ProviderLinux:
		return "Linux"
	case ProviderNextMN:
		return "NextMN"
	case ProviderNextMNWithController:
		return "NextMN (via controller)"
	default:
		return "Unknown provider"
	}
}

// Unmarshal YAML to Provider
func (p *Provider) UnmarshalYAML(n *yaml.Node) error {
	switch strings.ToLower(n.Value) {
	case "linux":
		*p = ProviderLinux
	case "nextmn":
		*p = ProviderNextMN
	case "nextmn-ctrl", "nextmnwithcontroller", "nextmn-with-controller", "nextmn-via-controller":
		*p = ProviderNextMNWithController
	default:
		return fmt.Errorf("Unknown provider")
	}
	return nil
}
