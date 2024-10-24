package kvstore

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/DeAI-Artist/Linkis/abci/example/kvstore/txs"
	cryptoenc "github.com/DeAI-Artist/Linkis/crypto/encoding"
	pc "github.com/DeAI-Artist/Linkis/proto/tendermint/crypto"
	"strconv"
	"strings"

	dbm "github.com/tendermint/tm-db"

	"github.com/DeAI-Artist/Linkis/abci/example/code"
	"github.com/DeAI-Artist/Linkis/abci/types"
	"github.com/DeAI-Artist/Linkis/libs/log"
	"github.com/DeAI-Artist/Linkis/version"
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")

	ProtocolVersion uint64 = 0x1
)

type State struct {
	db                   dbm.DB
	Size                 int64            `json:"size"`
	Height               int64            `json:"height"`
	AppHash              []byte           `json:"app_hash"`
	MinerActivityRecords MinerWorkRecords `json:"miner_activity_records"`
}

func loadState(db dbm.DB) State {
	var state State
	state.db = db
	stateBytes, err := db.Get(stateKey)
	if err != nil {
		panic(err)
	}
	if len(stateBytes) == 0 {
		return state
	}
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		panic(err)
	}
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	err = state.db.Set(stateKey, stateBytes)
	if err != nil {
		panic(err)
	}
}

//---------------------------------------------------

var _ types.Application = (*Application)(nil)

type Application struct {
	types.BaseApplication

	state        State
	RetainBlocks int64 // blocks to retain after commit (via ResponseCommit.RetainHeight)
	// validator set
	ValUpdates []types.ValidatorUpdate

	valAddrToPubKeyMap map[string]pc.PublicKey

	logger log.Logger
}

func NewApplication(dbDir string) *Application {
	name := "kvstore"
	db, err := dbm.NewGoLevelDB(name, dbDir)
	if err != nil {
		panic(err)
	}
	state := loadState(db)

	return &Application{
		state:              state,
		valAddrToPubKeyMap: make(map[string]pc.PublicKey),
		logger:             log.NewNopLogger(),
	}
}

func (app *Application) SetLogger(l log.Logger) {
	app.logger = l
}

func (app *Application) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return app.SetOption(req)
}

func (app *Application) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       ProtocolVersion,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.AppHash,
	}
}

// Track the block hash and header information
func (app *Application) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	// reset valset changes
	app.ValUpdates = make([]types.ValidatorUpdate, 0)

	// Punish validators who committed equivocation.
	for _, ev := range req.ByzantineValidators {
		if ev.Type == types.EvidenceType_DUPLICATE_VOTE {
			addr := string(ev.Validator.Address)
			if pubKey, ok := app.valAddrToPubKeyMap[addr]; ok {
				app.updateValidator(types.ValidatorUpdate{
					PubKey: pubKey,
					Power:  ev.Validator.Power - 1,
				})
				app.logger.Info("Decreased val power by 1 because of the equivocation",
					"val", addr)
			} else {
				app.logger.Error("Wanted to punish val, but can't find it",
					"val", addr)
			}
		}
	}

	return types.ResponseBeginBlock{}
}

// Update the validator set
func (app *Application) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

func (app *Application) ListSnapshots(
	req types.RequestListSnapshots) types.ResponseListSnapshots {
	return types.ResponseListSnapshots{}
}

func (app *Application) LoadSnapshotChunk(
	req types.RequestLoadSnapshotChunk) types.ResponseLoadSnapshotChunk {
	return types.ResponseLoadSnapshotChunk{}
}

func (app *Application) OfferSnapshot(
	req types.RequestOfferSnapshot) types.ResponseOfferSnapshot {
	return types.ResponseOfferSnapshot{Result: types.ResponseOfferSnapshot_ABORT}
}

func (app *Application) ApplySnapshotChunk(
	req types.RequestApplySnapshotChunk) types.ResponseApplySnapshotChunk {
	return types.ResponseApplySnapshotChunk{Result: types.ResponseApplySnapshotChunk_ABORT}
}

//---------------------------------------------
// update validators

func (app *Application) Validators() (validators []types.ValidatorUpdate) {
	itr, err := app.state.db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	for ; itr.Valid(); itr.Next() {
		if isValidatorTx(itr.Key()) {
			validator := new(types.ValidatorUpdate)
			err := types.ReadMessage(bytes.NewBuffer(itr.Value()), validator)
			if err != nil {
				panic(err)
			}
			validators = append(validators, *validator)
		}
	}
	if err = itr.Error(); err != nil {
		panic(err)
	}
	return
}

// format is "val:pubkey!power"
// pubkey is a base64-encoded 32-byte ed25519 key
func (app *Application) execValidatorTx(tx []byte) types.ResponseDeliverTx {
	tx = tx[len(ValidatorSetChangePrefix):]

	//  get the pubkey and power
	pubKeyAndPower := strings.Split(string(tx), "!")
	if len(pubKeyAndPower) != 2 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected 'pubkey!power'. Got %v", pubKeyAndPower)}
	}
	pubkeyS, powerS := pubKeyAndPower[0], pubKeyAndPower[1]

	// decode the pubkey
	pubkey, err := base64.StdEncoding.DecodeString(pubkeyS)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Pubkey (%s) is invalid base64", pubkeyS)}
	}

	// decode the power
	power, err := strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Power (%s) is not an int", powerS)}
	}

	// update
	return app.updateValidator(types.UpdateValidator(pubkey, power, ""))
}

// tx is either "key=value" or just arbitrary bytes
func (app *Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {

	tx := string(req.Tx)

	msg, sig, _ := txs.DecodeMessageAndSignature(tx)
	msgBytes := []byte(msg)
	sigBytes, _ := hex.DecodeString(sig)
	pubKey, _ := txs.RecoverPubKey(txs.HashPersonalMessage(msgBytes), sigBytes)
	senderAddr := txs.AddressFromPublicKey(pubKey)

	var transaction txs.Transaction
	err := transaction.FromString(string(req.Tx))
	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError, GasWanted: 0, Log: err.Error()}
	}

	switch transaction.Msg.Type {
	case txs.ClientRegistrationType:
		err := app.handleClientRegistration(senderAddr, transaction.Msg)
		if err != nil {
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: err.Error()}
		}
		app.state.Size++

	case txs.ServiceRequestType:
		err := app.handleServiceRequest(senderAddr, transaction.Msg)
		if err != nil {
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: err.Error()}
		}
		app.state.Size++

	case txs.ClientRatingMsgType:
		err := app.handleClientRating(senderAddr, transaction.Msg)
		if err != nil {
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: err.Error()}
		}
		app.state.Size++

	case txs.MinerRegistrationType:
		err := app.handleMinerRegistration(senderAddr, transaction.Msg)
		if err != nil {
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: err.Error()}
		}
		app.state.Size++

	case txs.MinerServiceDoneType:
		err := app.handleMinerServiceDone(senderAddr, transaction.Msg)
		if err != nil {
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: err.Error()}
		}
		app.state.Size++

	case txs.MinerStatusUpdateType:
		err := app.handleMinerStatusUpdate(senderAddr, transaction.Msg)
		if err != nil {
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: err.Error()}
		}
		app.state.Size++

	case txs.MinerRewardClaimType:
		err := app.handleMinerRewardClaim(senderAddr, transaction.Msg)
		if err != nil {
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: err.Error()}
		}
		app.state.Size++

	case txs.MinerServiceStartingType:
		err := app.handleMinerServiceStarting(senderAddr, transaction.Msg)
		if err != nil {
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: err.Error()}
		}
		app.state.Size++

	default:
		return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, GasWanted: 0, Log: "Unknown message type"}
	}

	events := []types.Event{
		{
			Type: "app",
			Attributes: []types.EventAttribute{
				{Key: []byte("transaction sender"), Value: []byte(senderAddr), Index: true},
			},
		},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: events}
}

func (app *Application) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	message, signature, err := txs.DecodeMessageAndSignature(string(req.Tx))
	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError, GasWanted: 0}
	}

	// Hash the message using personal_sign spec
	hashed := txs.HashPersonalMessage([]byte(message))
	sigBytes, err := txs.HexToBytes(signature)
	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError, GasWanted: 0}
	}
	_, err = txs.RecoverPubKey(hashed, sigBytes)
	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeUnauthorized, GasWanted: 0}
	}

	// Optionally, validate the public key or other aspects of the transaction
	// For this example, we assume success if we reach this point
	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 0}
}

func (app *Application) Commit() types.ResponseCommit {
	hasher := sha256.New()

	// Serialize and hash critical parts of the state
	// This approach assumes you keep track of changes rather than iterate through the entire DB
	stateData, err := json.Marshal(struct {
		Size   int64  `json:"size"`
		Height int64  `json:"height"`
		Hash   []byte `json:"app_hash"`
	}{
		Size:   app.state.Size,
		Height: app.state.Height,
		Hash:   app.state.AppHash,
	})
	if err != nil {
		panic(err) // handle the error more gracefully in production
	}
	hasher.Write(stateData)

	// Optionally include miner activity records
	activityData, err := json.Marshal(app.state.MinerActivityRecords)
	if err != nil {
		panic(err) // handle the error appropriately
	}
	hasher.Write(activityData)

	// Compute the new app hash
	appHash := hasher.Sum(nil)

	// Update the state with the new hash and increment the height
	app.state.AppHash = appHash
	app.state.Height++
	saveState(app.state) // Save the updated state to the database or config

	// Prepare the commit response
	resp := types.ResponseCommit{Data: appHash}
	if app.RetainBlocks > 0 && app.state.Height >= app.RetainBlocks {
		resp.RetainHeight = app.state.Height - app.RetainBlocks + 1
	}
	return resp
}

// Returns an associated value or nil if missing.
func (app *Application) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	if reqQuery.Prove {
		value, err := app.state.db.Get(reqQuery.Data)
		if err != nil {
			panic(err)
		}
		if value == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		resQuery.Height = app.state.Height

		return
	}

	resQuery.Key = reqQuery.Data
	value, err := app.state.db.Get(reqQuery.Data)
	if err != nil {
		panic(err)
	}
	if value == nil {
		resQuery.Log = "does not exist"
	} else {
		resQuery.Log = "exists"
	}
	resQuery.Value = value
	resQuery.Height = app.state.Height

	return resQuery
}

// handleClientRegistration processes a client registration message.
func (app *Application) handleClientRegistration(sender string, msg txs.Message) error {
	// Assume the Content of the message is a JSON with client information.
	// You would need to define a ClientRegistrationInfo struct according to your application's needs.
	cr, err := msg.DecodeContent()
	if err != nil {
		return fmt.Errorf("error unmarshaling client registration info: %v", err)
	}

	crm, ok := cr.(txs.ClientRegistrationMsg)
	if !ok {
		return fmt.Errorf("type assertion to ClientRegistrationMsg failed")
	}

	clientInfo := ClientInfo{
		Name:  crm.ClientName,
		Power: 10,
	}

	// Perform registration logic, e.g., storing info in a database.
	err = StoreClientInfo(app.state.db, sender, clientInfo)
	if err != nil {
		return fmt.Errorf("StoreClientInfo failed: %v", err)
	}

	// Additional logic here, such as sending confirmation emails, logging, etc.

	return nil
}

// handleMinerRegistration processes a miner registration message.
func (app *Application) handleMinerRegistration(sender string, msg txs.Message) error {
	// Decode the content from the Message struct, expecting miner registration information.
	mr, err := msg.DecodeContent()
	if err != nil {
		return fmt.Errorf("error decoding miner registration info: %v", err)
	}

	// Type assertion to MinerRegistrationMsg
	mrm, ok := mr.(txs.MinerRegistrationMsg)
	if !ok {
		return fmt.Errorf("type assertion to MinerRegistrationMsg failed")
	}

	// Construct MinerInfo from decoded content
	minerInfo := MinerInfo{
		Name:          mrm.MinerName,
		Power:         10, // Assume a fixed value or derive it from another source if necessary
		ServiceTypes:  mrm.ServiceTypes,
		IP:            mrm.IP,
		InitialStatus: mrm.Status,
	}

	// Store miner information in the database
	err = StoreMinerInfo(app.state.db, sender, minerInfo)
	if err != nil {
		return fmt.Errorf("StoreMinerInfo failed: %v", err)
	}

	// Register the miner in service type mappings
	err = RegisterMiner(app.state.db, minerInfo, sender)
	if err != nil {
		return fmt.Errorf("RegisterMiner failed: %v", err)
	}

	err = AddOrUpdateMinerStatus(app.state.db, sender, minerInfo.InitialStatus)
	if err != nil {
		return fmt.Errorf("AddOrUpdateMinerStatus failed: %v", err)
	}

	// Optionally, additional logic such as logging the registration, notifying other systems, etc.
	fmt.Printf("Registered new miner: %s, IP: %s\n", minerInfo.Name, minerInfo.IP)

	return nil
}

// handleMinerStatusUpdate processes a status update message for a miner.
func (app *Application) handleMinerStatusUpdate(senderAddr string, msg txs.Message) error {
	// Decode the content from the Message struct, expecting a status update information.
	ms, err := msg.DecodeContent()
	if err != nil {
		return fmt.Errorf("error decoding miner status update info: %v", err)
	}

	// Type assertion to MinerStatusUpdateMsg
	msm, ok := ms.(txs.MinerStatusUpdateMsg)
	if !ok {
		return fmt.Errorf("type assertion to MinerStatusUpdateMsg failed")
	}

	// Retrieve the existing miner info from the database
	minerInfo, err := GetMinerInfo(app.state.db, senderAddr)
	if err != nil {
		return fmt.Errorf("failed to get miner info: %v", err)
	}

	// Update service types: Remove first, then add
	currentServiceTypes := make(map[uint64]bool)
	for _, st := range minerInfo.ServiceTypes {
		currentServiceTypes[st] = true
	}

	// Remove service types
	for _, st := range msm.RemoveServiceTypes {
		delete(currentServiceTypes, st)
		// Remove the miner from the service type mapping
		if err := RemoveMinerFromServiceTypeMapping(app.state.db, st, senderAddr); err != nil {
			return fmt.Errorf("failed to remove miner from service type mapping: %v", err)
		}
	}

	// Add new service types
	for _, st := range msm.AddServiceTypes {
		currentServiceTypes[st] = true
		// Add the miner to the service type mapping
		if err := AddMinerToServiceTypeMapping(app.state.db, st, senderAddr); err != nil {
			return fmt.Errorf("failed to add miner to service type mapping: %v", err)
		}
	}

	// Convert map back to slice for storage
	updatedServiceTypes := make([]uint64, 0, len(currentServiceTypes))
	for st := range currentServiceTypes {
		updatedServiceTypes = append(updatedServiceTypes, st)
	}
	minerInfo.ServiceTypes = updatedServiceTypes

	// Store the updated miner info back into the database
	if err := StoreMinerInfo(app.state.db, senderAddr, minerInfo); err != nil {
		return fmt.Errorf("failed to update miner info: %v", err)
	}

	// Update the miner's status in the database
	err = AddOrUpdateMinerStatus(app.state.db, senderAddr, msm.Status)
	if err != nil {
		return fmt.Errorf("AddOrUpdateMinerStatus failed: %v", err)
	}

	// Optionally, logging the status update.
	fmt.Printf("Updated miner status: Address=%s, New Status=%d\n", senderAddr, msm.Status)

	return nil
}

// handleMinerRewardClaim processes a miner reward claim message.
// Currently, this function does nothing and always returns nil.
func (app *Application) handleMinerRewardClaim(senderAddr string, msg txs.Message) error {
	// This function is intentionally left empty for future implementations
	// or for fulfilling interface requirements.

	// Log the operation if needed for debug purposes
	// fmt.Printf("Received reward claim from %s, but no action taken.\n", senderAddr)

	return nil
}

func (app *Application) handleClientRating(senderAddr string, msg txs.Message) error {
	// Decode the content from the Message struct, assuming it includes the miner's address and the rating.
	cr, err := msg.DecodeContent()
	if err != nil {
		return fmt.Errorf("error decoding client rating info: %v", err)
	}

	// Type assertion to ClientRatingMsg
	crm, ok := cr.(txs.ClientRatingMsg)
	if !ok {
		return fmt.Errorf("type assertion to ClientRatingMsg failed")
	}

	// Retrieve existing ratings from the database
	ratings, err := GetClientRating(app.state.db, crm.ReviewedMinerAddr)
	if err != nil {
		return fmt.Errorf("failed to get ratings for miner %s: %v", crm.ReviewedMinerAddr, err)
	}

	// Update the rating map with the new rating
	ratings[senderAddr] = uint8(crm.Rating)

	// Store the updated ratings back into the database
	err = StoreClientRating(app.state.db, crm.ReviewedMinerAddr, ratings)
	if err != nil {
		return fmt.Errorf("failed to update ratings for miner %s: %v", crm.ReviewedMinerAddr, err)
	}

	// Optionally, log the update
	fmt.Printf("Updated rating for miner: %s by client: %s to %d\n", crm.ReviewedMinerAddr, senderAddr, crm.Rating)

	return nil
}

// handleServiceRequest processes a service request from a client.
func (app *Application) handleServiceRequest(senderAddr string, msg txs.Message) error {
	// Decode the content from the Message struct, expecting a service request message.
	sr, err := msg.DecodeContent()
	if err != nil {
		return fmt.Errorf("error decoding service request info: %v", err)
	}

	// Type assertion to ServiceRequestMsg
	srm, ok := sr.(txs.ServiceRequestMsg)
	if !ok {
		return fmt.Errorf("type assertion to ServiceRequestMsg failed")
	}

	// Retrieve the current block height and hash from the application state
	currentHeight := app.state.Height // Assuming app.state has a BlockHeight field\
	currentHash := app.state.AppHash

	// Generate a unique service ID using hash of the client's address, metadata, and block height
	serviceID := GenerateHashForServiceInfo(senderAddr, srm.Meta, currentHeight)

	// Retrieve the list of miners registered for the specific service type from the database
	miners, err := GetMinersForServiceType(app.state.db, srm.ServiceID)
	if err != nil {
		return fmt.Errorf("failed to retrieve miners for service type %d: %v", srm.ServiceID, err)
	}
	if len(miners) == 0 {
		return fmt.Errorf("no miners registered for service type %d", srm.ServiceID)
	}

	// Select a pseudorandom miner based on the block height, app hash, and service ID
	selectedMiner := txs.SelectPseudorandomMiner(miners, currentHeight, currentHash, serviceID)

	// Create a JobInfo struct (details to be defined elsewhere)
	jobInfo := JobInfo{
		ServiceID:   serviceID,
		ServiceType: srm.ServiceID,
		ClientID:    senderAddr,
		JobStatus:   Registered,
	}

	// Store JobInfo in the database under the key derived from the miner's ID
	if err := StoreJobInfo(app.state.db, selectedMiner, jobInfo); err != nil {
		return fmt.Errorf("failed to store job info for miner ID '%s': %v", selectedMiner, err)
	}

	err = AddServiceRequest(app.state.db, serviceID, selectedMiner, currentHeight)
	if err != nil {
		return fmt.Errorf("failed to add service request: %v", err)
	}

	return nil
}

// handleServiceStarting processes the beginning of a service with conditions.
func (app *Application) handleMinerServiceStarting(senderAddr string, msg txs.Message) error {
	// Decode the content from the Message struct, expecting a service starting message.
	ss, err := msg.DecodeContent()
	if err != nil {
		return fmt.Errorf("error decoding service starting info: %v", err)
	}

	// Type assertion to ServiceStartingMsg
	ssm, ok := ss.(txs.ServiceStartingMsg)
	if !ok {
		return fmt.Errorf("type assertion to ServiceStartingMsg failed")
	}

	serviceID := ssm.ServiceID
	minerID := senderAddr
	blockOffset := ssm.MaxTimeoutBlock
	currentBlock := app.state.Height

	// Retrieve the job information by service ID
	jobInfo, err := GetJobInfoByServiceID(app.state.db, minerID, serviceID)
	if err != nil {
		return fmt.Errorf("failed to retrieve job info for ServiceID '%s': %v", serviceID, err)
	}

	// Check if the jobInfo is effectively empty
	if (JobInfo{}) == jobInfo {
		return fmt.Errorf("no valid job info found for ServiceID '%s'", serviceID)
	}

	jobInfo.TimeoutBlock = blockOffset + currentBlock
	jobInfo.JobStatus = Processing

	// Store the updated job info back into the database
	if err := StoreJobInfo(app.state.db, minerID, jobInfo); err != nil {
		return fmt.Errorf("failed to store updated job info for ServiceID '%s': %v", serviceID, err)
	}

	// Remove the service request associated with this service ID
	if err := RemoveServiceRequest(app.state.db, serviceID); err != nil {
		return fmt.Errorf("failed to remove service request for ServiceID '%s': %v", serviceID, err)
	}

	return nil
}

// handleMinerServiceDone processes a miner service done event.
func (app *Application) handleMinerServiceDone(senderAddr string, msg txs.Message) error {
	// Decode the content from the Message struct, expecting a miner service done message.
	msd, err := msg.DecodeContent()
	if err != nil {
		return fmt.Errorf("error decoding miner service done info: %v", err)
	}

	// Type assertion to MinerServiceDoneMsg
	msdm, ok := msd.(txs.MinerServiceDoneMsg)
	if !ok {
		return fmt.Errorf("type assertion to MinerServiceDoneMsg failed")
	}

	serivceType := msdm.ServiceType
	minerID := senderAddr

	// Increment the service type count in the work records
	currentHeight := app.state.Height // Assuming you track current block height in app.state
	if err := app.state.MinerActivityRecords.IncrementServiceType(currentHeight, minerID, serivceType); err != nil {
		return fmt.Errorf("failed to increment service type count: %v", err)
	}

	return nil
}

// add, update, or remove a validator
func (app *Application) updateValidator(v types.ValidatorUpdate) types.ResponseDeliverTx {
	pubkey, err := cryptoenc.PubKeyFromProto(v.PubKey)
	if err != nil {
		panic(fmt.Errorf("can't decode public key: %w", err))
	}
	key := []byte("val:" + string(pubkey.Bytes()))

	if v.Power == 0 {
		// remove validator
		hasKey, err := app.state.db.Has(key)
		if err != nil {
			panic(err)
		}
		if !hasKey {
			pubStr := base64.StdEncoding.EncodeToString(pubkey.Bytes())
			return types.ResponseDeliverTx{
				Code: code.CodeTypeUnauthorized,
				Log:  fmt.Sprintf("Cannot remove non-existent validator %s", pubStr)}
		}
		if err = app.state.db.Delete(key); err != nil {
			panic(err)
		}
		delete(app.valAddrToPubKeyMap, string(pubkey.Address()))
	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&v, value); err != nil {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Error encoding validator: %v", err)}
		}
		if err = app.state.db.Set(key, value.Bytes()); err != nil {
			panic(err)
		}
		app.valAddrToPubKeyMap[string(pubkey.Address())] = v.PubKey
	}

	// we only update the changes array if we successfully updated the tree
	app.ValUpdates = append(app.ValUpdates, v)

	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}
