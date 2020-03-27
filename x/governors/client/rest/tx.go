package rest

import (
	"bytes"
	"github.com/xar-network/xar-network/x/governors/types"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/governors/depositors/{depositorAddr}/deposits",
		postDepositsHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/governors/depositors/{depositorAddr}/unbonding_deposits",
		postUnbondingDepositsHandlerFn(cliCtx),
	).Methods("POST")
}

type (
	// BondRequest defines the properties of a deposit request's body.
	BondRequest struct {
		BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
		DepositorAddress sdk.AccAddress `json:"depositor_address" yaml:"depositor_address"` // in bech32
		Amount           sdk.Coin       `json:"amount" yaml:"amount"`
	}

	// UnbondRequest defines the properties of a unbond request's body.
	UnbondRequest struct {
		BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
		DepositorAddress sdk.AccAddress `json:"depositor_address" yaml:"depositor_address"` // in bech32
		Amount           sdk.Coin       `json:"amount" yaml:"amount"`
	}
)

func postDepositsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BondRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgBond(req.DepositorAddress, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, req.DepositorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own depositor address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postUnbondingDepositsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UnbondRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgUnbond(req.DepositorAddress, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, req.DepositorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own depositor address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
