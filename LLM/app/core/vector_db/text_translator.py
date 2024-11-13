from deep_translator import GoogleTranslator


class TextTranslator:
    def __init__(self, dest_lang='en'):
        self.translator = GoogleTranslator(source='auto', target=dest_lang)

    def translate_text(self, text):
        return self.translator.translate(text)


class SingletonTextTranslator:
    _instance = None

    @classmethod
    def get_instance(cls):
        if cls._instance is None:
            cls._instance = TextTranslator()

        return cls._instance
