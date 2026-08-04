# This file calculates overall running analytics from activity data.

from python.app.analytics.models import RunSummary


def _runs(activities):
    return [activity for activity in activities if activity.sport_type == "Run"]


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

    total_distance = sum(activity.distance for activity in runs)
    total_moving_time = sum(activity.moving_time for activity in runs)
    total_elevation_gain = sum(activity.total_elevation_gain for activity in runs)

    average_run_distance = 0
    if runs:
        average_run_distance = total_distance / len(runs)

    average_pace = 0
    if total_distance > 0:
        average_pace = total_moving_time / (total_distance / 1000)

    fastest_run_pace = 0
    runs_with_distance = [activity for activity in runs if activity.distance > 0]
    if runs_with_distance:
        fastest_run_pace = min(_pace_seconds_per_km(activity) for activity in runs_with_distance)

    longest_run_distance = 0
    if runs:
        longest_run_distance = max(activity.distance for activity in runs)

    return RunSummary(
        run_count=len(runs),
        total_distance_meters=total_distance,
        average_run_distance_meters=average_run_distance,
        total_moving_time_seconds=total_moving_time,
        average_pace_seconds_per_km=average_pace,
        average_heartrate=_average(activity.average_heartrate for activity in runs),
        average_cadence=_average(activity.average_cadence for activity in runs),
        total_elevation_gain_meters=total_elevation_gain,
        fastest_run_pace_seconds_per_km=fastest_run_pace,
        longest_run_distance_meters=longest_run_distance,
    )
