// Package uiconfigmodel
package uiconfigmodel

type UIConfig = map[string]string

type UIConfigChangeHook = func(newConfig UIConfig)
