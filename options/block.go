package options

import (
	"fmt"
	cid "github.com/ipfs/go-cid"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
)

type BlockPutSettings struct {
	StoreCodec string
	// FIXME: Rename to Format (and possibly mark as deprecated).
	Codec    string
	MhType   uint64
	MhLength int
	Pin      bool
}

type BlockRmSettings struct {
	Force bool
}

type BlockPutOption func(*BlockPutSettings) error
type BlockRmOption func(*BlockRmSettings) error

func BlockPutOptions(opts ...BlockPutOption) (*BlockPutSettings, cid.Prefix, error) {
	options := &BlockPutSettings{
		Codec:      "",
		StoreCodec: "",
		MhType:     mh.SHA2_256,
		MhLength:   -1,
		Pin:        false,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, cid.Prefix{}, err
		}
	}

	if options.Codec != "" && options.StoreCodec != "" {
		return nil, cid.Prefix{}, fmt.Errorf("incompatible format (%s) and store-codec options set (%s)",
			options.Codec, options.StoreCodec)
	}

	if options.Codec == "" && options.StoreCodec == "" {
		// FIXME(BLOCKING): Do we keep the old default v0 here?
		options.Codec = "v0"
		// FIXME(BLOCKING): Review how to handle "protobuf". For now we simplify the code only with "v0".
	}

	var pref cid.Prefix
	pref.Version = 1

	// Old format option.
	if options.Codec != "" {
		if options.Codec == "v0" {
			if options.MhType != mh.SHA2_256 || (options.MhLength != -1 && options.MhLength != 32) {
				return nil, cid.Prefix{}, fmt.Errorf("only sha2-255-32 is allowed with CIDv0")
			}
			pref.Version = 0
		}

		// FIXME(BLOCKING): Do we actually want to consult the CID codecs table
		//  even with the old --format options? Or do we always want to check
		//  the multicodec one?
		cidCodec, ok := cid.Codecs[options.Codec]
		if !ok {
			return nil, cid.Prefix{}, fmt.Errorf("unrecognized format: %s", options.Codec)
		}
		pref.Codec = cidCodec
	} else {
		// New store-codec options. We handle it as it's done for `ipfs dag put`.
		var storeCodec mc.Code
		if err := storeCodec.Set(options.StoreCodec); err != nil {
			return nil, cid.Prefix{}, err
		}
		pref.Codec = uint64(storeCodec)
	}

	// FIXME: The entire codec manipulation/validation needs to be encapsulated
	//  outside this funtion to clearly demark that it is the only option we are
	//  overwriting here.

	pref.MhType = options.MhType
	pref.MhLength = options.MhLength

	return options, pref, nil
}

func BlockRmOptions(opts ...BlockRmOption) (*BlockRmSettings, error) {
	options := &BlockRmSettings{
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

type blockOpts struct{}

var Block blockOpts

// Format is an option for Block.Put which specifies the multicodec to use to
// serialize the object. Default is "v0"
func (blockOpts) Format(codec string) BlockPutOption {
	return func(settings *BlockPutSettings) error {
		settings.Codec = codec
		return nil
	}
}

// StoreCodec is an option for Block.Put which specifies the multicodec to use to
// serialize the object. It replaces the old Format now with CIDv1 as the default.
func (blockOpts) StoreCodec(storeCodec string) BlockPutOption {
	return func(settings *BlockPutSettings) error {
		settings.StoreCodec = storeCodec
		return nil
	}
}

// Hash is an option for Block.Put which specifies the multihash settings to use
// when hashing the object. Default is mh.SHA2_256 (0x12).
// If mhLen is set to -1, default length for the hash will be used
func (blockOpts) Hash(mhType uint64, mhLen int) BlockPutOption {
	return func(settings *BlockPutSettings) error {
		settings.MhType = mhType
		settings.MhLength = mhLen
		return nil
	}
}

// Pin is an option for Block.Put which specifies whether to (recursively) pin
// added blocks
func (blockOpts) Pin(pin bool) BlockPutOption {
	return func(settings *BlockPutSettings) error {
		settings.Pin = pin
		return nil
	}
}

// Force is an option for Block.Rm which, when set to true, will ignore
// non-existing blocks
func (blockOpts) Force(force bool) BlockRmOption {
	return func(settings *BlockRmSettings) error {
		settings.Force = force
		return nil
	}
}
