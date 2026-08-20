# This file calculates weekly running analytics from supplied activities.

from collections import defaultdict
from datetime import date, datetime, timedelta

from python.app.analytics.models import (
    WeeklyRunAnalytics,
    WeeklyRunningAnalytics,
)


def _parse_date(value):
    if not value:
        return None

    try:
        return datetime.fromisoformat(
            value.replace("Z", "+00:00")
        ).date()
    except ValueError:
        return None


def _week_start(activity_date):
    return activity_date - timedelta(
        days=activity_date.weekday()
    )


def _pace_seconds_per_km(activity):
    if activity.distance <= 0:
        return 0

    if activity.moving_time <= 0:
        return 0

    return activity.moving_time / (
        activity.distance / 1000
    )


def _average(values):
    valid_values = [
        value
        for value in values
        if value > 0
    ]

    if not valid_values:
        return 0

    return sum(valid_values) / len(valid_values)


def _empty_week(week_start):
    return WeeklyRunAnalytics(
        week_start=week_start.isoformat(),
        week_end=(
            week_start + timedelta(days=6)
        ).isoformat(),
        run_count=0,
        total_distance_meters=0,
        average_run_distance_meters=0,
        average_pace_seconds_per_km=0,
        average_heartrate=0,
        average_cadence=0,
        total_elevation_gain_meters=0,
        longest_run_distance_meters=0,
        fastest_run_pace_seconds_per_km=0,
    )


def calculate_weekly_running_analytics(
    activities,
    analysis_start=None,
    analysis_end=None,
):
    """
    Calculate Monday-Sunday analytics for the supplied period.

    The default period is the last 60 calendar days.
    """

    if analysis_end is None:
        analysis_end = date.today()

    if analysis_start is None:
        analysis_start = analysis_end - timedelta(days=59)

    dated_runs = []

    for activity in activities:
        if activity.sport_type != "Run":
            continue

        activity_date = _parse_date(
            activity.start_date
        )

        if activity_date is None:
            continue

        if activity_date < analysis_start:
            continue

        if activity_date > analysis_end:
            continue

        dated_runs.append(
            (activity_date, activity)
        )

    first_week = _week_start(analysis_start)
    last_week = _week_start(analysis_end)

    grouped_runs = defaultdict(list)

    for activity_date, activity in dated_runs:
        grouped_runs[
            _week_start(activity_date)
        ].append(activity)

    weekly_results = []
    current_week = first_week

    while current_week <= last_week:
        runs = grouped_runs.get(
            current_week,
            [],
        )

        if not runs:
            weekly_results.append(
                _empty_week(current_week)
            )

            current_week += timedelta(days=7)
            continue

        total_distance = sum(
            activity.distance
            for activity in runs
        )

        valid_distance_runs = [
            activity
            for activity in runs
            if activity.distance > 0
        ]

        valid_pace_runs = [
            activity
            for activity in valid_distance_runs
            if activity.moving_time > 0
        ]

        average_pace = _average(
            _pace_seconds_per_km(activity)
            for activity in valid_pace_runs
        )

        fastest_pace = 0

        if valid_pace_runs:
            fastest_pace = min(
                _pace_seconds_per_km(activity)
                for activity in valid_pace_runs
            )

        longest_run = max(
            (
                activity.distance
                for activity in runs
            ),
            default=0,
        )

        weekly_results.append(
            WeeklyRunAnalytics(
                week_start=current_week.isoformat(),
                week_end=(
                    current_week + timedelta(days=6)
                ).isoformat(),
                run_count=len(runs),
                total_distance_meters=total_distance,
                average_run_distance_meters=(
                    total_distance / len(runs)
                    if runs
                    else 0
                ),
                average_pace_seconds_per_km=average_pace,
                average_heartrate=_average(
                    activity.average_heartrate
                    for activity in runs
                ),
                average_cadence=_average(
                    activity.average_cadence
                    for activity in runs
                ) * 2,
                total_elevation_gain_meters=sum(
                    activity.total_elevation_gain
                    for activity in runs
                ),
                longest_run_distance_meters=longest_run,
                fastest_run_pace_seconds_per_km=fastest_pace,
            )
        )

        current_week += timedelta(days=7)

    print('Weekly analytics calculated:', weekly_results)
    return WeeklyRunningAnalytics(
        weeks=weekly_results
    )