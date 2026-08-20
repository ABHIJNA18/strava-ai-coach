# This file tests weekly running analytics and edge cases.

from types import SimpleNamespace

from python.app.analytics.weekly import (
    calculate_weekly_running_analytics,
)


def activity(
    date,
    distance,
    moving_time,
    heartrate=150,
):
    return SimpleNamespace(
        sport_type="Run",
        start_date=date,
        distance=distance,
        moving_time=moving_time,
        average_heartrate=heartrate,
        average_cadence=80,
        total_elevation_gain=10,
    )


def test_groups_activities_into_monday_sunday_weeks():
    activities = [
        activity(
            "2026-08-03T08:00:00Z",
            5000,
            1500,
        ),
        activity(
            "2026-08-05T08:00:00Z",
            7000,
            2100,
        ),
    ]

    result = calculate_weekly_running_analytics(
        activities
    )

    assert len(result.weeks) == 1
    assert result.weeks[0].week_start == "2026-08-03"
    assert result.weeks[0].run_count == 2
    assert result.weeks[0].total_distance_meters == 12000
    assert result.weeks[0].longest_run_distance_meters == 7000


def test_preserves_empty_calendar_weeks():
    activities = [
        activity(
            "2026-08-03T08:00:00Z",
            5000,
            1500,
        ),
        activity(
            "2026-08-17T08:00:00Z",
            6000,
            1800,
        ),
    ]

    result = calculate_weekly_running_analytics(
        activities
    )

    assert len(result.weeks) == 3
    assert result.weeks[1].run_count == 0
    assert result.weeks[1].total_distance_meters == 0
    assert result.weeks[1].average_pace_seconds_per_km == 0


def test_ignores_invalid_dates():
    activities = [
        activity(
            "not-a-date",
            5000,
            1500,
        )
    ]

    result = calculate_weekly_running_analytics(
        activities
    )

    assert result.weeks == []


def test_zero_distance_does_not_crash():
    activities = [
        activity(
            "2026-08-03T08:00:00Z",
            0,
            1000,
        )
    ]

    result = calculate_weekly_running_analytics(
        activities
    )

    assert result.weeks[0].run_count == 1
    assert result.weeks[0].average_pace_seconds_per_km == 0