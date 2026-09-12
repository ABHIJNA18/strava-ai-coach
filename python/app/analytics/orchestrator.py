# This file combines the analytics needed by 30 Days run summary and coaching pipeline.

from app.analytics.models import (
    CoachingAnalytics,
    RunningAnalytics,
)
from app.analytics.running import (
    calculate_run_summary,
)
from app.analytics.weekly import (
    calculate_weekly_running_analytics,
)


def calculate_running_analytics(activities):
    """
    Existing analytics path used by the 30-day summary.
    """

    return RunningAnalytics(
        summary=calculate_run_summary(
            activities
        ),
    )


def calculate_coaching_analytics(activities):
    """
    Analytics path used by personalized coaching.
    """
    #we send both 30 days running and weekly analytics 
    
    return CoachingAnalytics(
        summary=calculate_run_summary(
            activities
        ),
        weekly=calculate_weekly_running_analytics(
            activities
        ),
    )
