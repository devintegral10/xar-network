package rest

import (
	"fmt"
	"github.com/xar-network/xar-network/x/bonds"
	"github.com/xar-network/xar-network/x/governors/types"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Get all deposits from a depositor
	r.HandleFunc(
		"/governors/depositors/{depositorAddr}/deposits",
		depositorDepositsHandlerFn(cliCtx),
	).Methods("GET")

	// Get all unbonding deposits from a depositor
	r.HandleFunc(
		"/governors/depositors/{depositorAddr}/unbonding_deposits",
		depositorUnbondingDepositsHandlerFn(cliCtx),
	).Methods("GET")

	// Get all governors txs (i.e msgs) from a depositor
	r.HandleFunc(
		"/governors/depositors/{depositorAddr}/txs",
		depositorTxsHandlerFn(cliCtx),
	).Methods("GET")

	// Query a deposit
	r.HandleFunc(
		"/governors/depositors/{depositorAddr}/deposits/",
		depositHandlerFn(cliCtx),
	).Methods("GET")

	// Query all unbonding deposits
	r.HandleFunc(
		"/governors/depositors/{depositorAddr}/unbonding_deposits/",
		unbondingDepositHandlerFn(cliCtx),
	).Methods("GET")

	// Get the current state of the pool
	r.HandleFunc(
		"/governors/pool",
		poolHandlerFn(cliCtx),
	).Methods("GET")

	// Get the current parameter values
	r.HandleFunc(
		"/governors/parameters",
		paramsHandlerFn(cliCtx),
	).Methods("GET")

}

// HTTP request handler to query a depositor deposits
func depositorDepositsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryDepositor(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute))
}

// HTTP request handler to query a depositor unbonding deposits
func depositorUnbondingDepositsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryDepositor(cliCtx, "custom/governors/depositorUnbondingDeposits")
}

// HTTP request handler to query all txs (msgs) from a depositor
func depositorTxsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var typesQuerySlice []string
		vars := mux.Vars(r)
		depositorAddr := vars["depositorAddr"]

		_, err := sdk.AccAddressFromBech32(depositorAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		typesQuery := r.URL.Query().Get("type")
		trimmedQuery := strings.TrimSpace(typesQuery)
		if len(trimmedQuery) != 0 {
			typesQuerySlice = strings.Split(trimmedQuery, " ")
		}

		noQuery := len(typesQuerySlice) == 0
		isBondTx := contains(typesQuerySlice, "bond")
		isUnbondTx := contains(typesQuerySlice, "unbond")

		var (
			txs     []*sdk.SearchTxsResult
			actions []string
		)

		switch {
		case isBondTx:
			actions = append(actions, types.MsgBond{}.Type())

		case isUnbondTx:
			actions = append(actions, types.MsgUnbond{}.Type())

		case noQuery:
			actions = append(actions, types.MsgBond{}.Type())
			actions = append(actions, types.MsgUnbond{}.Type())

		default:
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for _, action := range actions {
			foundTxs, errQuery := queryTxs(cliCtx, action, depositorAddr)
			if errQuery != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, errQuery.Error())
			}
			txs = append(txs, foundTxs)
		}

		res, err := cliCtx.Codec.MarshalJSON(txs)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}

// HTTP request handler to query an unbonding-deposit
func unbondingDepositHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryBonds(cliCtx, "custom/governors/unbondingDeposit")
}

// HTTP request handler to query a deposit
func depositHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryBonds(cliCtx, fmt.Sprintf("custom/%s/%s", types.QuerierRoute, bonds.QueryDeposit))
}

// HTTP request handler to query the pool information
func poolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/governors/pool", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query the params values
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/governors/parameters", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
