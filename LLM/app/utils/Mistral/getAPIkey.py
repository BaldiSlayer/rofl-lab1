import os

from ...utils.Mistral.config import api_key


def make_api_key():
    """
    place your API key in environment
    """
    os.environ["MISTRAL_API_KEY"] = api_key