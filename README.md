# udf-api

The udf-api repository defines an API allowing shared objects running inside the Couchbase Query Service
to execute arbitrary N1QL statements.

This could be used by User Defined Functions written in Golang and dynamically loaded by the Query service
or by the Javascript runner executing a Javascript user defined function needing to switch to N1QL.

## introduction

The API has been introduced due to the general difficulty of building and deploying shared objects to be
used with golang applications: generally speaking, the plugin package requires that plugins be compiled
against the exact same version of the packages they reference as the running program, including source
package paths.

Having to build shared objects as part of a full Couchbase Server build would defy the main purpose of
having shared objects in the first place, which is to execute user provided code built independently of
the main binary.

In order to get round these limitations, udf-api defines a set of interfaces that the N1QL service
implements, but it never actually directly references.

This means that the API itself never references any complex types defined by the N1QL code - or the plugin
code would trigger version checking - which makes its usage slightly awkward.

API methods are accessed directly via arguments passed in shared object calls and never via
referencing N1QL packages directly, thereby bypassing the need for most runtime version checks.

The repository will have a package for each new version of the APIs, with each version being backwards
compatible.

The reason for this is to allow shared objects that have already been built to keep on working without
being rebuilt even on different versions of the N1QL service (with the proviso that all the N1QL
versions support that particular version of the API).

## usage

- set GOPATH as usual
- create you shared object package in $GOPATH/src
- git clone http://github.com/couchbase/udf-api
- in your code, import "github.com/couchbase/udf-api/v1" (substitute v1 with the version of the api you intend to use)
- exported symbols should have names starting with an uppercase character, and have a signature of the form `func(interface{},
interface{}) (interface{}, error)`, where the first argument is the argument list, and the second the execution context
- go build -buildmode=plugin .

## sample user defined function code

    package main
    
    import (
            "fmt"
    
            "github.com/couchbase/udf-api/v1"
    )
    
    func Udfquery(a interface{}, c interface{}) (interface{}, error) {
            args, context, err:= api.Args(a, c)
            if err != nil {
                    return nil, err
            }
            argMap, ok := args.Actual().(map[string]interface{})
            if !ok  || len(argMap) != 1 {
                    return nil, fmt.Errorf("invalid arguments, %T %v", args, args)
            }
            val := argMap["pType"]
            if val == nil  {
                    return nil, fmt.Errorf("missing argument 'pType'")
            }
            rv, _, err := context.ExecuteStatement("select * from `travel-sample` where type=$1", nil, []interface{}{ val })
            return rv, err
    }

## limitations

Windows is not supported.

Although there is prototype code that loads shared objects, Windows Go ports do not support building shared objects in the first place.
