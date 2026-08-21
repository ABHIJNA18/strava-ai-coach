# This file sends prompts to OpenAI and returns generated text.

import os
from dotenv import load_dotenv
from openai import OpenAI
import json

load_dotenv()


def generate_summary(prompt: str):
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        raise ValueError("OPENAI_API_KEY is not set")

    model = os.getenv("OPENAI_MODEL", "gpt-5.4-mini")
    client = OpenAI(api_key=api_key)

    response = client.responses.create(
        model=model,
        input=prompt,
        reasoning={
            "effort": "medium",
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

# function to generate structured coaching

def generate_structured_summary(prompt: str) -> dict:
    api_key = os.getenv("OPENAI_API_KEY")

    if not api_key:
        raise ValueError(
            "OPENAI_API_KEY is not set"
        )

    model = os.getenv(
        "OPENAI_MODEL",
        "gpt-5.4-mini",
    )

    client = OpenAI(
        api_key=api_key
    )

    response = client.responses.create(
        model=model,
        input=prompt,
        reasoning={
            "effort": "low",
        },
        text={
            "format": {
                "type": "json_schema",
                "name": "coaching_response",
                "strict": True,
                "schema": {
                    "type": "object",
                    "properties": {
                        "observations": {
                            "type": "array",
                            "items": {
                                "type": "string",
                            },
                        },
                        "progress": {
                            "type": "array",
                            "items": {
                                "type": "string",
                            },
                        },
                        "risks": {
                            "type": "array",
                            "items": {
                                "type": "string",
                            },
                        },
                        "recommendations": {
                            "type": "array",
                            "items": {
                                "type": "string",
                            },
                        },
                    },
                    "required": [
                        "observations",
                        "progress",
                        "risks",
                        "recommendations",
                    ],
                    "additionalProperties": False,
                },
            },
        },
        max_output_tokens=1600,
    )

    output = response.output_text.strip()

    if not output:
        raise ValueError(
            "OpenAI returned an empty structured response"
        )

    parsed = json.loads(output)

    return {
        "observations": parsed.get(
            "observations",
            [],
        ),
        "progress": parsed.get(
            "progress",
            [],
        ),
        "risks": parsed.get(
            "risks",
            [],
        ),
        "recommendations": parsed.get(
            "recommendations",
            [],
        ),
    }