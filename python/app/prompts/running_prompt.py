# This file builds the prompt used to generate a running coaching summary.
# It only turns structured analytics into text instructions.

from datetime import datetime

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


def _format_duration(seconds):
    if seconds <= 0:
        return "N/A"

    hours = int(seconds // 3600)
    minutes = int((seconds % 3600) // 60)

    if hours > 0:
        return f"{hours}h {minutes}m"

    return f"{minutes} min"


def _format_date(date_string):
    if not date_string:
        return "N/A"

    dt = datetime.fromisoformat(
        date_string.replace("Z", "+00:00")
    )

    return dt.strftime("%d %b %Y")


def build_running_prompt(analytics: RunningAnalytics):
    summary = analytics.summary

    prompt = f"""
You are an experienced endurance running coach.

Write a personalized coaching summary based ONLY on the analytics below.

Do not recalculate any numbers.
Use the analytics as the source of truth.

Running Analytics (Last 30 Days):

- Runs analyzed: {summary.run_count}
- Total distance: {_format_distance(summary.total_distance_meters)}
- Average run distance: {_format_distance(summary.average_run_distance_meters)}
- Total moving time: {_format_duration(summary.total_moving_time_seconds)}
- Average run duration: {_format_duration(summary.average_run_duration_seconds)}
- Average pace: {_format_pace(summary.average_pace_seconds_per_km)}
- Average heart rate: {summary.average_heartrate:.0f} bpm
- Average cadence: {summary.average_cadence:.0f} spm
- Total elevation gain: {summary.total_elevation_gain_meters:.0f} meters
- Fastest Run:
- Name: {summary.fastest_run_name}
- Distance: {_format_distance(summary.fastest_run_distance_meters)}
- Pace: {_format_pace(summary.fastest_run_pace_seconds_per_km)}
- Date: {_format_date(summary.fastest_run_date)}

- Longest Run:
- Name: {summary.longest_run_name}
- Distance: {_format_distance(summary.longest_run_distance_meters)}
- Pace: {_format_pace(summary.longest_run_pace_seconds_per_km)}
- Date: {_format_date(summary.longest_run_date)}

Guidelines:

- Keep the summary under 150 words.
- Use 3-5 sentences.
- Be encouraging but honest.
- Summarize the athlete's recent training.
- Comment on overall activity level.
- Mention the total distance and average run distance.
- Mention average pace.
- Mention average heart rate.
- Mention cadence only if it adds useful coaching context.
- Mention the fastest run with it's name, date and pace and explain why it stood out.
- Mention the longest run with it's name, date and pace.
- Mention elevation only if it provides meaningful coaching context.Otherwise omit it.
- Provide one practical suggestion for the next few weeks.

Important:

- A lower pace value means faster running.
- Describe only the supplied analytics.
- Do not infer trends that were not calculated.
- Do not compare weeks or months.
- Do not infer training consistency.
- Do not infer fitness gains or losses.
- Do not infer races or goals.
- Do not infer injuries.
- Do not judge whether the running volume is good or bad.
- Do not invent numbers or facts.
- Do not give medical advice.
- Use plain ASCII punctuation only.

Return only the coaching summary.
""".strip()

    return prompt