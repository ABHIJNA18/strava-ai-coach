# This file combines the analytics needed by the coaching pipeline.

from python.app.analytics.models import RunningAnalytics
from python.app.analytics.running import calculate_run_summary


def calculate_running_analytics(activities):
    return RunningAnalytics(
        summary=calculate_run_summary(activities),
    )
