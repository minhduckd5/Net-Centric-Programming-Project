# Text-Based Clash Royale (TCR)

## Protocol Data Unit (PDU) Specification

---

## Table of Contents

1. [Overview](#overview)
2. [PDU Structure](#pdu-structure)
3. [Message Types](#message-types)
4. [PDU Definitions](#pdu-definitions)

   * 4.1 [Authentication](#authentication-pdus)
   * 4.2 [Game Management](#game-management-pdus)
   * 4.3 [Game Actions](#game-action-pdus)
   * 4.4 [System Messages](#system-message-pdus)
5. [Error Handling](#error-handling)
6. [Sequence Diagrams](#sequence-examples)

---

## 1. Overview

This document defines the JSON-based PDUs for TCR’s client-server communication. Each PDU conveys a discrete operation or event, ensuring consistency and extensibility across authentication, matchmaking, gameplay, and system control flows.

---

## 2. PDU Structure

All PDUs adhere to the following JSON template:

```json
{
  "type": "<MESSAGE_TYPE>",
  "data": { /* payload object */ }
}
```

* **type**: Identifier for message semantics.
* **data**: Context-specific payload containing required parameters.

---

## 3. Message Types

PDUs are classified into four categories:

| Category            | Message Types                                                                            |
| ------------------- | ---------------------------------------------------------------------------------------- |
| **Authentication**  | LOGIN\_REQUEST, LOGIN\_RESPONSE, LOGOUT\_REQUEST, LOGOUT\_RESPONSE                       |
| **Game Management** | MATCHMAKING\_REQUEST, MATCHMAKING\_RESPONSE, GAME\_START, GAME\_END, GAME\_STATE\_UPDATE |
| **Game Actions**    | DEPLOY\_TROOP, TOWER\_ATTACK, TROOP\_ATTACK, MANA\_UPDATE, EXP\_UPDATE                   |
| **System Messages** | ERROR, PING, PONG, DISCONNECT                                                            |

---

## 4. PDU Definitions

### 4.1 Authentication PDUs {#authentication-pdus}

#### LOGIN\_REQUEST

```json
{
  "type": "LOGIN_REQUEST",
  "data": {
    "username": "<string>",
    "password": "<string>"
  }
}
```

#### LOGIN\_RESPONSE

```json
{
  "type": "LOGIN_RESPONSE",
  "data": {
    "success": <boolean>,
    "message": "<string>",
    "user": {
      "username": "<string>",
      "experience": <int>
    }
  }
}
```

#### LOGOUT\_REQUEST

```json
{
  "type": "LOGOUT_REQUEST",
  "data": { }
}
```

#### LOGOUT\_RESPONSE

```json
{
  "type": "LOGOUT_RESPONSE",
  "data": {
    "success": <boolean>,
    "message": "<string>"
  }
}
```

---

### 4.2 Game Management PDUs {#game-management-pdus}

#### MATCHMAKING\_REQUEST

```json
{
  "type": "MATCHMAKING_REQUEST",
  "data": {
    "username": "<string>"
  }
}
```

#### MATCHMAKING\_RESPONSE

```json
{
  "type": "MATCHMAKING_RESPONSE",
  "data": {
    "session_id": "<string>",
    "opponent": "<string>"
  }
}
```

#### GAME\_START

```json
{
  "type": "GAME_START",
  "data": {
    "session_id": "<string>",
    "players": [
      {
        "username": "<string>",
        "towers": [
          {
            "type": "<string>",
            "hp": <int>,
            "atk": <int>,
            "def": <int>
          }
        ]
      }
    ]
  }
}
```

#### GAME\_END

```json
{
  "type": "GAME_END",
  "data": {
    "winner": "<string>",
    "duration": <int>  // seconds
  }
}
```

#### GAME\_STATE\_UPDATE

```json
{
  "type": "GAME_STATE_UPDATE",
  "data": {
    "turn": <int>,
    "mana": <int>,
    "troops": [
      {
        "type": "<string>",
        "hp": <int>,
        "position": { "x": <int>, "y": <int> }
      }
    ],
    "towers": [ { "type": "<string>", "hp": <int> } ]
  }
}
```

---

### 4.3 Game Action PDUs {#game-action-pdus}

#### DEPLOY\_TROOP

```json
{
  "type": "DEPLOY_TROOP",
  "data": {
    "troop_type": "<string>",
    "position": { "x": <int>, "y": <int> }
  }
}
```

#### TOWER\_ATTACK

```json
{
  "type": "TOWER_ATTACK",
  "data": {
    "tower_id": "<string>",
    "target_id": "<string>"
  }
}
```

#### TROOP\_ATTACK

```json
{
  "type": "TROOP_ATTACK",
  "data": {
    "troop_id": "<string>",
    "target_id": "<string>"
  }
}
```

#### MANA\_UPDATE

```json
{
  "type": "MANA_UPDATE",
  "data": { "mana": <int> }
}
```

#### EXP\_UPDATE

```json
{
  "type": "EXP_UPDATE",
  "data": { "experience": <int> }
}
```

---

### 4.4 System Message PDUs {#system-message-pdus}

#### PING

```json
{
  "type": "PING",
  "data": { }
}
```

#### PONG

```json
{
  "type": "PONG",
  "data": { "timestamp": <int> }
}
```

#### DISCONNECT

```json
{
  "type": "DISCONNECT",
  "data": { "reason": "<string>" }
}
```

#### ERROR

Defined in Section 5.

---

## 5. Error Handling

All errors are reported via the **ERROR** PDU:

```json
{
  "type": "ERROR",
  "data": {
    "code": <int>,
    "message": "<string>",
    "details": { /* optional object */ }
  }
}
```

**Error Code Ranges:**

* 1000–1999: Authentication
* 2000–2999: Game State
* 3000–3999: Network
* 4000–4999: System

---

## 6. Sequence Examples {#sequence-examples}

### 6.1 Login Flow

```text
Client → Server: LOGIN_REQUEST → Server: LOGIN_RESPONSE
```

### 6.2 Matchmaking & Start

```text
Client → Server: MATCHMAKING_REQUEST → Server: MATCHMAKING_RESPONSE → Server: GAME_START → Server: GAME_STATE_UPDATE
```

### 6.3 In-Game Actions

```text
Client → Server: DEPLOY_TROOP → Server: GAME_STATE_UPDATE → Server: TROOP_ATTACK → Server: GAME_STATE_UPDATE
```

---

*End of PDU Specification*
