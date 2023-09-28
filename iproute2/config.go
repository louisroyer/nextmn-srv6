// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package iproute2

import (
	"fmt"
	"log"
	"os"
)

// This name will appear in configuration files
const generatorName = "NextMN"

type Kind int

const (
	RTTable Kind = iota + 1 // RT Table
	Proto                   // Proto
)

// Holds configuration for a RT Table or a Proto
type Config struct {
	name   string // Name
	kind   Kind   // RT Table or Proto
	number uint32 // RT Table / Proto number
}

// Creates new configuration
func NewConfig(name string, kind Kind, number uint32) (*Config, error) {
	if kind == Proto && number > 255 {
		return nil, fmt.Errorf("Proto number maximum is 255")
	}
	return &Config{name: name, kind: kind, number: number}, nil
}

// Returns filename used to store configuration
func (i Config) filename() (string, error) {
	switch i.kind {
	case RTTable:
		return fmt.Sprintf("/etc/iproute2/rt_tables.d/%s.conf", i.name), nil
	case Proto:
		return fmt.Sprintf("/etc/iproute2/rt_protos.d/%s.conf", i.name), nil
	default:
		return "", fmt.Errorf("Unknown IPRoute2Kind: %s", i.kind)
	}
}

// Return name of the RT Table or of the Proto
func (i Config) Name() string {
	return i.name
}

// Create config file
func (i Config) Create() error {
	// TODO: check that i.number is unused
	filename, err := i.filename()
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()
	str := fmt.Sprintf("# This file is autogenerated by %s\n%d\t%s\n", generatorName, i.number, i.name)
	if _, err = f.WriteString(str); err != nil {
		return err
	}
	return nil
}

// Delete config file
func (i Config) Delete() error {
	filename, err := i.filename()
	if err != nil {
		return err
	}
	if err := os.Remove(filename); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
