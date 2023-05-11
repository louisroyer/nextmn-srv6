// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package srv6

import (
	"fmt"
	"log"
	"os"
)

func ipRoute2Init() error {
	if err := createRtTable(); err != nil {
		return err
	}
	if err := createRtProto(); err != nil {
		return err
	}
	if err := addIPRules(); err != nil {
		return err
	}
	if err := addDefaultRoute(); err != nil {
		return err
	}
	return nil
}

func ipRoute2Exit() error {
	if err := removeDefaultRoute(); err != nil {
		return err
	}
	if err := removeIPRules(); err != nil {
		return err
	}
	if err := removeRtProto(); err != nil {
		return err
	}
	if err := removeRtTable(); err != nil {
		return err
	}
	return nil
}

func addIPRules() error {
	return runIP("-6", "rule", "add", "to", SRv6.Locator, "lookup", RTTableName, "protocol", RTProtoName)
}

func addDefaultRoute() error {
	// This default route will be replaced later with a route to sr TUN interface
	return runIP("-6", "route", "add", "blackhole", "default", "table", RTTableName, "proto", RTProtoName)
}

func removeDefaultRoute() error {
	return runIP("-6", "route", "del", "blackhole", "default", "table", RTTableName, "proto", RTProtoName)
}

func removeIPRules() error {
	return runIP("-6", "rule", "del", "to", SRv6.Locator, "lookup", RTTableName, "protocol", RTProtoName)
}

func createRtTable() error {
	// TODO: check that RTTableNumber is unused
	f, err := os.Create(RTTableFileName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()
	str := fmt.Sprintf("# This file is autogenerated by nextmn-srv6\n%d\t%s\n", SRv6.IPRoute2.RTTableNumber, RTTableName)
	if _, err = f.WriteString(str); err != nil {
		return err
	}
	return nil
}

func removeRtTable() error {
	if err := os.Remove(RTTableFileName); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func createRtProto() error {
	// TODO: check that RTProtoNumber is unused
	f, err := os.Create(RTProtoFileName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()
	str := fmt.Sprintf("# This file is autogenerated by nextmn-srv6\n%d\t%s\n", SRv6.IPRoute2.RTProtoNumber, RTProtoName)
	if _, err = f.WriteString(str); err != nil {
		return err
	}
	return nil
}

func removeRtProto() error {
	if err := os.Remove(RTProtoFileName); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
