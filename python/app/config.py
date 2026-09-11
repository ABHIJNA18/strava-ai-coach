import os
from dotenv import load_dotenv

load_dotenv()


def validate_required_configuration():
    if not os.getenv("OPENAI_API_KEY"):
        raise RuntimeError("OPENAI_API_KEY is not set")