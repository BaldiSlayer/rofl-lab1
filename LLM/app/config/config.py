import faiss
import yaml

from app.core.readiness_probe.readiness_probe import ReadinessProbe

class Config:
    def __init__(self):
        # база знаний
        with open('data.yaml', 'r') as file:
            self.data = yaml.safe_load(file)

        for item in self.data:
            if 'question' in item:
                item['question'] = item['question'].replace('\n', ' ').strip()
            if 'answer' in item:
                item['answer'] = item['answer'].replace('\n', ' ').strip()

        self.index = faiss.read_index('vectorized_data.faiss')

    def get_data(self):
        return self.data

    def get_index(self):
        return self.index


class SingletonConfig:
    _instance = None

    @classmethod
    def get_instance(cls):
        if cls._instance is None:
            cls._instance = Config()

            # проинициализировали конфиг и ставим ready в True
            # так как загрузили индекс в память и готовы к обработке запросов
            ReadinessProbe.set_value(True)

        return cls._instance

