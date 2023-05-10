package options

// Deprecated: use github.com/ipfs/boxo/coreiface/options.ApiSettings
type ApiSettings struct {
	Offline     bool
	FetchBlocks bool
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.ApiOption
type ApiOption func(*ApiSettings) error

// Deprecated: use github.com/ipfs/boxo/coreiface/options.ApiOptions
func ApiOptions(opts ...ApiOption) (*ApiSettings, error) {
	options := &ApiSettings{
		Offline:     false,
		FetchBlocks: true,
	}

	return ApiOptionsTo(options, opts...)
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.ApiOptionsTo
func ApiOptionsTo(options *ApiSettings, opts ...ApiOption) (*ApiSettings, error) {
	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

type apiOpts struct{}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.Api
var Api apiOpts

func (apiOpts) Offline(offline bool) ApiOption {
	return func(settings *ApiSettings) error {
		settings.Offline = offline
		return nil
	}
}

// FetchBlocks when set to false prevents api from fetching blocks from the
// network while allowing other services such as IPNS to still be online
func (apiOpts) FetchBlocks(fetch bool) ApiOption {
	return func(settings *ApiSettings) error {
		settings.FetchBlocks = fetch
		return nil
	}
}
