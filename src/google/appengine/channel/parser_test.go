package channel

import (
	"log"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	s := `[[3,["c",['8A6DD7714F2813E7',["me","me@talkgadget.google.com",0]
]]
]
,[4,["c",['8A6DD7714F2813E7',["cfj","w01q6amlqu69kj4thga19v9gvna4b8a5t5ilmdl9tq6cjdht3stjgl07jhj95o1hpd1lcg31oj4bulvruu1nknf0iiuu0zaez@guest.talk.google.com/appengine-channelC440E16B"]
]]
]
,[5,["c",['8A6DD7714F2813E7',["ru","go-rpi@appspot.com",,,0,0,0,0,3,,,[]
,1,4,0,[]
,,,,0,0,,0,0,1]
]]
]
,[6,["c",['8A6DD7714F2813E7',["vc","go-rpi@appspot.com",,"7b8ac4c3c5b468f9ceabd8eea8a9df61574ad997"]
]]
]
,[7,["c",['8A6DD7714F2813E7',["sus",[["defaultMuvcRoom",""]
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
,[8,["c",['8A6DD7714F2813E7',["sst",1]
]]
]
,[9,["c",['8A6DD7714F2813E7',["p","go-rpi@appspot.com/bot",0,5,,[]
,0,8957635240165547]
]]
]
]
`

	s = `[[6,["c",['D37147C04F2813E7',["ru","go-rpi@appspot.com",,,0,0,0,0,3,,,[]
,1,4,0,[]
,,,,0,0,,0,0,1]
]]
]`
	// s = `[[7,["c",['C3020D5B4F2813E7',["sus",[["defaultMuvcRoom",""]
	// ,["displayOob","true"]
	// ]
	// ]
	// ]]
	// ]
	// ]`

	// s = `[[14,["noop"]
	// ]
	// ]`

	arr := strings.Split(s[1:len(s)-2], "]\n]\n,")
	log.Println(arr[0])
	parser := &Parser{}
	root := parser.Parse([]byte(arr[0] + "]]"))
	log.Println(root.ToString())
}
