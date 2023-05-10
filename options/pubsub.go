package options

// Deprecated: use github.com/ipfs/boxo/coreiface/options.PubSubPeersSettings
type PubSubPeersSettings struct {
	Topic string
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.PubSubSubscribeSettings
type PubSubSubscribeSettings struct {
	Discover bool
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.PubSubPeersOption
type PubSubPeersOption func(*PubSubPeersSettings) error

// Deprecated: use github.com/ipfs/boxo/coreiface/options.PubSubSubscribeOption
type PubSubSubscribeOption func(*PubSubSubscribeSettings) error

// Deprecated: use github.com/ipfs/boxo/coreiface/options.PubSubPeersOptions
func PubSubPeersOptions(opts ...PubSubPeersOption) (*PubSubPeersSettings, error) {
	options := &PubSubPeersSettings{
		Topic: "",
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.PubSubSubscribeOptions
func PubSubSubscribeOptions(opts ...PubSubSubscribeOption) (*PubSubSubscribeSettings, error) {
	options := &PubSubSubscribeSettings{
		Discover: false,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

type pubsubOpts struct{}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.PubSub
var PubSub pubsubOpts

func (pubsubOpts) Topic(topic string) PubSubPeersOption {
	return func(settings *PubSubPeersSettings) error {
		settings.Topic = topic
		return nil
	}
}

func (pubsubOpts) Discover(discover bool) PubSubSubscribeOption {
	return func(settings *PubSubSubscribeSettings) error {
		settings.Discover = discover
		return nil
	}
}
