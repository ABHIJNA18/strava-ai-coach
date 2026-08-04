# This file defines the structured analytics objects used by the coaching pipeline.

from dataclasses import dataclass


@dataclass
class RunSummary:
    run_count: int
    total_distance_meters: float
    average_run_distance_meters: float
    total_moving_time_seconds: int
    average_pace_seconds_per_km: float
    average_heartrate: float
    average_cadence: float
    total_elevation_gain_meters: float
    fastest_run_pace_seconds_per_km: float
    longest_run_distance_meters: float
    average_run_duration_seconds: float
    longest_run_date: str


@dataclass
class RunningAnalytics:
    summary: RunSummary
