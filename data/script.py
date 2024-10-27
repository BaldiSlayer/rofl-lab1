import yaml


def parse(source_file='data.yaml'):
    with open(source_file, 'r') as file:
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


def generate_markdown(data_grouped_by_author, filename='output.md'):
    with open(filename, 'w', encoding='utf-8') as file:
        file.write("# База знаний \n\n\n")
        for author, entries in data_grouped_by_author.items():
            file.write(f'## {author}\n\n')
            for entry in entries:
                file.write(f'**Вопрос:** {entry["question"]}\n\n')
                file.write(f'**Ответ:** {entry["answer"]}\n\n')
            file.write('\n---\n\n')


if __name__ == '__main__':
    grouped_data = parse()
    generate_markdown(grouped_data)
