# This file calculates running metrics from activity data.
# It does not call OpenAI, gRPC, or the database.
#Currently this file is not used as code switched to analytics folder, later once everything is working it could be deleted

from dataclasses import dataclass

@dataclass
class RunningMetrics:
    run_count: int
    total_distance_meters: float
    total_moving_time_seconds: int
    average_heartrate: float
    average_pace_seconds_per_km: float


def _runs(activities):
    return [activity for activity in activities if activity.sport_type == "Run"]


def calculate_run_count(activities):
    return len(_runs(activities))


def calculate_total_distance(activities):
    return sum(activity.distance for activity in _runs(activities))


def calculate_total_moving_time(activities):
    return sum(activity.moving_time for activity in _runs(activities))


def calculate_average_heartrate(activities):
    heart_rates = [
        activity.average_heartrate
        for activity in _runs(activities)
        if activity.average_heartrate > 0
    ]

    if not heart_rates:
        return 0

    return sum(heart_rates) / len(heart_rates)


def calculate_average_pace(activities):
    total_distance_meters = calculate_total_distance(activities)
    total_moving_time_seconds = calculate_total_moving_time(activities)

    if total_distance_meters <= 0:
        return 0

    total_distance_km = total_distance_meters / 1000
    return total_moving_time_seconds / total_distance_km

def calculate_running_metrics(activities):
    metrics= RunningMetrics(
        run_count=calculate_run_count(activities),
        total_distance_meters=calculate_total_distance(activities),
        total_moving_time_seconds=calculate_total_moving_time(activities),
        average_heartrate=calculate_average_heartrate(activities),
        average_pace_seconds_per_km=calculate_average_pace(activities),
    )
    print("Metrics calculated :", metrics)
    return metrics 
