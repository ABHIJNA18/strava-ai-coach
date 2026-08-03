# This file coordinates the coaching pipeline.
# It calculates metrics, builds a prompt, calls OpenAI, and returns the summary.

from python.app.ai.openai_client import generate_summary
from python.app.metrics.running import calculate_running_metrics
from python.app.prompts.running_prompt import build_running_prompt

def generate_coaching_summary(activities):
    metrics=calculate_running_metrics(activities)
    prompt=build_running_prompt(metrics)
    summary=generate_summary(prompt)
    print("Generated summary :",summary)
    return summary