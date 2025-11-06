# DLP Event Monitor

A real-time Data Loss Prevention (DLP) event monitoring system built to practice my Go skills and learn about rate limiting and DLP detection. This system would be receiving events from clients in real-time and flagging suspicious activity.

## Overview

This event monitor will simulate a stream of business events with mock data and variable rates. It uses rate limiting and pattern analysis to detect anomalous behavior that may indicate data exfiltration or security violations and presents it in a TUI dashboard.

## Features

### Mock Events
- **File Operations**: Downloads, bulk downloads, file access, modification, deletion
- **Data Transfer**: USB copies, clipboard operations, cloud uploads
- **External Sharing**: External emails, cloud shares, external file transfers
- **Sensitive Resource Tracking**: Automatic detection of sensitive files (PII, financial data, HR records)

### Anomaly Detection Rules

The system flags suspicious from individual user activity based on configurable thresholds within a 5-minute sliding window. This is a basic but effective rule based approach but could be improved with a more sophisticated statistical modeling or machine learning approach. an example configuration might be:

- **Downloads**: > 10 download actions
- **USB Copies**: > 3 USB copy operations
- **External Actions**: > 5 external sharing operations
- **Sensitive Access**: > 2 accesses to sensitive resources
- **Policy Violations**:
  - Sensitive resources shared externally
  - Contractors accessing sensitive data

### Real-Time Dashboard

To present it all we've put together a TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) which shows an event stream of the events as they come in, any alerts generated, and some basic alert rate metrics.

### Architecture

- **Producer-Consumer Pattern**: Multiple event producers feeding a centralized channel to practice working with Go channels and goroutines.
- **Sliding Window Algorithm**: Time-based activity tracking with automatic cleanup as I wanted to learn more about rate limiting and sliding windows in practice.
- **Channel-Based Processing**: Efficient concurrent event filtering and alert generation to practice consuming from Go channels and goroutines to simulate working with concurrent systems.

## Installation

### Prerequisites

- Go 1.24.3 or higher

### Build

```bash
go mod download
go build -o dlp-monitor main.go
```

### Run

```bash
./dlp-monitor
```

Or directly with Go:

```bash
go run main.go
```

## Project Structure

```
.
├── main.go              # Application entry point
├── event/               # Event data structures
│   └── event.go
├── producer/            # Event producers (simulated data sources)
│   └── producer.go
├── filter/             # Anomaly detection and filtering logic
│   └── filter.go
├── tui/                # Terminal user interface
│   ├── tui.go
│   └── styles.go
└── consumer/           # Alert consumer (optional)
    └── consumer.go
```

## Technology Stack

- **Language**: Go
- **UI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
