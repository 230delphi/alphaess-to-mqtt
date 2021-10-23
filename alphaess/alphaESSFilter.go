package alphaess

import (
	"bytes"
	"fmt"
	"io"
)

const CLIENT = 0
const SERVER = 1
const CLIENTORSERVER = 2
const INJECT = 0
const DROP = 1
const RESPOND = 2

// TODO unit tests for filters

// actionType is a definition of possible actions that could be taken on a message across a proxy connection.
// 	name can be an arbitrary description of the step
// 	actionCode defines the type of action: INJECT(0)|DROP(1)|RESPOND(2)
// 	from denotes the source of message on which we respond: CLIENT(0)|SERVER(1)|CLIENTORSERVER(2)
//	message is the search string to check
// 	response is the message to send in the case of RESPOND actionCode
type actionType struct {
	name       string
	actionCode int    //INJECT|DROP|RESPOND
	from       int    //Client(0)|Server(1)|CLIENTORSERVER(2)
	message    []byte //message to insert or search string
	response   []byte //response to request for case of RESPOND
}

// conversationType is the statemachine definition of actionTypes required for a particular client/server conversation
//	indexOfNextAction is the 0 based id of the next action in the arrays actions.
//	actions is an array of actionTypes in order of execution.
//	lastUpdate is the epoch time of last update. might allow for expiry of broken conversation.s
type conversationType struct {
	indexOfNextAction int
	actions           []actionType
	lastUpdate        int
}

// MessageFilter is an interface to enable to update, ignore or inject new messages into a stream
type MessageFilter interface {
	// FilterMessages - return original or modified content, or nil to block forwarding
	FilterMessages([]byte, int) (message []byte, response []byte)
	// InjectMessage - return message to inject to stream or nil to ignore
	InjectMessage(io.Writer, []byte)
	getName() string
}

// PassFilter is a minimal do-nothing implementation.
type PassFilter struct {
	name string
}

func (into *PassFilter) FilterMessages(buf []byte, count int) (message []byte, response []byte) {
	return buf, nil
}
func (into *PassFilter) InjectMessage(dst io.Writer, lastMsg []byte) {
	//DebugLog("Pass Inject")
}
func (into *PassFilter) getName() string {
	return into.name
}

// ServerFilter processes messages FROM the server.
type ServerFilter struct {
	name string
}

func (into *ServerFilter) getName() string {
	return into.name
}

func (into *ServerFilter) FilterMessages(buf []byte, count int) (message []byte, response []byte) {
	return buf, nil
}

func (into *ServerFilter) InjectMessage(dst io.Writer, lastMsg []byte) {
	activeConversations := getActiveConversations(SERVER, INJECT)
	if len(activeConversations) > 0 {
		DebugLog("Found Server InjectMessage")
		// can only do one conversation at a time:
		if len(lastMsg) < 35 {
			InjectBytes(dst, activeConversations[0].actions[activeConversations[0].indexOfNextAction].message)
			setNextConversation(activeConversations[0])
			DebugLog(fmt.Sprint("TEMP Inject Write Intent successful. next index:", activeConversations[0].indexOfNextAction))
		}
	}
}

// ClientFilter processes messages FROM the Client.
type ClientFilter struct {
	name string
}

func (into *ClientFilter) getName() string {
	return into.name
}
func (into *ClientFilter) FilterMessages(buf []byte, count int) (rq []byte, rs []byte) {
	activeDROPConversations := getActiveConversations(CLIENT, DROP)
	activeREPSONDConversations := getActiveConversations(CLIENT, RESPOND)
	if len(activeDROPConversations) > 0 {
		DebugLog("Found Client DROP FilterMessages")
		for _, s := range activeDROPConversations {
			if bytes.Contains(buf, s.actions[s.indexOfNextAction].message) {
				//drop message
				buf = nil
				setNextConversation(s)
				break
			}
		}
	} else if len(activeREPSONDConversations) > 0 {
		for _, s := range activeREPSONDConversations {
			if bytes.Contains(buf, s.actions[s.indexOfNextAction].message) {
				//drop message
				rq = nil
				rs = s.actions[s.indexOfNextAction].response
				setNextConversation(s)
				DebugLog(fmt.Sprint("Found Client RESPOND FilterMessages for RESPONSE Intent shared. next index:", s.indexOfNextAction))
				break
			}
		}
	}
	if (rq == nil) && (rs == nil) {
		//nothing to apply.
		rq = buf
		rs = nil
	}
	return rq, rs
}
func (into *ClientFilter) InjectMessage(dst io.Writer, lastMessage []byte) {
	// required func, right now there is no need to inject new messages from Client TO server.
	//DebugLog("Client Inject")
}

// InjectBytes will put the specified myMsgBytes into the dst stream
func InjectBytes(dst io.Writer, myMsgBytes []byte) {
	clientStr := string(myMsgBytes)
	//DebugLog("InjectBytes() Writing: " + string(myMsgBytes))
	mutex.Lock()
	writen, err := dst.Write(myMsgBytes)
	mutex.Unlock()
	parseAndDebugMessage("InjectBytes()", myMsgBytes)
	if err != nil {
		ErrorLog("Error writing message: " + clientStr)
		ExceptionLog(err, "InjectBytes()")
	}
	if writen == len(myMsgBytes) {
		DebugLog("Inject Write successful.")
	}
}

// getActiveConversations gets the conversations a actionCode Type from a particular source.
func getActiveConversations(from int, actionCode int) (activeConversations []*conversationType) {
	if gActiveConversations != nil {
		for _, s := range gActiveConversations {
			// TODO if s.lastUpdate < 10 minutes
			if s.indexOfNextAction > -1 && s.indexOfNextAction < len(s.actions) {
				action := s.actions[s.indexOfNextAction]
				if action.actionCode == actionCode && (action.from == from || action.from == CLIENTORSERVER) {
					activeConversations = append(activeConversations, s)
				}
			}
		}
	}
	return activeConversations
}

// setNextConversation is a simple utility method to update indexOfNextAction to the next action in the conversation.
//	conversations are expired by setting indexOfNextAction to -1.
func setNextConversation(c *conversationType) {
	if c.indexOfNextAction < (len(c.actions) - 1) {
		c.indexOfNextAction++
	} else {
		c.indexOfNextAction = -1
	}
}
