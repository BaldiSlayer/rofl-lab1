import yaml
import sys


def parse(source_file: str):
    with open(source_file, 'r', encoding="utf-8") as file:
        data = yaml.safe_load(file)

    data_grouped_by_author = {}

    # группировка данных по автору
    for entry in data:
        author = entry['author']

        if author not in data_grouped_by_author:
            data_grouped_by_author[author] = []

        data_grouped_by_author[author].append({
            'question': entry['question'],
            'answer': entry['answer']
        })

    return data_grouped_by_author


def generate_markdown(data_grouped_by_author):
    result = "# База знаний \n\n\n"

    for author, entries in data_grouped_by_author.items():
        result += f'## {author}\n\n'

        for entry in entries:
            result += f'**Вопрос:** {entry["question"]}\n\n'

            result += f'**Ответ:** {entry["answer"]}\n\n'

        result += '\n---\n\n'

    return result


if __name__ == '__main__':
    grouped_data = parse(sys.argv[1])
    print(generate_markdown(grouped_data))
