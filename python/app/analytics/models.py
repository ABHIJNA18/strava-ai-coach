# This file defines the structured analytics objects used by the coaching pipeline.

from dataclasses import dataclass


@dataclass
class RunSummary:
    run_count: int
    total_distance_meters: float
    average_run_distance_meters: float
    total_moving_time_seconds: int
    average_run_duration_seconds: float
    average_pace_seconds_per_km: float
    average_heartrate: float
    average_cadence: float
    total_elevation_gain_meters: float

    fastest_run_name: str
    fastest_run_distance_meters: float
    fastest_run_pace_seconds_per_km: float
    fastest_run_date: str

    longest_run_name: str
    longest_run_distance_meters: float
    longest_run_pace_seconds_per_km: float
    longest_run_date: str


@dataclass
class RunningAnalytics:
    summary: RunSummary


@dataclass
class WeeklyRunAnalytics:
    week_start: str
    week_end: str

    run_count: int
    total_distance_meters: float
    average_run_distance_meters: float
    average_pace_seconds_per_km: float
    average_heartrate: float
    average_cadence: float
    total_elevation_gain_meters: float
    longest_run_distance_meters: float
    fastest_run_pace_seconds_per_km: float


@dataclass
class WeeklyRunningAnalytics:
    weeks: list[WeeklyRunAnalytics]


@dataclass
class CoachingAnalytics:
    summary: RunSummary
    weekly: WeeklyRunningAnalytics