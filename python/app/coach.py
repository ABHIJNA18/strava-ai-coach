# This file coordinates the coaching pipeline.
# It calculates analytics, builds a prompt, calls OpenAI, and returns the summary.

from python.app.ai.openai_client import generate_summary
from python.app.analytics.orchestrator import calculate_running_analytics
from python.app.prompts.running_prompt import build_running_prompt


def generate_coaching_summary(activities):
    analytics = calculate_running_analytics(activities)
    prompt = build_running_prompt(analytics)
    summary = generate_summary(prompt)
    print("Generated summary :", summary)
    return summary