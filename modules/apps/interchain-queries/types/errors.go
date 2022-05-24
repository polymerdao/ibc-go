package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ICQ Host sentinel errors
var (
	ErrHostSubModuleDisabled = sdkerrors.Register(ModuleName, 2, "host module is disabled")
)
