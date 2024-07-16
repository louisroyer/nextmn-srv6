// Code generated by gen.go; DO NOT EDIT.

// Copyright 2024 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package database

type procedure struct {
	num_in int
	num_out int
}

var procedures = map[string]procedure{
	"insert_uplink_rule": procedure{num_in: 5, num_out: 1},
	"insert_downlink_rule": procedure{num_in: 4, num_out: 1},
	"enable_rule": procedure{num_in: 1, num_out: 0},
	"disable_rule": procedure{num_in: 1, num_out: 0},
	"delete_rule": procedure{num_in: 1, num_out: 0},
	"get_uplink_action": procedure{num_in: 3, num_out: 2},
	"set_uplink_action": procedure{num_in: 3, num_out: 2},
	"get_downlink_action": procedure{num_in: 1, num_out: 2},
	"get_rule": procedure{num_in: 1, num_out: 6},
}
