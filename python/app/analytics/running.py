# This file calculates overall running analytics from activity data.

from python.app.analytics.models import RunSummary


def _runs(activities):
    return [
        activity
        for activity in activities
        if activity.sport_type == "Run"
    ]


def _pace_seconds_per_km(activity):
    if activity.distance <= 0:
        return 0

    return activity.moving_time / (activity.distance / 1000)


def _average(values):
    values = [value for value in values if value > 0]

    if not values:
        return 0

    return sum(values) / len(values)


def calculate_run_summary(activities):
    runs = _runs(activities)

    total_distance = sum(
        activity.distance
        for activity in runs
    )

    total_moving_time = sum(
        activity.moving_time
        for activity in runs
    )

    total_elevation_gain = sum(
        activity.total_elevation_gain
        for activity in runs
    )

    average_run_distance = 0
    average_run_duration = 0
    average_pace = 0

    if runs:
        average_run_distance = total_distance / len(runs)
        average_run_duration = total_moving_time / len(runs)

    if total_distance > 0:
        average_pace = (
            total_moving_time
            / (total_distance / 1000)
        )

    runs_with_distance = [
        activity
        for activity in runs
        if activity.distance > 0
    ]

    fastest_run = None

    if runs_with_distance:
        fastest_run = min(
            runs_with_distance,
            key=lambda activity: _pace_seconds_per_km(activity),
        )

    longest_run = None

    if runs:
        longest_run = max(
            runs,
            key=lambda activity: activity.distance,
        )

    return RunSummary(
        run_count=len(runs),

        total_distance_meters=total_distance,

        average_run_distance_meters=average_run_distance,

        total_moving_time_seconds=total_moving_time,

        average_run_duration_seconds=average_run_duration,

        average_pace_seconds_per_km=average_pace,

        average_heartrate=_average(
            activity.average_heartrate
            for activity in runs
        ),

        # Garmin/Strava reports running cadence per foot.
        # Convert to total steps per minute.
        average_cadence=_average(
            activity.average_cadence
            for activity in runs
        ) * 2,

        total_elevation_gain_meters=total_elevation_gain,

        fastest_run_name=(
            fastest_run.name
            if fastest_run
            else ""
        ),

        fastest_run_distance_meters=(
            fastest_run.distance
            if fastest_run
            else 0
        ),

        fastest_run_pace_seconds_per_km=(
            _pace_seconds_per_km(fastest_run)
            if fastest_run
            else 0
        ),

        fastest_run_date=(
            fastest_run.start_date
            if fastest_run
            else ""
        ),

        longest_run_name=(
            longest_run.name
            if longest_run
            else ""
        ),

        longest_run_distance_meters=(
            longest_run.distance
            if longest_run
            else 0
        ),

        longest_run_pace_seconds_per_km=(
            _pace_seconds_per_km(longest_run)
            if longest_run
            else 0
        ),

        longest_run_date=(
            longest_run.start_date
            if longest_run
            else ""
        ),
    )