# This file builds the prompt used to generate a running coaching summary.
# It only turns structured analytics into text instructions.

from python.app.analytics.models import RunningAnalytics


def _format_pace(seconds_per_km):
    if seconds_per_km <= 0:
        return "N/A"

    minutes = int(seconds_per_km // 60)
    seconds = int(round(seconds_per_km % 60))

    if seconds == 60:
        minutes += 1
        seconds = 0

    return f"{minutes}:{seconds:02d} min/km"


def _format_distance(meters):
    return f"{meters / 1000:.1f} km"


def build_running_prompt(analytics: RunningAnalytics):
    summary = analytics.summary

    prompt = f"""
You are an experienced endurance running coach.

Write a personalized coaching summary based ONLY on the analytics below.
Do not recalculate any numbers. Use the analytics as the source of truth.

Running Analytics (Last 30 Days):
- Runs analyzed: {summary.run_count}
- Total distance: {_format_distance(summary.total_distance_meters)}
- Average run distance: {_format_distance(summary.average_run_distance_meters)}
- Total moving time: {summary.total_moving_time_seconds} seconds
- Average pace: {_format_pace(summary.average_pace_seconds_per_km)}
- Average heart rate: {summary.average_heartrate:.0f} bpm
- Average cadence: {summary.average_cadence:.0f}
- Total elevation gain: {summary.total_elevation_gain_meters:.0f} meters
- Fastest run pace: {_format_pace(summary.fastest_run_pace_seconds_per_km)}
- Longest run: {_format_distance(summary.longest_run_distance_meters)}

Guidelines:
- Keep the summary under 150 words.
- Use 3-5 sentences.
- Be encouraging but honest.
- Summarize the athlete's recent training.
- Comment on overall activity level and running volume.
- Comment on average pace and heart rate.
- Comment on cadence if useful.
- Mention the longest run.
- Mention the fastest run.
- Mention elevation gain if relevant.
- Provide one practical suggestion for the next few weeks.
- Remember that a lower pace value means faster running.
- Describe the observed running volume without judging whether it is good or bad unless the supplied analytics explicitly support that conclusion.

- Avoid statements like the ones below:
- your volume is low
- your fitness is declining
- you are undertraining

- Instead use neutral language such as:
- you completed 21 km across 6 runs during the last 30 days.
- Do not invent numbers or facts.
- Do not compare weeks.
- Do not infer training streaks.
- Do not infer fitness gains or losses.
- Do not infer consistency scores.
- Do not mention trends that were not calculated.
- Do not assume races, injuries, or training goals.
- Do not give medical advice.
- Use plain ASCII punctuation only.
- Return only the coaching summary.
""".strip()

    return prompt
