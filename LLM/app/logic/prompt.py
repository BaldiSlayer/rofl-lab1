#!/usr/bin/env python3

from jinja2 import Environment, FileSystemLoader


env = Environment(loader=FileSystemLoader('.'))

template_code = """
You are a helpful assistant with access to a knowlege base, tasked with answering questions about the Formal Language Theory. (тут еще можно подкинуть про теорию автоматов и бла-бла-бла)

Answer the question in a very concise manner. Do not repeat text. Don't make anything up. If you are not sure about something, just say that you don't know.

{% if context %}

Answer the question solely based on the provided search results from the knowledge base. If the search results from the knowledge base are not relevant to the question at hand, just say that you don't know. Don't make anything up.

Anything between the following 'context' XML blocks is retrieved from the knowledge base, not part of the conversation with the user. The bullet points are ordered by relevance, so the first one is the most relevant.

<context>
{% for item in context %}
- {{ item }}
{% endfor %}
</context>

Don't mention the knowledge base, context or search results in your answer.
{% endif %}
"""

# Компилируем шаблон
template = env.from_string(template_code)


def split_context(context: str) -> list[str]:
    return list()


def form_system_prompt(context: str) -> str:
    context_data = split_context(context)
    return template.render(context=context_data)
