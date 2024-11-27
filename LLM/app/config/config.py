import os


class Config:
    mistral_api_key = ""
    transformer_name = "sentence-transformers/all-mpnet-base-v2"

    def __init__(self):
        self.mistral_api_key = os.environ.get("MISTRAL_API_KEY")

    def get_mistral_api_key(self):
        return self.mistral_api_key

    def get_sentence_transformer_name(self):
        return self.transformer_name


class SingletonConfig:
    _instance = None

    @classmethod
    def get_instance(cls):
        if cls._instance is None:
            cls._instance = Config()

        return cls._instance

