# This file builds the prompt used to generate a running coaching summary.
# It only turns calculated metrics into text instructions.

from python.app.metrics.running import RunningMetrics

def _format_pace(seconds_per_km):
    if seconds_per_km <= 0:
        return "N/A"

    minutes = int(seconds_per_km // 60)
    seconds = int(round(seconds_per_km % 60))

    if seconds == 60:
        minutes += 1
        seconds = 0

    return f"{minutes}:{seconds:02d} min/km"

def build_running_prompt(metrics: RunningMetrics):
    total_distance_km = metrics.total_distance_meters / 1000

    prompt=  f"""

You are an experienced endurance running coach.

Write a personalized coaching summary based ONLY on the running metrics below.

Recent Running Metrics:

- Runs completed: {metrics.run_count}

- Total distance: {total_distance_km:.1f} km

- Total moving time: {metrics.total_moving_time_seconds} seconds

- Average heart rate: {metrics.average_heartrate:.0f} bpm

- Average pace: {_format_pace(metrics.average_pace_seconds_per_km)}

Your goal is to evaluate how active the runner has been and whether these metrics suggest they are building endurance consistently.

Guidelines:

- Keep the response to 3–5 sentences.

- Be encouraging but honest.

- Mention the run count, total distance, average pace, and heart rate whenever available.

- Comment on overall training consistency and activity level.

- Explain what the average pace and heart rate suggest about the effort level.

- If the training appears consistent, explain why it is helping improve endurance.

- If the volume appears low, encourage gradually increasing consistency rather than making large jumps.

- If the pace appears sustainable, mention that maintaining consistent easy running is valuable for long-term improvement.

- Focus on long-term endurance development rather than short-term speed.

- Give one practical suggestion for the next week of training.

- Do not invent any numbers or facts.

- Do not assume races, injuries, or training goals unless they are explicitly provided.

- Do not give medical advice.
- Use plain ASCII punctuation only.
- Use regular hyphens (-), straight apostrophes ('), and straight quotation marks (").
- Do not use typographic punctuation such as em dashes (—), en dashes (–), or curly quotes.
- Keep the summary under 120 words.
- Use 3-4 sentences.

Return only the coaching summary.

""".strip()
    #print("Prompt to be sent to openai:", prompt)
    return prompt