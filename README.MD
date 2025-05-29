# Text‑Based Clash Royale (TCR)

**Network Programming (IT096IU) Term Project**
**Instructor:** Le Thanh Son

![Go ≥1.16](https://img.shields.io/badge/Go-%3E%3D1.16-blue) ![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)

## Table of Contents

1. [Project Overview](#project-overview)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Building](#building)
6. [Running the Game](#running-the-game)
7. [Directory Structure](#directory-structure)
8. [Logging](#logging)
9. [Troubleshooting](#troubleshooting)
10. [Security Considerations](#security-considerations)
11. [Maintenance](#maintenance)

---

## 1. Project Overview

Leverage Golang, TCP/UDP protocols, and JSON-based persistence to deliver a scalable, text-centric adaptation of Supercell’s flagship title.

**Objectives:**

* Provide hands‑on experience in designing, prototyping, testing, and evaluating network-based applications.
* Implement TCP and/or UDP communication at the programming level.
* Reinforce theoretical networking concepts via practical deployment.

**Key Features:**

* Two-player synchronous combat with turn- and continuous-play modes.
* Robust EXP and leveling systems, dynamic mana mechanics.
* Configurable troop/tower specs via JSON for rapid balance iteration.

---

## 2. Prerequisites

* **Go** ≥ 1.16
* **Git** (for repository cloning)
* Proficiency with command-line interfaces and JSON editing

---

## 3. Installation

### 3.1 Clone the Repository

```bash
git clone https://github.com/minhduckd5/Net-Centric-Programming-Project.git
cd tcr
```

### 3.2 Fetch Dependencies

```bash
go mod download
```

---

## 4. Configuration

### 4.1 Server Settings

* **Development:** `config/dev.json`
* **Production:** `config/prod.json`

**Required Parameters:**

* `host` & `port`
* Game balance variables
* Logging level (`DEBUG`, `INFO`, `WARN`, `ERROR`)

### 4.2 Game Specifications

Edit `specs/game_specs.json` to tailor:

* Troop and tower attributes (HP, ATK, DEF, CRIT, MANA)
* Balance parameters (EXP rewards, spawn rates)

### 4.3 User Data

Manage player profiles in `data/players.json`:

* Account credentials
* XP, levels, and statistics

---

## 5. Building

### 5.1 Server

```bash
go build -o tcr-server server/main.go

```

### 5.2 Client

```bash
go build -o tcr-client client/main.go
```

---

## 6. Running the Game

### 6.1 Start Server

```bash
cd to tcr\cmd\server
go run main.go 
```

### 6.2 Start Client

```bash
cd to tcr/client
go run client.go 
```

---

## 7. Directory Structure

```
tcr/
├── client/           # Client source code
├── server/           # Server source code
├── config/           # Environment-specific JSON
│   ├── dev.json
│   └── prod.json
├── specs/            # JSON-based game specs
│   └── game_specs.json
├── data/             # Persistent user data
│   └── players.json
└── logs/             # Runtime logs
    ├── server.log
    ├── client.log
    └── error.log
```

---

## 8. Logging

* **Log Files:**  `logs/server.log`, `logs/client.log`, `logs/error.log`
* **Levels:** `DEBUG`, `INFO`, `WARN`, `ERROR`

---

## 9. Troubleshooting

**Common Issues & Resolutions:**

* **Connection refused:** Confirm server is active on specified host/port.
* **Authentication failed:** Validate entries in `data/players.json`.
* **Game stall:** Verify JSON syntax and parameter values in config/spec files.

**Debug Mode:**

```bash
./tcr-server -config config/dev.json -debug
./tcr-client -server localhost:9000 -debug
```

---

## 10. Security Considerations

* Enforce stringent file permissions on `config/` & `data/`.
* Deploy firewalls and rate-limiting to mitigate network threats.
* Audit logs regularly for suspicious activity.

---

## 11. Maintenance

**Routine Tasks:**

* Monitor and rotate log archives.
* Backup JSON data and config snapshots.
* Iterate on `specs/game_specs.json` for balance tuning.
* Assess server resource utilization and optimize as needed.
