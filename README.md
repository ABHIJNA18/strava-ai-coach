# 🏃 Strava AI Coach

An AI-powered running coach built with Go, Python, gRPC, PostgreSQL and OpenAI.

The application connects to a user's Strava account, stores their activities, analyzes their recent running history, and generates personalized coaching summaries using GPT.

---

## Features

### ✅ Strava OAuth

- Login securely with Strava
- Store athlete information
- Store OAuth tokens
- Automatically refresh expired access tokens

---

### ✅ Activity Synchronization

- Fetch activities from Strava
- Store activities in PostgreSQL
- Support running, hiking and weight training activities

---

### ✅ AI Coaching Pipeline

Go Backend

↓

PostgreSQL

↓

gRPC

↓

Python Analytics Engine

↓

OpenAI GPT

↓

AI Coaching Summary

The backend fetches the athlete's running activities from the last 30 days and sends them to a Python coaching service over gRPC.

The Python service:

- Calculates structured running analytics
- Builds an LLM prompt
- Calls OpenAI
- Returns a personalized coaching summary

---

## Running Analytics

Current analytics include:

- Number of runs
- Total running distance
- Average run distance
- Total moving time
- Average run duration
- Average pace
- Average heart rate
- Average cadence
- Total elevation gain
- Fastest run
- Longest run

The analytics are calculated deterministically in Python before being sent to the LLM.

The LLM is responsible only for interpreting the analytics and generating coaching advice.

---

## Tech Stack

### Backend

- Go
- Python
- gRPC
- Protocol Buffers
- PostgreSQL
- OpenAI API

### APIs

- Strava API

---

## Project Structure

```
cmd/
    api/

internal/
    handlers/
    database/
    strava/
    coach/

proto/

python/
    app/
        analytics/
        prompts/
        ai/
        server.py
```

---

## Current Status

Backend is complete for the first MVP.

Implemented:

- Strava OAuth
- PostgreSQL persistence
- Activity synchronization
- Go ↔ Python communication using gRPC
- Analytics engine
- GPT-generated coaching summaries

---

## Next Milestones

### Frontend

Build a web dashboard where users can:

- Login with Strava
- View recent activities
- Generate a 30-day AI coaching summary

### AI Chat

Allow users to ask questions such as:

- How many kilometers did I run this month?
- What was my fastest run?
- How has my pace changed?

