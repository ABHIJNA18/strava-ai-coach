# This file coordinates the 30-day summary and personalized coaching pipelines.

from python.app.ai.openai_client import (
    generate_structured_summary,
    generate_summary,
)
from python.app.analytics.orchestrator import (
    calculate_coaching_analytics,
    calculate_running_analytics,
)
from python.app.knowledge.running.selector import (
    select_knowledge,
)
from python.app.prompts.coaching_prompt import (
    build_coaching_prompt,
)
from python.app.prompts.running_prompt import (
    build_running_prompt,
)

# needed for 30 day running summary
 
def generate_coaching_summary(activities):
    """
    Existing 30-day summary pipeline.
    """

    analytics = calculate_running_analytics(activities)

    prompt = build_running_prompt(analytics)

    summary = generate_summary(prompt)

    print("Generated summary:", summary)

    return summary

# needed for personalised coaching 

def generate_personalized_coaching(
    goal,
    activities,
):
    """
    Personalized coaching pipeline.

    Activities are used to calculate analytics.
    """

    analytics = calculate_coaching_analytics(activities)

    knowledge = select_knowledge(goal)

    prompt = build_coaching_prompt(
        goal=goal,
        analytics=analytics,
        knowledge=knowledge,
    )

    coaching = generate_structured_summary(prompt)

    print("Generated structured coaching:",coaching)

    return coaching