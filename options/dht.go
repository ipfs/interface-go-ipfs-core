package options

// Deprecated: use github.com/ipfs/boxo/coreiface/options.DhtProvideSettings
type DhtProvideSettings struct {
	Recursive bool
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.DhtFindProvidersSettings
type DhtFindProvidersSettings struct {
	NumProviders int
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.DhtProvideOption
type DhtProvideOption func(*DhtProvideSettings) error

// Deprecated: use github.com/ipfs/boxo/coreiface/options.DhtFindProvidersOption
type DhtFindProvidersOption func(*DhtFindProvidersSettings) error

// Deprecated: use github.com/ipfs/boxo/coreiface/options.DhtProvideOptions
func DhtProvideOptions(opts ...DhtProvideOption) (*DhtProvideSettings, error) {
	options := &DhtProvideSettings{
		Recursive: false,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.DhtFindProvidersOptions
func DhtFindProvidersOptions(opts ...DhtFindProvidersOption) (*DhtFindProvidersSettings, error) {
	options := &DhtFindProvidersSettings{
		NumProviders: 20,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

type dhtOpts struct{}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.Dht
var Dht dhtOpts

// Recursive is an option for Dht.Provide which specifies whether to provide
// the given path recursively
func (dhtOpts) Recursive(recursive bool) DhtProvideOption {
	return func(settings *DhtProvideSettings) error {
		settings.Recursive = recursive
		return nil
	}
}

// NumProviders is an option for Dht.FindProviders which specifies the
// number of peers to look for. Default is 20
func (dhtOpts) NumProviders(numProviders int) DhtFindProvidersOption {
	return func(settings *DhtFindProvidersSettings) error {
		settings.NumProviders = numProviders
		return nil
	}
}
