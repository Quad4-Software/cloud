// Copyright 2026 Quad4 Cloud GmbH <mail@quad4.io>
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"

	"github.com/Quad4-Software/cloud/pkg/log"
	"github.com/Quad4-Software/cloud/services/auth-guest/pkg/config"
	"github.com/Quad4-Software/cloud/services/auth-guest/pkg/service/authguest"
	"github.com/spf13/pflag"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger  log.Logger
	Context context.Context
	Config  *config.Config
	Service authguest.AuthGuest
	Flags   []pflag.Flag
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Context provides a function to set the context option.
func Context(val context.Context) Option {
	return func(o *Options) {
		o.Context = val
	}
}

// Config provides a function to set the config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// Service provides a function to set the service option.
func Service(val authguest.AuthGuest) Option {
	return func(o *Options) {
		o.Service = val
	}
}

// Flags provides a function to set the flags option.
func Flags(flags ...pflag.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, flags...)
	}
}
