# This file sends prompts to OpenAI and returns generated text.

import os
from dotenv import load_dotenv
from openai import OpenAI

load_dotenv()


def generate_summary(prompt: str):
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        raise ValueError("OPENAI_API_KEY is not set")

    model = os.getenv("OPENAI_MODEL", "gpt-5-mini")
    client = OpenAI(api_key=api_key)

    response = client.responses.create(
        model=model,
        input=prompt,
        reasoning={
            "effort": "low",
        },
        text={
            "verbosity": "low",
        },
        max_output_tokens=800,
    )

    summary = response.output_text.strip()

    if not summary:
        raise ValueError("OpenAI returned an empty summary")

    return summary