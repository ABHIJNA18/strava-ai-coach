

def _runs(activities):
    return [activity for activity in activities if activity.sport_type =="Run"]

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