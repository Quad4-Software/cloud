package log

import "github.com/Quad4-Software/cloud/pkg/shared"

// Configure initializes a service-specific logger instance.
func Configure(name string, commons *shared.Commons, localServiceLogLevel string) Logger {
	return NewLogger(
		Name(name),
		Level(localServiceLogLevel),
		Pretty(commons.Log.Pretty),
		Color(commons.Log.Color),
		File(commons.Log.File),
	)
}
