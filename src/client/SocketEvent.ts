import { SignalData } from 'simple-peer'

export interface Ready {
  room: string
  peerId: string
  nickname: string
}

export interface TrackMetadata {
  mid: string
  kind: string
  peerId: string
  streamId: string
}

export interface MetadataPayload {
  peerId: string
  metadata: TrackMetadata[]
}

export enum TrackEventType {
  Add = 1,
  Remove = 2,
  Sub = 3,
  Unsub = 4,
}

  // TrackId maps to identifiers.TrackID.
export interface TrackId {
  id: string
  streamId: string
}

// TrackKind maps to transport.TrackKind.
export type TrackKind = 'audio' | 'video'

export interface SocketEvent {
  users: {
    initiator: string
    // peers to connect to
    peerIds: string[]
    // mapping of peerId / nickname
    nicknames: Record<string, string>
  }
  metadata: MetadataPayload
  hangUp: {
    peerId: string
  }
  pubTrack: {
    trackId: TrackId
    pubClientId: string
    peerId: string
    kind: TrackKind
    type: TrackEventType.Add | TrackEventType.Remove
  }
  subTrack: {
    trackId: TrackId
    pubClientId: string
    type: TrackEventType.Sub | TrackEventType.Unsub
  }
  signal: {
    peerId: string
    // eslint-disable-next-line
    signal: SignalData
  }
  connect: undefined
  disconnect: undefined
  ready: Ready
}
