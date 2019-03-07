package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientRest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	csdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	sdkTypes "github.com/ironman0x7b2/sentinel-sdk/types"
	"github.com/ironman0x7b2/sentinel-sdk/x/vpn"
)

type msgRegisterNode struct {
	BaseReq      rest.BaseReq       `json:"base_req"`
	AmountToLock string             `json:"amount_to_lock"`
	PricesPerGB  string             `json:"prices_per_gb"`
	NetSpeed     sdkTypes.Bandwidth `json:"net_speed"`
	APIPort      uint16             `json:"api_port"`
	Encryption   string             `json:"encryption"`
	Version      string             `json:"version"`
	NodeType     string             `json:"node_type"`
}

func registerNodeHandlerFunc(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req msgRegisterNode

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := csdkTypes.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		amountToLock, err := csdkTypes.ParseCoin(req.AmountToLock)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		pricesPerGB, err := csdkTypes.ParseCoins(req.PricesPerGB)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := vpn.NewMsgRegisterNode(fromAddress,
			amountToLock, pricesPerGB, req.NetSpeed.Upload, req.NetSpeed.Download,
			req.APIPort, req.Encryption, req.NodeType, req.Version)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientRest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []csdkTypes.Msg{msg})
	}
}
