// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package config

type Headend struct {
	Name                string
	To                  string // IP Prefix this Headend will handle (can be the same as GTP4HeadendPrefix if you have a single Headend)
	Provider            Provider
	Behavior            HeadendBehavior
	Policy              *Policy
	SourceAddressPrefix *string `yaml:"set-source-prefix"`
}

type Headends []*Headend

func (he Headends) Filter(provider Provider) Headends {
	newList := make([]*Headend, 0)
	for _, e := range he {
		if e.Provider == provider {
			newList = append(newList, e)
		}
	}
	return newList
}

func (he Headends) FilterWithBehavior(provider Provider, behavior HeadendBehavior) Headends {
	newList := make([]*Headend, 0)
	for _, e := range he {
		if e.Provider == provider && e.Behavior == behavior {
			newList = append(newList, e)
		}
	}
	return newList
}
