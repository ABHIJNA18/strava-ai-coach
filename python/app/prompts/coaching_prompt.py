# This file builds the personalized coaching prompt from analytics and knowledge.

from python.app.analytics.models import (
    CoachingAnalytics,
)


def _format_distance(meters):
    return f"{meters / 1000:.1f} km"


def _format_pace(seconds_per_km):
    if seconds_per_km <= 0:
        return "N/A"

    minutes = int(seconds_per_km // 60)
    seconds = int(round(seconds_per_km % 60))

    if seconds == 60:
        minutes += 1
        seconds = 0

    return f"{minutes}:{seconds:02d} min/km"


def _format_week(week):
    return f"""
Week {week.week_start} to {week.week_end}:
- Runs: {week.run_count}
- Distance: {_format_distance(week.total_distance_meters)}
- Average run distance: {_format_distance(week.average_run_distance_meters)}
- Average pace: {_format_pace(week.average_pace_seconds_per_km)}
- Average heart rate: {week.average_heartrate:.0f} bpm
- Average cadence: {week.average_cadence:.0f} spm
- Elevation gain: {week.total_elevation_gain_meters:.0f} meters
- Longest run: {_format_distance(week.longest_run_distance_meters)}
- Fastest pace: {_format_pace(week.fastest_run_pace_seconds_per_km)}
""".strip()


def build_coaching_prompt(
    goal,
    analytics: CoachingAnalytics,
    knowledge,
):
    summary = analytics.summary

    weekly_text = "\n\n".join(
        _format_week(week)
        for week in analytics.weekly.weeks
    )

    return f"""
You are an experienced running coach.

User goal:
{goal}

The following overall analytics describe the athlete's
running activity over the last 60 days:

- Runs: {summary.run_count}
- Total distance: {_format_distance(summary.total_distance_meters)}
- Average run distance: {_format_distance(summary.average_run_distance_meters)}
- Average pace: {_format_pace(summary.average_pace_seconds_per_km)}
- Average heart rate: {summary.average_heartrate:.0f} bpm
- Average cadence: {summary.average_cadence:.0f} spm
- Total elevation gain: {summary.total_elevation_gain_meters:.0f} meters
- Longest run: {_format_distance(summary.longest_run_distance_meters)}
- Fastest pace: {_format_pace(summary.fastest_run_pace_seconds_per_km)}

The following weekly analytics describe how the athlete's
running was distributed over time:

{weekly_text}

Relevant running knowledge:

{knowledge}

Instructions:

- Use the user's goal as the main coaching context.
- Use only the supplied analytics as measured facts.
- The overall analytics describe the entire 60-day period.
- The weekly analytics describe changes in training over time.
- Identify trends only when the weekly values support them.
- Do not claim improvement unless the weekly data supports it.
- Do not calculate totals or averages yourself.
- Do not invent activities, metrics, injuries, races, or training history.
- Do not expect or request raw individual activities.
- Clearly distinguish measured data from general coaching knowledge.
- Recommendations must relate directly to the user's goal.
- Use Zone 2 or easy running when the athlete needs more aerobic endurance,
  better recovery, or lower-intensity training.
- Suggest intervals only when faster running or pace development is relevant
  and the recent training volume appears appropriate.
- Suggest tempo running only when sustained pace work is relevant to the goal.
- Do not prescribe all three types of workouts at the same time.
- Consider recovery before suggesting intervals or tempo work.
- Do not treat any Zone 2 heart-rate number as universal.
- Do not provide medical diagnosis.
- Keep every item concise.

Return a JSON object with exactly these fields:

{{
  "observations": [],
  "progress": [],
  "risks": [],
  "recommendations": []
}}
""".strip()