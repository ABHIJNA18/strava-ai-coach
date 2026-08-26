# This file selects relevant running knowledge articles based on the user's goal.

from pathlib import Path


KNOWLEDGE_DIRECTORY = Path(__file__).resolve().parent


KNOWLEDGE_RULES = {
    "10k": [
        "10k",
        "10 k",
        "10-k",
    ],
    "half_marathon": [
        "half marathon",
        "21k",
        "21 k",
    ],
    "marathon": [
        "marathon",
        "42k",
        "42 k",
    ],
    "pace": [
        "pace",
        "faster",
        "speed",
        "quick",
        "time",
        "personal best",
        "pb",
    ],
    "heart_rate": [
        "heart rate",
        "heart-rate",
        "bpm",
        "pulse",
    ],
    "cadence": [
        "cadence",
        "steps per minute",
        "spm",
    ],
    "running_power": [
        "power",
        "running power",
        "watts",
    ],
    "training_load": [
        "training load",
        "volume",
        "weekly distance",
        "more mileage",
    ],
    "recovery": [
        "recovery",
        "rest",
        "fatigue",
        "tired",
        "easy",
        "injury prevention",
    ],
    "endurance": [
        "endurance",
        "aerobic",
        "long run",
        "distance",
        "stamina",
    ],
    "zone2": [
        "zone 2",
        "zone2",
        "easy aerobic",
        "aerobic base",
        "low intensity",
        "conversational",
    ],
    "intervals": [
        "interval",
        "intervals",
        "speedwork",
        "speed work",
        "faster repetitions",
        "track",
    ],
    "tempo": [
        "tempo",
        "threshold",
        "sustained pace",
        "comfortably hard",
    ],
}


DEFAULT_DOCUMENTS = [
    "endurance",
    "zone2",
    "heart_rate",
    "training_load",
    "recovery",
]


def select_knowledge(goal: str) -> str:
    goal_text = goal.lower()

    selected_names = set(DEFAULT_DOCUMENTS)

    for document_name, keywords in KNOWLEDGE_RULES.items():
        if any(
            keyword in goal_text
            for keyword in keywords
        ):
            selected_names.add(document_name)

    documents = []

    for document_name in sorted(selected_names):
        document_path = (
            KNOWLEDGE_DIRECTORY
            / f"{document_name}.md"
        )

        if not document_path.exists():
            continue

        content = document_path.read_text(
            encoding="utf-8"
        ).strip()

        if content:
            documents.append(
                f"--- {document_name} ---\n{content}"
            )

    return "\n\n".join(documents)