package main

import (
	"google/appengine/channel"
	"log"
	"strings"
)

func main() {
	s := `[[3,["c",['8458A4D24F2813E7',["me","me@talkgadget.google.com",0]
]]
]
,[4,["c",['8458A4D24F2813E7',["cfj","w01q6amlqu69kj4thga19v9gvna4b8a5t5ilmdl9tq6cjdht3stjgl07jhj95o1hpd1lcg31oj4bulvruu1nknf0iiuu0zaez@guest.talk.google.com/appengine-channel3D23A2AA"]
]]
]
,[5,["c",['8458A4D24F2813E7',["sus",[["defaultMuvcRoom",""]
,["displayOob","true"]
,["displayRemotingOob","true"]
,["fullscreenBrowserAlways","true"]
,["gus",""]
,["hangoutsAbuseRecordingTermsAgreed","false"]
,["hangoutsCapturePromoAcked","false"]
,["hangoutsCapturePromoDialogAcked","false"]
,["hangoutsChatExpanded","false"]
,["hangoutsOnAirTermsAgreed","false"]
,["hangoutsOnAirV2NewCountryPromoAcked","false"]
,["hangoutsOnAirV2PromoAcked","false"]
,["hangoutsWaverlyHelpDisplayed","false"]
,["hangoutsWaverlySidebarCollapsed","false"]
,["hasSentMessage","false"]
,["highResolutionEnabled","false"]
,["introPromo","false"]
,["jmiQuality","true"]
,["lastFriendSuggestionPrompt","0"]
,["lastHatsSurvey","0"]
,["licensing",""]
,["moleAvatars","false"]
,["muvcRoomHistory",""]
,["namedRoomHistory",""]
,["pluginHelpPaneState","ok"]
,["rosterSortOption","MOST_POPULAR"]
,["showChatUpdated","true"]
,["showFriendSuggestionPrompts","true"]
,["signedOut","false"]
,["skipHangoutsExtrasPromo","false"]
,["soundEnabled","true"]
,["troubleshootingCaseId",""]
,["troubleshootingIssues",""]
,["troubleshootingStartTime","0"]
,["useChatCircles",""]
,["userCircleSetting",""]
]
]
]]
]
]
`
	s = `[5,["c",['8458A4D24F2813E7',["sus",[["defaultMuvcRoom",""]
,["displayOob","true"]
,["displayRemotingOob","true"]
,["fullscreenBrowserAlways","true"]
,["gus",""]
,["hangoutsAbuseRecordingTermsAgreed","false"]
,["hangoutsCapturePromoAcked","false"]
,["hangoutsCapturePromoDialogAcked","false"]
,["hangoutsChatExpanded","false"]
,["hangoutsOnAirTermsAgreed","false"]
,["hangoutsOnAirV2NewCountryPromoAcked","false"]
,["hangoutsOnAirV2PromoAcked","false"]
,["hangoutsWaverlyHelpDisplayed","false"]
,["hangoutsWaverlySidebarCollapsed","false"]
,["hasSentMessage","false"]
,["highResolutionEnabled","false"]
,["introPromo","false"]
,["jmiQuality","true"]
,["lastFriendSuggestionPrompt","0"]
,["lastHatsSurvey","0"]
,["licensing",""]
,["moleAvatars","false"]
,["muvcRoomHistory",""]
,["namedRoomHistory",""]
,["pluginHelpPaneState","ok"]
,["rosterSortOption","MOST_POPULAR"]
,["showChatUpdated","true"]
,["showFriendSuggestionPrompts","true"]
,["signedOut","false"]
,["skipHangoutsExtrasPromo","false"]
,["soundEnabled","true"]
,["troubleshootingCaseId",""]
,["troubleshootingIssues",""]
,["troubleshootingStartTime","0"]
,["useChatCircles",""]
,["userCircleSetting",""]
]]
]`
	parse([]byte(s), 0)
	log.Println(parent.Child.Child.Val)
}

type Element struct {
	Val   string
	Start int
	End   int
}

type KeyedElement struct {
	Key   string
	Val   *[]Element
	Child *KeyedElement
}

var parent *KeyedElement
var arr = &[]Element{}

func parse(bytes []byte, pos1 int) int {
	key := ""
	for i := pos1; i < len(bytes); i++ {
		b := bytes[i]
		switch b {
		case '[':
			key = strings.Trim(string(bytes[pos1:i]), "\"',")
			pos1 = i
			i = parse(bytes, i+1)
		case ']':
			e := Element{
				Start: pos1,
				End:   i,
				Val:   string(bytes[pos1:i]),
			}
			if key != "" {
				parent = &KeyedElement{
					Key:   key,
					Val:   arr,
					Child: parent,
				}
				arr = &[]Element{}
			} else {
				tmp := append(*arr, e)
				arr = &tmp
			}
			return i
		}
	}
	return len(bytes)
}

func main3() {
	stop_chan := make(chan bool)

	channel := channel.NewChannel("http://app.myalert.info/online_get_token?id=", "pi")
	socket := channel.Open()
	socket.OnOpened = func() {
		log.Println("socket opened!")
	}

	socket.OnClose = func() {
		log.Println("socket closed!")
		stop_chan <- true
	}

	socket.OnMessage = func(msg string) {
		// var v [][]interface{}
		// json.Unmarshal([]byte(msg), &v)
		log.Println(msg)
	}

	socket.OnError = func(err error) {
		log.Println("error:", err)
	}

	<-stop_chan
}
