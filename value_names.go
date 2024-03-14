package stagelinq

import "fmt"

const (
	ClientLibrarianDevicesControllerCurrentDevice         = "/Client/Librarian/DevicesController/CurrentDevice"
	ClientLibrarianDevicesControllerHasSDCardConnected    = "/Client/Librarian/DevicesController/HasSDCardConnected"
	ClientLibrarianDevicesControllerHasUsbDeviceConnected = "/Client/Librarian/DevicesController/HasUsbDeviceConnected"
	ClientPreferencesLayerA                               = "/Client/Preferences/LayerA"
	ClientPreferencesLayerB                               = "/Client/Preferences/LayerB"
	ClientPreferencesPlayer                               = "/Client/Preferences/Player"
	ClientPreferencesPlayerJogColorA                      = "/Client/Preferences/PlayerJogColorA"
	ClientPreferencesPlayerJogColorB                      = "/Client/Preferences/PlayerJogColorB"
	ClientPreferencesProfileApplicationPlayerColor1       = "/Client/Preferences/Profile/Application/PlayerColor1"
	ClientPreferencesProfileApplicationPlayerColor1A      = "/Client/Preferences/Profile/Application/PlayerColor1A"
	ClientPreferencesProfileApplicationPlayerColor1B      = "/Client/Preferences/Profile/Application/PlayerColor1B"
	ClientPreferencesProfileApplicationPlayerColor2       = "/Client/Preferences/Profile/Application/PlayerColor2"
	ClientPreferencesProfileApplicationPlayerColor2A      = "/Client/Preferences/Profile/Application/PlayerColor2A"
	ClientPreferencesProfileApplicationPlayerColor2B      = "/Client/Preferences/Profile/Application/PlayerColor2B"
	ClientPreferencesProfileApplicationPlayerColor3       = "/Client/Preferences/Profile/Application/PlayerColor3"
	ClientPreferencesProfileApplicationPlayerColor3A      = "/Client/Preferences/Profile/Application/PlayerColor3A"
	ClientPreferencesProfileApplicationPlayerColor3B      = "/Client/Preferences/Profile/Application/PlayerColor3B"
	ClientPreferencesProfileApplicationPlayerColor4       = "/Client/Preferences/Profile/Application/PlayerColor4"
	ClientPreferencesProfileApplicationPlayerColor4A      = "/Client/Preferences/Profile/Application/PlayerColor4A"
	ClientPreferencesProfileApplicationPlayerColor4B      = "/Client/Preferences/Profile/Application/PlayerColor4B"
	ClientPreferencesProfileApplicationSyncMode           = "/Client/Preferences/Profile/Application/SyncMode"
	EngineDeckCount                                       = "/Engine/DeckCount"
	EngineMasterMasterTempo                               = "/Engine/Master/MasterTempo"
	EngineSyncNetworkMasterStatus                         = "/Engine/Sync/Network/MasterStatus"
	GUIDecksDeckActiveDeck                                = "/GUI/Decks/Deck/ActiveDeck"
	GUIViewLayerLayerB                                    = "/GUI/ViewLayer/LayerB"
	MixerCH1faderPosition                                 = "/Mixer/CH1faderPosition"
	MixerCH2faderPosition                                 = "/Mixer/CH2faderPosition"
	MixerCH3faderPosition                                 = "/Mixer/CH3faderPosition"
	MixerCH4faderPosition                                 = "/Mixer/CH4faderPosition"
	MixerChannelAssignment1                               = "/Mixer/ChannelAssignment1"
	MixerChannelAssignment2                               = "/Mixer/ChannelAssignment2"
	MixerChannelAssignment3                               = "/Mixer/ChannelAssignment3"
	MixerChannelAssignment4                               = "/Mixer/ChannelAssignment4"
	MixerCrossfaderPosition                               = "/Mixer/CrossfaderPosition"
	MixerNumberOfChannels                                 = "/Mixer/NumberOfChannels"
)

// Deck 1 legacy variables
//
// Deprecated: Use values from [EngineDeck1] instead.
var (
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.CurrentBPM] method instead.
	EngineDeck1CurrentBPM = EngineDeck1.CurrentBPM()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.DeckIsMaster] method instead.
	EngineDeck1DeckIsMaster = EngineDeck1.DeckIsMaster()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.ExternalMixerVolume] method instead.
	EngineDeck1ExternalMixerVolume = EngineDeck1.ExternalMixerVolume()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.ExternalScratchWheelTouch] method instead.
	EngineDeck1ExternalScratchWheelTouch = EngineDeck1.ExternalScratchWheelTouch()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.PadsView] method instead.
	EngineDeck1PadsView = EngineDeck1.PadsView()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.Play] method instead.
	EngineDeck1Play = EngineDeck1.Play()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.PlayState] method instead.
	EngineDeck1PlayState = EngineDeck1.PlayState()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.PlayStatePath] method instead.
	EngineDeck1PlayStatePath = EngineDeck1.PlayStatePath()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.Speed] method instead.
	EngineDeck1Speed = EngineDeck1.Speed()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.SpeedNeutral] method instead.
	EngineDeck1SpeedNeutral = EngineDeck1.SpeedNeutral()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.SpeedOffsetDown] method instead.
	EngineDeck1SpeedOffsetDown = EngineDeck1.SpeedOffsetDown()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.SpeedOffsetUp] method instead.
	EngineDeck1SpeedOffsetUp = EngineDeck1.SpeedOffsetUp()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.SpeedRange] method instead.
	EngineDeck1SpeedRange = EngineDeck1.SpeedRange()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.SpeedState] method instead.
	EngineDeck1SpeedState = EngineDeck1.SpeedState()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.SyncMode] method instead.
	EngineDeck1SyncMode = EngineDeck1.SyncMode()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackArtistName] method instead.
	EngineDeck1TrackArtistName = EngineDeck1.TrackArtistName()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackBleep] method instead.
	EngineDeck1TrackBleep = EngineDeck1.TrackBleep()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackCuePosition] method instead.
	EngineDeck1TrackCuePosition = EngineDeck1.TrackCuePosition()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackCurrentBPM] method instead.
	EngineDeck1TrackCurrentBPM = EngineDeck1.TrackCurrentBPM()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackCurrentKeyIndex] method instead.
	EngineDeck1TrackCurrentKeyIndex = EngineDeck1.TrackCurrentKeyIndex()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackCurrentLoopInPosition] method instead.
	EngineDeck1TrackCurrentLoopInPosition = EngineDeck1.TrackCurrentLoopInPosition()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackCurrentLoopOutPosition] method instead.
	EngineDeck1TrackCurrentLoopOutPosition = EngineDeck1.TrackCurrentLoopOutPosition()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackCurrentLoopSizeInBeats] method instead.
	EngineDeck1TrackCurrentLoopSizeInBeats = EngineDeck1.TrackCurrentLoopSizeInBeats()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackKeyLock] method instead.
	EngineDeck1TrackKeyLock = EngineDeck1.TrackKeyLock()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopEnableState] method instead.
	EngineDeck1TrackLoopEnableState = EngineDeck1.TrackLoopEnableState()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopQuickLoop1] method instead.
	EngineDeck1TrackLoopQuickLoop1 = EngineDeck1.TrackLoopQuickLoop1()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopQuickLoop2] method instead.
	EngineDeck1TrackLoopQuickLoop2 = EngineDeck1.TrackLoopQuickLoop2()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopQuickLoop3] method instead.
	EngineDeck1TrackLoopQuickLoop3 = EngineDeck1.TrackLoopQuickLoop3()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopQuickLoop4] method instead.
	EngineDeck1TrackLoopQuickLoop4 = EngineDeck1.TrackLoopQuickLoop4()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopQuickLoop5] method instead.
	EngineDeck1TrackLoopQuickLoop5 = EngineDeck1.TrackLoopQuickLoop5()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopQuickLoop6] method instead.
	EngineDeck1TrackLoopQuickLoop6 = EngineDeck1.TrackLoopQuickLoop6()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopQuickLoop7] method instead.
	EngineDeck1TrackLoopQuickLoop7 = EngineDeck1.TrackLoopQuickLoop7()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackLoopQuickLoop8] method instead.
	EngineDeck1TrackLoopQuickLoop8 = EngineDeck1.TrackLoopQuickLoop8()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackPlayPauseLEDState] method instead.
	EngineDeck1TrackPlayPauseLEDState = EngineDeck1.TrackPlayPauseLEDState()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackSampleRate] method instead.
	EngineDeck1TrackSampleRate = EngineDeck1.TrackSampleRate()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackSongAnalyzed] method instead.
	EngineDeck1TrackSongAnalyzed = EngineDeck1.TrackSongAnalyzed()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackSongLoaded] method instead.
	EngineDeck1TrackSongLoaded = EngineDeck1.TrackSongLoaded()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackSongName] method instead.
	EngineDeck1TrackSongName = EngineDeck1.TrackSongName()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackSoundSwitchGuid] method instead.
	EngineDeck1TrackSoundSwitchGUID = EngineDeck1.TrackSoundSwitchGuid()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackTrackBytes] method instead.
	EngineDeck1TrackTrackBytes = EngineDeck1.TrackTrackBytes()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackTrackData] method instead.
	EngineDeck1TrackTrackData = EngineDeck1.TrackTrackData()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackTrackLength] method instead.
	EngineDeck1TrackTrackLength = EngineDeck1.TrackTrackLength()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackTrackName] method instead.
	EngineDeck1TrackTrackName = EngineDeck1.TrackTrackName()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackTrackNetworkPath] method instead.
	EngineDeck1TrackTrackNetworkPath = EngineDeck1.TrackTrackNetworkPath()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackTrackUri] method instead.
	EngineDeck1TrackTrackURI = EngineDeck1.TrackTrackUri()
	// Deprecated: Use [EngineDeck1] and its [DeckValueNames.TrackTrackWasPlayed] method instead.
	EngineDeck1TrackTrackWasPlayed = EngineDeck1.TrackTrackWasPlayed()
)

// Deck 2 legacy variables
//
// Deprecated: Use values from [EngineDeck2] instead.
var (
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.CurrentBPM] method instead.
	EngineDeck2CurrentBPM = EngineDeck2.CurrentBPM()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.DeckIsMaster] method instead.
	EngineDeck2DeckIsMaster = EngineDeck2.DeckIsMaster()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.ExternalMixerVolume] method instead.
	EngineDeck2ExternalMixerVolume = EngineDeck2.ExternalMixerVolume()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.ExternalScratchWheelTouch] method instead.
	EngineDeck2ExternalScratchWheelTouch = EngineDeck2.ExternalScratchWheelTouch()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.PadsView] method instead.
	EngineDeck2PadsView = EngineDeck2.PadsView()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.Play] method instead.
	EngineDeck2Play = EngineDeck2.Play()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.PlayState] method instead.
	EngineDeck2PlayState = EngineDeck2.PlayState()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.PlayStatePath] method instead.
	EngineDeck2PlayStatePath = EngineDeck2.PlayStatePath()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.Speed] method instead.
	EngineDeck2Speed = EngineDeck2.Speed()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.SpeedNeutral] method instead.
	EngineDeck2SpeedNeutral = EngineDeck2.SpeedNeutral()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.SpeedOffsetDown] method instead.
	EngineDeck2SpeedOffsetDown = EngineDeck2.SpeedOffsetDown()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.SpeedOffsetUp] method instead.
	EngineDeck2SpeedOffsetUp = EngineDeck2.SpeedOffsetUp()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.SpeedRange] method instead.
	EngineDeck2SpeedRange = EngineDeck2.SpeedRange()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.SpeedState] method instead.
	EngineDeck2SpeedState = EngineDeck2.SpeedState()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.SyncMode] method instead.
	EngineDeck2SyncMode = EngineDeck2.SyncMode()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackArtistName] method instead.
	EngineDeck2TrackArtistName = EngineDeck2.TrackArtistName()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackBleep] method instead.
	EngineDeck2TrackBleep = EngineDeck2.TrackBleep()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackCuePosition] method instead.
	EngineDeck2TrackCuePosition = EngineDeck2.TrackCuePosition()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackCurrentBPM] method instead.
	EngineDeck2TrackCurrentBPM = EngineDeck2.TrackCurrentBPM()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackCurrentKeyIndex] method instead.
	EngineDeck2TrackCurrentKeyIndex = EngineDeck2.TrackCurrentKeyIndex()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackCurrentLoopInPosition] method instead.
	EngineDeck2TrackCurrentLoopInPosition = EngineDeck2.TrackCurrentLoopInPosition()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackCurrentLoopOutPosition] method instead.
	EngineDeck2TrackCurrentLoopOutPosition = EngineDeck2.TrackCurrentLoopOutPosition()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackCurrentLoopSizeInBeats] method instead.
	EngineDeck2TrackCurrentLoopSizeInBeats = EngineDeck2.TrackCurrentLoopSizeInBeats()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackKeyLock] method instead.
	EngineDeck2TrackKeyLock = EngineDeck2.TrackKeyLock()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopEnableState] method instead.
	EngineDeck2TrackLoopEnableState = EngineDeck2.TrackLoopEnableState()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopQuickLoop1] method instead.
	EngineDeck2TrackLoopQuickLoop1 = EngineDeck2.TrackLoopQuickLoop1()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopQuickLoop2] method instead.
	EngineDeck2TrackLoopQuickLoop2 = EngineDeck2.TrackLoopQuickLoop2()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopQuickLoop3] method instead.
	EngineDeck2TrackLoopQuickLoop3 = EngineDeck2.TrackLoopQuickLoop3()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopQuickLoop4] method instead.
	EngineDeck2TrackLoopQuickLoop4 = EngineDeck2.TrackLoopQuickLoop4()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopQuickLoop5] method instead.
	EngineDeck2TrackLoopQuickLoop5 = EngineDeck2.TrackLoopQuickLoop5()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopQuickLoop6] method instead.
	EngineDeck2TrackLoopQuickLoop6 = EngineDeck2.TrackLoopQuickLoop6()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopQuickLoop7] method instead.
	EngineDeck2TrackLoopQuickLoop7 = EngineDeck2.TrackLoopQuickLoop7()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackLoopQuickLoop8] method instead.
	EngineDeck2TrackLoopQuickLoop8 = EngineDeck2.TrackLoopQuickLoop8()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackPlayPauseLEDState] method instead.
	EngineDeck2TrackPlayPauseLEDState = EngineDeck2.TrackPlayPauseLEDState()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackSampleRate] method instead.
	EngineDeck2TrackSampleRate = EngineDeck2.TrackSampleRate()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackSongAnalyzed] method instead.
	EngineDeck2TrackSongAnalyzed = EngineDeck2.TrackSongAnalyzed()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackSongLoaded] method instead.
	EngineDeck2TrackSongLoaded = EngineDeck2.TrackSongLoaded()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackSongName] method instead.
	EngineDeck2TrackSongName = EngineDeck2.TrackSongName()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackSoundSwitchGuid] method instead.
	EngineDeck2TrackSoundSwitchGUID = EngineDeck2.TrackSoundSwitchGuid()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackTrackBytes] method instead.
	EngineDeck2TrackTrackBytes = EngineDeck2.TrackTrackBytes()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackTrackData] method instead.
	EngineDeck2TrackTrackData = EngineDeck2.TrackTrackData()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackTrackLength] method instead.
	EngineDeck2TrackTrackLength = EngineDeck2.TrackTrackLength()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackTrackName] method instead.
	EngineDeck2TrackTrackName = EngineDeck2.TrackTrackName()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackTrackNetworkPath] method instead.
	EngineDeck2TrackTrackNetworkPath = EngineDeck2.TrackTrackNetworkPath()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackTrackUri] method instead.
	EngineDeck2TrackTrackURI = EngineDeck2.TrackTrackUri()
	// Deprecated: Use [EngineDeck2] and its [DeckValueNames.TrackTrackWasPlayed] method instead.
	EngineDeck2TrackTrackWasPlayed = EngineDeck2.TrackTrackWasPlayed()
)

// Deck 3 legacy variables
//
// Deprecated: Use values from [EngineDeck3] instead.
var (
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.CurrentBPM] method instead.
	EngineDeck3CurrentBPM = EngineDeck3.CurrentBPM()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.DeckIsMaster] method instead.
	EngineDeck3DeckIsMaster = EngineDeck3.DeckIsMaster()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.ExternalMixerVolume] method instead.
	EngineDeck3ExternalMixerVolume = EngineDeck3.ExternalMixerVolume()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.ExternalScratchWheelTouch] method instead.
	EngineDeck3ExternalScratchWheelTouch = EngineDeck3.ExternalScratchWheelTouch()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.PadsView] method instead.
	EngineDeck3PadsView = EngineDeck3.PadsView()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.Play] method instead.
	EngineDeck3Play = EngineDeck3.Play()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.PlayState] method instead.
	EngineDeck3PlayState = EngineDeck3.PlayState()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.PlayStatePath] method instead.
	EngineDeck3PlayStatePath = EngineDeck3.PlayStatePath()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.Speed] method instead.
	EngineDeck3Speed = EngineDeck3.Speed()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.SpeedNeutral] method instead.
	EngineDeck3SpeedNeutral = EngineDeck3.SpeedNeutral()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.SpeedOffsetDown] method instead.
	EngineDeck3SpeedOffsetDown = EngineDeck3.SpeedOffsetDown()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.SpeedOffsetUp] method instead.
	EngineDeck3SpeedOffsetUp = EngineDeck3.SpeedOffsetUp()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.SpeedRange] method instead.
	EngineDeck3SpeedRange = EngineDeck3.SpeedRange()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.SpeedState] method instead.
	EngineDeck3SpeedState = EngineDeck3.SpeedState()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.SyncMode] method instead.
	EngineDeck3SyncMode = EngineDeck3.SyncMode()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackArtistName] method instead.
	EngineDeck3TrackArtistName = EngineDeck3.TrackArtistName()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackBleep] method instead.
	EngineDeck3TrackBleep = EngineDeck3.TrackBleep()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackCuePosition] method instead.
	EngineDeck3TrackCuePosition = EngineDeck3.TrackCuePosition()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackCurrentBPM] method instead.
	EngineDeck3TrackCurrentBPM = EngineDeck3.TrackCurrentBPM()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackCurrentKeyIndex] method instead.
	EngineDeck3TrackCurrentKeyIndex = EngineDeck3.TrackCurrentKeyIndex()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackCurrentLoopInPosition] method instead.
	EngineDeck3TrackCurrentLoopInPosition = EngineDeck3.TrackCurrentLoopInPosition()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackCurrentLoopOutPosition] method instead.
	EngineDeck3TrackCurrentLoopOutPosition = EngineDeck3.TrackCurrentLoopOutPosition()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackCurrentLoopSizeInBeats] method instead.
	EngineDeck3TrackCurrentLoopSizeInBeats = EngineDeck3.TrackCurrentLoopSizeInBeats()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackKeyLock] method instead.
	EngineDeck3TrackKeyLock = EngineDeck3.TrackKeyLock()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopEnableState] method instead.
	EngineDeck3TrackLoopEnableState = EngineDeck3.TrackLoopEnableState()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopQuickLoop1] method instead.
	EngineDeck3TrackLoopQuickLoop1 = EngineDeck3.TrackLoopQuickLoop1()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopQuickLoop2] method instead.
	EngineDeck3TrackLoopQuickLoop2 = EngineDeck3.TrackLoopQuickLoop2()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopQuickLoop3] method instead.
	EngineDeck3TrackLoopQuickLoop3 = EngineDeck3.TrackLoopQuickLoop3()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopQuickLoop4] method instead.
	EngineDeck3TrackLoopQuickLoop4 = EngineDeck3.TrackLoopQuickLoop4()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopQuickLoop5] method instead.
	EngineDeck3TrackLoopQuickLoop5 = EngineDeck3.TrackLoopQuickLoop5()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopQuickLoop6] method instead.
	EngineDeck3TrackLoopQuickLoop6 = EngineDeck3.TrackLoopQuickLoop6()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopQuickLoop7] method instead.
	EngineDeck3TrackLoopQuickLoop7 = EngineDeck3.TrackLoopQuickLoop7()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackLoopQuickLoop8] method instead.
	EngineDeck3TrackLoopQuickLoop8 = EngineDeck3.TrackLoopQuickLoop8()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackPlayPauseLEDState] method instead.
	EngineDeck3TrackPlayPauseLEDState = EngineDeck3.TrackPlayPauseLEDState()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackSampleRate] method instead.
	EngineDeck3TrackSampleRate = EngineDeck3.TrackSampleRate()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackSongAnalyzed] method instead.
	EngineDeck3TrackSongAnalyzed = EngineDeck3.TrackSongAnalyzed()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackSongLoaded] method instead.
	EngineDeck3TrackSongLoaded = EngineDeck3.TrackSongLoaded()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackSongName] method instead.
	EngineDeck3TrackSongName = EngineDeck3.TrackSongName()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackSoundSwitchGuid] method instead.
	EngineDeck3TrackSoundSwitchGUID = EngineDeck3.TrackSoundSwitchGuid()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackTrackBytes] method instead.
	EngineDeck3TrackTrackBytes = EngineDeck3.TrackTrackBytes()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackTrackData] method instead.
	EngineDeck3TrackTrackData = EngineDeck3.TrackTrackData()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackTrackLength] method instead.
	EngineDeck3TrackTrackLength = EngineDeck3.TrackTrackLength()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackTrackName] method instead.
	EngineDeck3TrackTrackName = EngineDeck3.TrackTrackName()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackTrackNetworkPath] method instead.
	EngineDeck3TrackTrackNetworkPath = EngineDeck3.TrackTrackNetworkPath()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackTrackUri] method instead.
	EngineDeck3TrackTrackURI = EngineDeck3.TrackTrackUri()
	// Deprecated: Use [EngineDeck3] and its [DeckValueNames.TrackTrackWasPlayed] method instead.
	EngineDeck3TrackTrackWasPlayed = EngineDeck3.TrackTrackWasPlayed()
)

// Deck 4 legacy variables
//
// Deprecated: Use values from [EngineDeck4] instead.
var (
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.CurrentBPM] method instead.
	EngineDeck4CurrentBPM = EngineDeck4.CurrentBPM()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.DeckIsMaster] method instead.
	EngineDeck4DeckIsMaster = EngineDeck4.DeckIsMaster()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.ExternalMixerVolume] method instead.
	EngineDeck4ExternalMixerVolume = EngineDeck4.ExternalMixerVolume()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.ExternalScratchWheelTouch] method instead.
	EngineDeck4ExternalScratchWheelTouch = EngineDeck4.ExternalScratchWheelTouch()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.PadsView] method instead.
	EngineDeck4PadsView = EngineDeck4.PadsView()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.Play] method instead.
	EngineDeck4Play = EngineDeck4.Play()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.PlayState] method instead.
	EngineDeck4PlayState = EngineDeck4.PlayState()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.PlayStatePath] method instead.
	EngineDeck4PlayStatePath = EngineDeck4.PlayStatePath()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.Speed] method instead.
	EngineDeck4Speed = EngineDeck4.Speed()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.SpeedNeutral] method instead.
	EngineDeck4SpeedNeutral = EngineDeck4.SpeedNeutral()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.SpeedOffsetDown] method instead.
	EngineDeck4SpeedOffsetDown = EngineDeck4.SpeedOffsetDown()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.SpeedOffsetUp] method instead.
	EngineDeck4SpeedOffsetUp = EngineDeck4.SpeedOffsetUp()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.SpeedRange] method instead.
	EngineDeck4SpeedRange = EngineDeck4.SpeedRange()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.SpeedState] method instead.
	EngineDeck4SpeedState = EngineDeck4.SpeedState()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.SyncMode] method instead.
	EngineDeck4SyncMode = EngineDeck4.SyncMode()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackArtistName] method instead.
	EngineDeck4TrackArtistName = EngineDeck4.TrackArtistName()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackBleep] method instead.
	EngineDeck4TrackBleep = EngineDeck4.TrackBleep()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackCuePosition] method instead.
	EngineDeck4TrackCuePosition = EngineDeck4.TrackCuePosition()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackCurrentBPM] method instead.
	EngineDeck4TrackCurrentBPM = EngineDeck4.TrackCurrentBPM()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackCurrentKeyIndex] method instead.
	EngineDeck4TrackCurrentKeyIndex = EngineDeck4.TrackCurrentKeyIndex()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackCurrentLoopInPosition] method instead.
	EngineDeck4TrackCurrentLoopInPosition = EngineDeck4.TrackCurrentLoopInPosition()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackCurrentLoopOutPosition] method instead.
	EngineDeck4TrackCurrentLoopOutPosition = EngineDeck4.TrackCurrentLoopOutPosition()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackCurrentLoopSizeInBeats] method instead.
	EngineDeck4TrackCurrentLoopSizeInBeats = EngineDeck4.TrackCurrentLoopSizeInBeats()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackKeyLock] method instead.
	EngineDeck4TrackKeyLock = EngineDeck4.TrackKeyLock()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopEnableState] method instead.
	EngineDeck4TrackLoopEnableState = EngineDeck4.TrackLoopEnableState()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopQuickLoop1] method instead.
	EngineDeck4TrackLoopQuickLoop1 = EngineDeck4.TrackLoopQuickLoop1()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopQuickLoop2] method instead.
	EngineDeck4TrackLoopQuickLoop2 = EngineDeck4.TrackLoopQuickLoop2()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopQuickLoop3] method instead.
	EngineDeck4TrackLoopQuickLoop3 = EngineDeck4.TrackLoopQuickLoop3()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopQuickLoop4] method instead.
	EngineDeck4TrackLoopQuickLoop4 = EngineDeck4.TrackLoopQuickLoop4()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopQuickLoop5] method instead.
	EngineDeck4TrackLoopQuickLoop5 = EngineDeck4.TrackLoopQuickLoop5()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopQuickLoop6] method instead.
	EngineDeck4TrackLoopQuickLoop6 = EngineDeck4.TrackLoopQuickLoop6()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopQuickLoop7] method instead.
	EngineDeck4TrackLoopQuickLoop7 = EngineDeck4.TrackLoopQuickLoop7()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackLoopQuickLoop8] method instead.
	EngineDeck4TrackLoopQuickLoop8 = EngineDeck4.TrackLoopQuickLoop8()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackPlayPauseLEDState] method instead.
	EngineDeck4TrackPlayPauseLEDState = EngineDeck4.TrackPlayPauseLEDState()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackSampleRate] method instead.
	EngineDeck4TrackSampleRate = EngineDeck4.TrackSampleRate()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackSongAnalyzed] method instead.
	EngineDeck4TrackSongAnalyzed = EngineDeck4.TrackSongAnalyzed()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackSongLoaded] method instead.
	EngineDeck4TrackSongLoaded = EngineDeck4.TrackSongLoaded()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackSongName] method instead.
	EngineDeck4TrackSongName = EngineDeck4.TrackSongName()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackSoundSwitchGuid] method instead.
	EngineDeck4TrackSoundSwitchGUID = EngineDeck4.TrackSoundSwitchGuid()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackTrackBytes] method instead.
	EngineDeck4TrackTrackBytes = EngineDeck4.TrackTrackBytes()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackTrackData] method instead.
	EngineDeck4TrackTrackData = EngineDeck4.TrackTrackData()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackTrackLength] method instead.
	EngineDeck4TrackTrackLength = EngineDeck4.TrackTrackLength()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackTrackName] method instead.
	EngineDeck4TrackTrackName = EngineDeck4.TrackTrackName()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackTrackNetworkPath] method instead.
	EngineDeck4TrackTrackNetworkPath = EngineDeck4.TrackTrackNetworkPath()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackTrackUri] method instead.
	EngineDeck4TrackTrackURI = EngineDeck4.TrackTrackUri()
	// Deprecated: Use [EngineDeck4] and its [DeckValueNames.TrackTrackWasPlayed] method instead.
	EngineDeck4TrackTrackWasPlayed = EngineDeck4.TrackTrackWasPlayed()
)

var (
	EngineDeck1 = DeckValueNames{DeckIndex: 1}
	EngineDeck2 = DeckValueNames{DeckIndex: 2}
	EngineDeck3 = DeckValueNames{DeckIndex: 3}
	EngineDeck4 = DeckValueNames{DeckIndex: 4}
)

type DeckValueNames struct {
	DeckIndex int
}

func (n *DeckValueNames) PadsView() string {
	return fmt.Sprintf("/Engine/Deck%d/Pads/View", n.DeckIndex)
}

func (n *DeckValueNames) TrackArtistName() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/ArtistName", n.DeckIndex)
}

func (n *DeckValueNames) TrackBleep() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Bleep", n.DeckIndex)
}

func (n *DeckValueNames) TrackCuePosition() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/CuePosition", n.DeckIndex)
}

func (n *DeckValueNames) TrackCurrentBPM() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/CurrentBPM", n.DeckIndex)
}

func (n *DeckValueNames) TrackCurrentKeyIndex() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/CurrentKeyIndex", n.DeckIndex)
}

func (n *DeckValueNames) TrackCurrentLoopInPosition() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/CurrentLoopInPosition", n.DeckIndex)
}

func (n *DeckValueNames) TrackCurrentLoopOutPosition() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/CurrentLoopOutPosition", n.DeckIndex)
}

func (n *DeckValueNames) TrackCurrentLoopSizeInBeats() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/CurrentLoopSizeInBeats", n.DeckIndex)
}

func (n *DeckValueNames) TrackKeyLock() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/KeyLock", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopEnableState() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/LoopEnableState", n.DeckIndex)
}

func (n *DeckValueNames) TrackPlayPauseLEDState() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/PlayPauseLEDState", n.DeckIndex)
}

func (n *DeckValueNames) TrackSampleRate() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/SampleRate", n.DeckIndex)
}

func (n *DeckValueNames) TrackSongAnalyzed() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/SongAnalyzed", n.DeckIndex)
}

func (n *DeckValueNames) TrackSongLoaded() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/SongLoaded", n.DeckIndex)
}

func (n *DeckValueNames) TrackSongName() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/SongName", n.DeckIndex)
}

func (n *DeckValueNames) TrackSoundSwitchGuid() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/SoundSwitchGuid", n.DeckIndex)
}

func (n *DeckValueNames) TrackTrackBytes() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/TrackBytes", n.DeckIndex)
}

func (n *DeckValueNames) TrackTrackData() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/TrackData", n.DeckIndex)
}

func (n *DeckValueNames) TrackTrackLength() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/TrackLength", n.DeckIndex)
}

func (n *DeckValueNames) TrackTrackName() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/TrackName", n.DeckIndex)
}

func (n *DeckValueNames) TrackTrackNetworkPath() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/TrackNetworkPath", n.DeckIndex)
}

func (n *DeckValueNames) TrackTrackUri() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/TrackUri", n.DeckIndex)
}

func (n *DeckValueNames) TrackTrackWasPlayed() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/TrackWasPlayed", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopQuickLoop1() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Loop/QuickLoop1", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopQuickLoop2() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Loop/QuickLoop2", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopQuickLoop3() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Loop/QuickLoop3", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopQuickLoop4() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Loop/QuickLoop4", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopQuickLoop5() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Loop/QuickLoop5", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopQuickLoop6() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Loop/QuickLoop6", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopQuickLoop7() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Loop/QuickLoop7", n.DeckIndex)
}

func (n *DeckValueNames) TrackLoopQuickLoop8() string {
	return fmt.Sprintf("/Engine/Deck%d/Track/Loop/QuickLoop8", n.DeckIndex)
}

func (n *DeckValueNames) CurrentBPM() string {
	return fmt.Sprintf("/Engine/Deck%d/CurrentBPM", n.DeckIndex)
}

func (n *DeckValueNames) DeckIsMaster() string {
	return fmt.Sprintf("/Engine/Deck%d/DeckIsMaster", n.DeckIndex)
}

func (n *DeckValueNames) ExternalMixerVolume() string {
	return fmt.Sprintf("/Engine/Deck%d/ExternalMixerVolume", n.DeckIndex)
}

func (n *DeckValueNames) ExternalScratchWheelTouch() string {
	return fmt.Sprintf("/Engine/Deck%d/ExternalScratchWheelTouch", n.DeckIndex)
}

func (n *DeckValueNames) Play() string {
	return fmt.Sprintf("/Engine/Deck%d/Play", n.DeckIndex)
}

func (n *DeckValueNames) PlayState() string {
	return fmt.Sprintf("/Engine/Deck%d/PlayState", n.DeckIndex)
}

func (n *DeckValueNames) PlayStatePath() string {
	return fmt.Sprintf("/Engine/Deck%d/PlayStatePath", n.DeckIndex)
}

func (n *DeckValueNames) Speed() string {
	return fmt.Sprintf("/Engine/Deck%d/Speed", n.DeckIndex)
}

func (n *DeckValueNames) SpeedNeutral() string {
	return fmt.Sprintf("/Engine/Deck%d/SpeedNeutral", n.DeckIndex)
}

func (n *DeckValueNames) SpeedOffsetDown() string {
	return fmt.Sprintf("/Engine/Deck%d/SpeedOffsetDown", n.DeckIndex)
}

func (n *DeckValueNames) SpeedOffsetUp() string {
	return fmt.Sprintf("/Engine/Deck%d/SpeedOffsetUp", n.DeckIndex)
}

func (n *DeckValueNames) SpeedRange() string {
	return fmt.Sprintf("/Engine/Deck%d/SpeedRange", n.DeckIndex)
}

func (n *DeckValueNames) SpeedState() string {
	return fmt.Sprintf("/Engine/Deck%d/SpeedState", n.DeckIndex)
}

func (n *DeckValueNames) SyncMode() string {
	return fmt.Sprintf("/Engine/Deck%d/SyncMode", n.DeckIndex)
}
