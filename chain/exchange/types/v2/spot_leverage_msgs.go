package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// 确保现货杠杆消息实现了 sdk.Msg 接口
var (
	_ sdk.Msg = &MsgEnableSpotLeverage{}
	_ sdk.Msg = &MsgDisableSpotLeverage{}
	_ sdk.Msg = &MsgSupplyToLendingPool{}
	_ sdk.Msg = &MsgWithdrawFromLendingPool{}
	_ sdk.Msg = &MsgCreateSpotLeverageLimitOrder{}
	_ sdk.Msg = &MsgCreateSpotLeverageMarketOrder{}
	_ sdk.Msg = &MsgCancelSpotLeverageOrder{}
)

// ============================================================================
// MsgEnableSpotLeverage
// ============================================================================

// Route implements the sdk.Msg interface
func (msg *MsgEnableSpotLeverage) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg *MsgEnableSpotLeverage) Type() string { return "enableSpotLeverage" }

// GetSigners implements the sdk.Msg interface
func (msg *MsgEnableSpotLeverage) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgEnableSpotLeverage) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgEnableSpotLeverage) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.SubaccountId == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("subaccount_id cannot be empty")
	}
	return nil
}

// ============================================================================
// MsgDisableSpotLeverage
// ============================================================================

// Route implements the sdk.Msg interface
func (msg *MsgDisableSpotLeverage) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg *MsgDisableSpotLeverage) Type() string { return "disableSpotLeverage" }

// GetSigners implements the sdk.Msg interface
func (msg *MsgDisableSpotLeverage) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgDisableSpotLeverage) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgDisableSpotLeverage) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.SubaccountId == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("subaccount_id cannot be empty")
	}
	return nil
}

// ============================================================================
// MsgSupplyToLendingPool
// ============================================================================

// Route implements the sdk.Msg interface
func (msg *MsgSupplyToLendingPool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg *MsgSupplyToLendingPool) Type() string { return "supplyToLendingPool" }

// GetSigners implements the sdk.Msg interface
func (msg *MsgSupplyToLendingPool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgSupplyToLendingPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgSupplyToLendingPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.Denom == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("denom cannot be empty")
	}
	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("amount must be positive")
	}
	return nil
}

// ============================================================================
// MsgWithdrawFromLendingPool
// ============================================================================

// Route implements the sdk.Msg interface
func (msg *MsgWithdrawFromLendingPool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg *MsgWithdrawFromLendingPool) Type() string { return "withdrawFromLendingPool" }

// GetSigners implements the sdk.Msg interface
func (msg *MsgWithdrawFromLendingPool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgWithdrawFromLendingPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgWithdrawFromLendingPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.Denom == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("denom cannot be empty")
	}
	if msg.Shares.IsNegative() || msg.Shares.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("shares must be positive")
	}
	return nil
}

// ============================================================================
// MsgCreateSpotLeverageLimitOrder
// ============================================================================

// Route implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageLimitOrder) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageLimitOrder) Type() string { return "createSpotLeverageLimitOrder" }

// GetSigners implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageLimitOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageLimitOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageLimitOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.Order.MarketId == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("order market_id cannot be empty")
	}
	return nil
}

// ============================================================================
// MsgCreateSpotLeverageMarketOrder
// ============================================================================

// Route implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageMarketOrder) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageMarketOrder) Type() string { return "createSpotLeverageMarketOrder" }

// GetSigners implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageMarketOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageMarketOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgCreateSpotLeverageMarketOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.Order.MarketId == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("order market_id cannot be empty")
	}
	return nil
}

// ============================================================================
// MsgCancelSpotLeverageOrder
// ============================================================================

// Route implements the sdk.Msg interface
func (msg *MsgCancelSpotLeverageOrder) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg *MsgCancelSpotLeverageOrder) Type() string { return "cancelSpotLeverageOrder" }

// GetSigners implements the sdk.Msg interface
func (msg *MsgCancelSpotLeverageOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgCancelSpotLeverageOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgCancelSpotLeverageOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}
	if msg.MarketId == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("market_id cannot be empty")
	}
	if msg.OrderHash == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("order_hash cannot be empty")
	}
	return nil
}
