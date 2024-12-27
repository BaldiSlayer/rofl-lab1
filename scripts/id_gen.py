import yaml

FILE_PATH = '../data/data.yaml'


def gen_ids(data):
    new_data = []

    for current_idx, item in enumerate(data):
        new_data.append({
            'questions': item['questions'],
            'answer': item['answer'],
            'author': item['author'],
            'id': current_idx,
        })

    return new_data


def main():
    with open(FILE_PATH, 'r', encoding='utf-8') as file:
        data = yaml.safe_load(file)

    new_data = gen_ids(data)

    with open(FILE_PATH, 'w', encoding='utf-8') as f:
        yaml.dump(new_data, f, allow_unicode=True)


if __name__ == '__main__':
    main()
