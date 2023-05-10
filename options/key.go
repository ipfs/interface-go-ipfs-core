package options

const (
	// Deprecated: use github.com/ipfs/boxo/coreiface/options.RSAKey
	RSAKey = "rsa"
	// Deprecated: use github.com/ipfs/boxo/coreiface/options.Ed25519Key
	Ed25519Key = "ed25519"

	// Deprecated: use github.com/ipfs/boxo/coreiface/options.DefaultRSALen
	DefaultRSALen = 2048
)

// Deprecated: use github.com/ipfs/boxo/coreiface/options.KeyGenerateSettings
type KeyGenerateSettings struct {
	Algorithm string
	Size      int
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.KeyRenameSettings
type KeyRenameSettings struct {
	Force bool
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.KeyGenerateOption
type KeyGenerateOption func(*KeyGenerateSettings) error

// Deprecated: use github.com/ipfs/boxo/coreiface/options.KeyRenameOption
type KeyRenameOption func(*KeyRenameSettings) error

// Deprecated: use github.com/ipfs/boxo/coreiface/options.KeyGenerateOptions
func KeyGenerateOptions(opts ...KeyGenerateOption) (*KeyGenerateSettings, error) {
	options := &KeyGenerateSettings{
		Algorithm: RSAKey,
		Size:      -1,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.KeyRenameOptions
func KeyRenameOptions(opts ...KeyRenameOption) (*KeyRenameSettings, error) {
	options := &KeyRenameSettings{
		Force: false,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

type keyOpts struct{}

// Deprecated: use github.com/ipfs/boxo/coreiface/options.Key
var Key keyOpts

// Type is an option for Key.Generate which specifies which algorithm
// should be used for the key. Default is options.RSAKey
//
// Supported key types:
// * options.RSAKey
// * options.Ed25519Key
func (keyOpts) Type(algorithm string) KeyGenerateOption {
	return func(settings *KeyGenerateSettings) error {
		settings.Algorithm = algorithm
		return nil
	}
}

// Size is an option for Key.Generate which specifies the size of the key to
// generated. Default is -1
//
// value of -1 means 'use default size for key type':
//   - 2048 for RSA
func (keyOpts) Size(size int) KeyGenerateOption {
	return func(settings *KeyGenerateSettings) error {
		settings.Size = size
		return nil
	}
}

// Force is an option for Key.Rename which specifies whether to allow to
// replace existing keys.
func (keyOpts) Force(force bool) KeyRenameOption {
	return func(settings *KeyRenameSettings) error {
		settings.Force = force
		return nil
	}
}
