// Copyright 2023 Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package config

type Policy struct {
	Match *Match `yaml:"match,omitempty"`
	Bsid  Bsid   `yaml:"bsid"`
}
