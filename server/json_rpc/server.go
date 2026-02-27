package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/osmosis-labs/osmosis/v27/server/api/eth"
	"github.com/osmosis-labs/osmosis/v27/x/evm/keeper"
)

// ContextProvider is a function that returns the current SDK context.
type ContextProvider func() sdk.Context

// Server is the JSON-RPC server for EVM compatibility.
type Server struct {
	config     Config
	router     *mux.Router
	ethAPI     *eth.API
	httpServer *http.Server
	getContext ContextProvider
}

// Config holds the JSON-RPC server configuration.
type Config struct {
	// Enable enables the JSON-RPC server.
	Enable bool `mapstructure:"enable"`
	// Address is the HTTP server address (default: "localhost:8545").
	Address string `mapstructure:"address"`
	// WSAPI defines the WebSocket API namespaces to enable.
	WsAddress string `mapstructure:"ws-address"`
	// API defines the HTTP API namespaces to enable.
	API string `mapstructure:"api"`
}

// DefaultConfig returns the default JSON-RPC server configuration.
func DefaultConfig() Config {
	return Config{
		Enable:    false,
		Address:   "localhost:8545",
		WsAddress: "localhost:8546",
		API:       "eth,net,web3",
	}
}

// NewServer creates a new JSON-RPC server.
// The contextProvider function should return the current SDK context for processing requests.
func NewServer(cfg Config, evmKeeper *keeper.Keeper, contextProvider ContextProvider) (*Server, error) {
	if !cfg.Enable {
		return nil, nil
	}

	router := mux.NewRouter()
	ethAPI := eth.NewAPI(evmKeeper)

	srv := &Server{
		config:     cfg,
		router:     router,
		ethAPI:     ethAPI,
		getContext: contextProvider,
	}

	// Register JSON-RPC handler
	router.HandleFunc("/", srv.handleJSONRPC).Methods("POST")

	return srv, nil
}

// Start starts the JSON-RPC HTTP server.
func (s *Server) Start() error {
	if s == nil {
		return nil
	}

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(s.router)

	s.httpServer = &http.Server{
		Addr:    s.config.Address,
		Handler: handler,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("JSON-RPC server error: %v\n", err)
		}
	}()

	fmt.Printf("JSON-RPC server started on %s\n", s.config.Address)
	return nil
}

// Stop stops the JSON-RPC server.
func (s *Server) Stop() error {
	if s == nil || s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(context.Background())
}

// handleJSONRPC handles JSON-RPC 2.0 requests.
func (s *Server) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nil, -32700, "Parse error", nil)
		return
	}

	// Get the current SDK context for processing the request
	ctx := s.getContext()

	// Route to appropriate handler based on method
	var result interface{}
	var err error

	switch req.Method {
	case "eth_chainId":
		result, err = s.ethAPI.ChainId(ctx)
	case "eth_blockNumber":
		result, err = s.ethAPI.BlockNumber(ctx)
	case "eth_getBalance":
		if len(req.Params) < 2 {
			writeError(w, req.ID, -32602, "Invalid params", nil)
			return
		}
		address, _ := req.Params[0].(string)
		blockNum, _ := req.Params[1].(string)
		result, err = s.ethAPI.GetBalance(ctx, address, blockNum)
	case "eth_getTransactionCount":
		if len(req.Params) < 2 {
			writeError(w, req.ID, -32602, "Invalid params", nil)
			return
		}
		address, _ := req.Params[0].(string)
		blockNum, _ := req.Params[1].(string)
		result, err = s.ethAPI.GetTransactionCount(ctx, address, blockNum)
	case "eth_getCode":
		if len(req.Params) < 2 {
			writeError(w, req.ID, -32602, "Invalid params", nil)
			return
		}
		address, _ := req.Params[0].(string)
		blockNum, _ := req.Params[1].(string)
		result, err = s.ethAPI.GetCode(ctx, address, blockNum)
	case "eth_call":
		if len(req.Params) < 2 {
			writeError(w, req.ID, -32602, "Invalid params", nil)
			return
		}
		// callArgs, _ := req.Params[0].(map[string]interface{})
		// blockNum, _ := req.Params[1].(string)
		result = "0x" // Placeholder
	case "eth_sendRawTransaction":
		if len(req.Params) < 1 {
			writeError(w, req.ID, -32602, "Invalid params", nil)
			return
		}
		txData, _ := req.Params[0].(string)
		result, err = s.ethAPI.SendRawTransaction(txData)
	case "eth_gasPrice":
		result, err = s.ethAPI.GasPrice()
	case "eth_estimateGas":
		result = "0x5208" // 21000 gas (placeholder)
	case "net_version":
		result, err = s.ethAPI.ChainId(ctx)
	case "web3_clientVersion":
		result = "Symphony/v1.0.0"
	default:
		writeError(w, req.ID, -32601, "Method not found", nil)
		return
	}

	if err != nil {
		writeError(w, req.ID, -32000, err.Error(), nil)
		return
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// JSONRPCRequest represents a JSON-RPC 2.0 request.
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response.
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error.
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// writeError writes a JSON-RPC error response.
func writeError(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // JSON-RPC errors are still 200 OK
	json.NewEncoder(w).Encode(response)
}
