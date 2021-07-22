//  Copyright 2021-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included in
//  the file licenses/Couchbase-BSL.txt.  As of the Change Date specified in that
//  file, in accordance with the Business Source License, use of this software will
//  be governed by the Apache License, Version 2.0, included in the file
//  licenses/APL.txt.

package api

import "fmt"

type Value interface {
	String() string
	MarshalJSON() ([]byte, error)
	Actual() interface{}
	ToString() string
	Truth() bool
	Recycle()
	Track()
}

type Context interface {
	NewValue(val interface{}) interface{}
	ExecuteStatement(statement string, namedArgs map[string]interface{}, positionalArgs []interface{}) (interface{}, uint64, error)
}

func Args(args interface{}, context interface{}) (Value, Context, error) {
	a, ok := args.(Value)
	if !ok {
		return nil, nil, fmt.Errorf("invalid function arguments type %T", args)
	}
	c, ok := context.(Context)
	if !ok {
		return nil, nil, fmt.Errorf("invalid function context type %T", context)
	}
	return a, c, nil
}